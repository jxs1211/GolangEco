# K8s apiserver watch 机制浅析

Original 段朦 [CNCF](javascript:void(0);) *2022-08-23 10:16* *Posted on 香港*

![Image](https://mmbiz.qpic.cn/mmbiz_png/DmBLZYMe830KiaxnDzOj58X4gjezTbd4TvWaJIwgL2ib0mn0PVGGBw1tk8ZndibRocKC4JoficoSMWRCNibBpbpxDYg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

最近有一个业务需求，需要实现多集群watch功能，多集群的控制面apiserver需要在每个子集群的资源发生改变后跟k8s一样将资源的事件发送给客户端。客户端的client-go通过多集群控制面的kubeconfig文件新建一个informer并list and watch 所有的子集群事件，从而在一个统一的控制面观察和处理多个集群的资源变化。因此抱着学习的目的，读了几遍k8s相关的源码，感受颇深。



K8s的apiserver是k8s所有组件的流量入口，其他的所有的组件包括kube-controller-manager，kubelet，kube-scheduler等通过list-watch 机制向apiserver 发起list watch 请求，根据收到的事件处理后续的请求。watch机制本质上是使客户端和服务端建立长连接，并将服务端的变化实时发送给客户端的方式减小服务端的压力。



k8s的apiserver实现了两种长连接方式：Chunked transfer encoding(分块传输编码)和 Websocket，其中基于chunked的方式是apiserver的默认配置。k8s的watch机制的实现依赖etcd v3的watch机制，etcd v3使用的是基于 HTTP/2 的 gRPC 协议，双向流的 Watch API 设计，实现了连接多路复用。etcd 里存储的key的任何变化都会发送给客户端。



下面我们以1.24.3版本的k8s源码为例，从两个方面介绍k8s的watch机制。



作者：段朦， 中国移动云能力中心软件开发工程师，专注于云原生领域。

01

kube­apiserver对etcd的list­watch机制



说到kube-apiserver对etcd的list-watch，不得不提到一个关键的struct：cacher。为了减轻etcd的压力，kube-apiserver本身对etcd实现了list-watch机制，将所有对象的最新状态和最近的事件存放到cacher里，所有外部组件对资源的访问都经过cacher。我们看下cacher的数据结构（为了篇幅考虑，这里保留了几个关键的子结构）：

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go 

```go
type Cacher struct {
   // incoming 事件管道, 会被分发给所有的watchers
   incoming chan watchCacheEvent
   
   //storage 的底层实现
   storage storage.Interface


   // 对象类型
   objectType reflect.Type


   // watchCache 滑动窗口，维护了当前kind的所有的资源，和一个基于滑动窗口的最近的事件数组
   watchCache *watchCache
   
   // reflector list并watch etcd 并将事件和资源存到watchCache中
   reflector  *cache.Reflector
   
   // watchersBuffer 代表着所有client-go客户端跟apiserver的连接
   watchersBuffer []*cacheWatcher
   ....
}
```

下面看下cacher的创建过程

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go

```go

func NewCacherFromConfig(config Config) (*Cacher, error) {
      ...
   cacher := &Cacher{
      ...
      incoming:              make(chan watchCacheEvent, 100),
      ...
   }
      ...
    watchCache := newWatchCache(
      config.KeyFunc, cacher.processEvent, config.GetAttrsFunc, config.Versioner, config.Indexers, config.Clock, objType)
   listerWatcher := NewCacherListerWatcher(config.Storage, config.ResourcePrefix, config.NewListFunc)
   reflectorName := "storage/cacher.go:" + config.ResourcePrefix


   reflector := cache.NewNamedReflector(reflectorName, listerWatcher, obj, watchCache, 0)
   // Configure reflector's pager to for an appropriate pagination chunk size for fetching data from
   // storage. The pager falls back to full list if paginated list calls fail due to an "Expired" error.
   reflector.WatchListPageSize = storageWatchListPageSize


   cacher.watchCache = watchCache
   cacher.reflector = reflector


   go cacher.dispatchEvents() // 1


   cacher.stopWg.Add(1)
   go func() {
      defer cacher.stopWg.Done()
      defer cacher.terminateAllWatchers()
      wait.Until(
         func() {
            if !cacher.isStopped() {
               cacher.startCaching(stopCh)  // 2
            }
         }, time.Second, stopCh,
      )
   }()


   return cacher, nil
}
```

可以看到，在创建cacher的时候，也创建了watchCache（用于保存事件和所有资源）和reflactor（执行对etcd的list-watch并更新watchCache）。创建cacher的时候同时开启了两个协程，注释1 处cacher.dispatchEvents()用于从cacher的incoming管道里获取事件，并放到cacheWatcher的input里。

处理逻辑可以看下面两段代码

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go

```go

func (c *Cacher) dispatchEvents() {
   ...
   for {
      select {
      case event, ok := <-c.incoming:
         if !ok {
            return
         }
         if event.Type != watch.Bookmark {
         // 从incoming通道中获取事件，并发送给交给dispatchEvent方法处理
            c.dispatchEvent(&event)
         }
         lastProcessedResourceVersion = event.ResourceVersion
         metrics.EventsCounter.WithLabelValues(c.objectType.String()).Inc()
      ...
      case <-c.stopCh:
         return
      }
   }
}
```

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go

```go
func (c *Cacher) dispatchEvent(event *watchCacheEvent) {
   c.startDispatching(event)
   defer c.finishDispatching()
   if event.Type == watch.Bookmark {
      for _, watcher := range c.watchersBuffer {
         watcher.nonblockingAdd(event)
      }
   } else {
      wcEvent := *event
      setCachingObjects(&wcEvent, c.versioner)
      event = &wcEvent


      c.blockedWatchers = c.blockedWatchers[:0]
      // watchersBuffer 是一个数组，维护着所有client-go跟apiserver的watch连接，产生的cacheWatcher。
      for _, watcher := range c.watchersBuffer {
         if !watcher.nonblockingAdd(event) {
            c.blockedWatchers = append(c.blockedWatchers, watcher)
         }
      }
      ...
   }
}
```

 watchersBuffer 是一个数组，维护着所有client-go跟apiserver的watch连接产生的cacheWatcher，因此CacheWatcher跟发起watch请求的client-go的客户端是一对一的关系。当apiserver收到一个etcd的事件之后，会将这个事件发送到所有的cacheWatcher的input channel里。

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go 

```go
func (c *cacheWatcher) nonblockingAdd(event *watchCacheEvent) bool {
   select {
   case c.input <- event:
      return true
   default:
      return false
   }
}
```

cacherWatcher的struct结构如下

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go

```go
type cacheWatcher struct {
   input     chan *watchCacheEvent
   result    chan watch.Event
   done      chan struct{}
   filter    filterWithAttrsFunc
   stopped   bool
   forget    func()
   versioner storage.Versioner
   // The watcher will be closed by server after the deadline,
   // save it here to send bookmark events before that.
   deadline            time.Time
   allowWatchBookmarks bool
   // Object type of the cache watcher interests
   objectType reflect.Type


   // human readable identifier that helps assigning cacheWatcher
   // instance with request
   identifier string
}
```

可以看到，cacherWatcher不用于存储数据，只是实现了watch接口，并且维护了两个channel，input channel用于获取从cacher中的incoming通道中的事件，result channel 用于跟client-go的客户端交互，客户端的informer发起watch请求后，会从这个chanel里获取事件进行后续的处理。

注释2处开启了另外一个协程，cacher.startCaching(stopCh) ，实际上调用了cacher的reflector的listAndWatch方法，这里的reflector跟informer的reflector一样，list方法是获取etcd里的所有资源并对reflector的store做一次整体的replace替换，这里的store就是上面说的watchCache，watchCache实现了store接口，watch方法是watch etcd的资源，并从watcher的resultChan里拿到事件，根据事件的类型，调用watchCache的add，update，或delete方法。startCaching 执行对etcd的listAndWatch

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go

```go
func (c *Cacher) startCaching(stopChannel <-chan struct{}) {
    ...
   if err := c.reflector.ListAndWatch(stopChannel); err != nil {
      klog.Errorf("cacher (%v): unexpected ListAndWatch error: %v; reinitializing...", c.objectType.String(), err)
   }
}
```

reflector的list方法里的syncWith方法将list得到的结果替换放到watchCache里

staging/src/k8s.io/client-go/tools/cache/reflector.go

```go
func (r *Reflector) syncWith(items []runtime.Object, resourceVersion string) error {
   found := make([]interface{}, 0, len(items))
   for _, item := range items {
      found = append(found, item)
   }
   return r.store.Replace(found, resourceVersion)
}
```

reflector的list方法里的watchHandler函数传入watch etcd得到的watcher和store（即watchCache），并根据watcher的resultChan通道里收到的事件类型执行watchCache相应的方法（Add，Delete，Update）。

staging/src/k8s.io/client-go/tools/cache/reflector.go

```go

func watchHandler(start time.Time,
   w watch.Interface,
   store Store,
   expectedType reflect.Type,
   expectedGVK *schema.GroupVersionKind,
   name string,
   expectedTypeName string,
   setLastSyncResourceVersion func(string),
   clock clock.Clock,
   errc chan error,
   stopCh <-chan struct{},
) error {
   eventCount := 0


   // Stopping the watcher should be idempotent and if we return from this function there's no way
   // we're coming back in with the same watch interface.
   defer w.Stop()


loop:
   for {
      select {
      case <-stopCh:
         return errorStopRequested
      case err := <-errc:
         return err
      case event, ok := <-w.ResultChan():
         if !ok {
            break loop
         }
         if event.Type == watch.Error {
            return apierrors.FromObject(event.Object)
         }
         if expectedType != nil {
            if e, a := expectedType, reflect.TypeOf(event.Object); e != a {
               utilruntime.HandleError(fmt.Errorf("%s: expected type %v, but watch event object had type %v", name, e, a))
               continue
            }
         }
         if expectedGVK != nil {
            if e, a := *expectedGVK, event.Object.GetObjectKind().GroupVersionKind(); e != a {
               utilruntime.HandleError(fmt.Errorf("%s: expected gvk %v, but watch event object had gvk %v", name, e, a))
               continue
            }
         }
         meta, err := meta.Accessor(event.Object)
         if err != nil {
            utilruntime.HandleError(fmt.Errorf("%s: unable to understand watch event %#v", name, event))
            continue
         }
         resourceVersion := meta.GetResourceVersion()
         switch event.Type {
         case watch.Added:
            err := store.Add(event.Object)
            if err != nil {
               utilruntime.HandleError(fmt.Errorf("%s: unable to add watch event object (%#v) to store: %v", name, event.Object, err))
            }
         ...  
         setLastSyncResourceVersion(resourceVersion)
         if rvu, ok := store.(ResourceVersionUpdater); ok {
            rvu.UpdateResourceVersion(resourceVersion)
         }
         eventCount++
      }
   }
}
```

上文说到，reflector执行ListAndWatch更新watchCache保存的资源数据，下面看下watchCache的replace和add 方法，看下reflector是如何操作watchCache保存的资源的。

replace 执行了watchCache的store的replace方法，store是threadSafeMap的实现，实际上更新了底层的threadSafeMap，用于当前资源的所有实例。

staging/src/k8s.io/apiserver/pkg/storage/cacher/watch_cache.go

```go
func (w *watchCache) Replace(objs []interface{}, resourceVersion string) error {
   ... 
   if err := w.store.Replace(toReplace, resourceVersion); err != nil {
      return err
   }
   ...
}
```

add方法同样更新了底层了threadSafeMap，同时执行了一个processEvent 方法，上文说到watchCache维护了一个基于事件的数组[]*watchCacheEvent，数组采用滑动窗口算法，长度固定为100，processEvent 会一直更新这个数组，后面的事件会挤掉最前面的事件，代码如下

staging/src/k8s.io/apiserver/pkg/storage/cacher/watch_cache.go

```go
func (w *watchCache) processEvent(event watch.Event, resourceVersion uint64, updateFunc func(*storeElement) error) error {
   ...
   if err := func() error {
      // TODO: We should consider moving this lock below after the watchCacheEvent
      // is created. In such situation, the only problematic scenario is Replace(
      // happening after getting object from store and before acquiring a lock.
      // Maybe introduce another lock for this purpose.
      w.Lock()
      defer w.Unlock()
      w.updateCache(wcEvent)
      w.resourceVersion = resourceVersion
      defer w.cond.Broadcast()
    ...
      return updateFunc(elem)
   }(); err != nil {
      return err
   }
   if w.eventHandler != nil {
     w.eventHandler(wcEvent)
   }
}
```

staging/src/k8s.io/apiserver/pkg/storage/cacher/watch_cache.go

```go
func (w *watchCache) updateCache(event *watchCacheEvent) {
   w.resizeCacheLocked(event.RecordTime)
   if w.isCacheFullLocked() {
      // Cache is full - remove the oldest element.
      w.startIndex++
   }
   w.cache[w.endIndex%w.capacity] = event
   w.endIndex++
}
```

至此，cacher 创建时创建的两个协程处理过程分析完了，我们做下简单的总结，创建cacher的时候开启了两个协程：

- 第一个协程从cacher的incoming 通道里取出事件放到cacheWatcher的input通道里，而cacheWatcher是本地客户端创建一个watch请求都会生成一个，这个我们下一章再说。

- 另外一个协程主要做的事情就是reflector执行listAndWatch 方法并更新cacher里的watchCache，具体的来说，就是更新watchCache里的基于滑动窗口算法的事件数组和维护当前kind的资源的所有实例的threadSafeMap。

  这里还有两个点我们没有明确：

  1. cacher是什么时候及谁创建的 

  2. cacher的incoming 通道里的事件是哪里来的，这个通道里的时间跟reflector的listAndWatch方法里的执行对etcd的watch请求的watcher的通道里事件是否同步？

带着这些问题，我们继续看下代码，

第一个问题，可以看到apiserver在创建storage的时候创建了cacher，说明apiserver在GVK注册到apiserver的时候就创建了相应资源的cacher，这里调用链太深，因此不贴代码了。

第二个问题，我们先看下incoming通道里事件是如何来的，注意这里是cacher的processEvent 方法处理的。

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go 

```go

func (c *Cacher) processEvent(event *watchCacheEvent) {
   if curLen := int64(len(c.incoming)); c.incomingHWM.Update(curLen) {
      // Monitor if this gets backed up, and how much.
      klog.V(1).Infof("cacher (%v): %v objects queued in incoming channel.", c.objectType.String(), curLen)
   }
   c.incoming <- *event
}
```

看下这个方法是哪里用到的，由上面的NewCacherFromConfig方法可以看到是创建cacher的时候，创建watchCache的时候传入的。watchCache定义了一个eventHandler 用于处理listAndWatch收到的事件，由上面的代码 watchCache 的processEvent方法可以看到，在更新watchCache之后，会根据是否有eventHandler 执行eventHandler的func，即上面的cacher的processEvent。至此，第二个问题也变的很清晰，cacher的incoming 通道里的事件是watch etcd收到的事件更新watchCache之后处理的。这一章讲了apiserver对etcd的list和watch机制，apiserver收到事件之后本身做了缓存，并将事件发送给cacheWatcher的input通道里，由cacherWatcher处理跟客户端的连接，下一章我们讲一下本地客户端跟apiserver的watch机制的实现。

## 02 客户端对apiserver的watch机制实现

apiserver对list接口增加了一个watch参数，客户端可以向apiserver通过增加一个watch=true 参数发起watch请求：

https://{host:6443}/apis/apps/v1/namespaces/default/deployments?watch=true

apiserver 的hander 在解析到watch参数为true时，进行watch请求的处理

staging/src/k8s.io/apiserver/pkg/endpoints/handlers/get.go

```go
func ListResource(r rest.Lister, rw rest.Watcher, scope *RequestScope, forceWatch bool, minRequestTimeout time.Duration) http.HandlerFunc {
   return func(w http.ResponseWriter, req *http.Request) {
      ...
      if opts.Watch || forceWatch {
         if rw == nil {
         ...
         watcher, err := rw.Watch(ctx, &opts)
         ...
         metrics.RecordLongRunning(req, requestInfo, metrics.APIServerComponent, func() {
            serveWatch(watcher, scope, outputMediaType, req, w, timeout)
         })
         return
      }
      ...
   }
}
```

可以看到，当客户端发起watch请求时，实际上调用了watcher的watch接口，这里的watcher实际上是watch接口的实现，apiserver根据url的路径参数，针对不同的watch请求强转为不同类型的watcher实现，k8s的内置资源大都继承了REST结构体，他的底层storage就是cacher，因此这里实际上就是调用了cacher的watch方法，在看一下serveWatch的实现

staging/src/k8s.io/apiserver/pkg/endpoints/handlers/watch.go

```go
func serveWatch(watcher watch.Interface, scope *RequestScope, mediaTypeOptions negotiation.MediaTypeOptions, req *http.Request, w http.ResponseWriter, timeout time.Duration) {
   ...
   server := &WatchServer{
      Watching: watcher,
      Scope:    scope,


      UseTextFraming:  useTextFraming,
      MediaType:       mediaType,
      Framer:          framer,
      Encoder:         encoder,
      EmbeddedEncoder: embeddedEncoder,


      Fixup: func(obj runtime.Object) runtime.Object {
         result, err := transformObject(ctx, obj, options, mediaTypeOptions, scope, req)
         if err != nil {
            utilruntime.HandleError(fmt.Errorf("failed to transform object %v: %v", reflect.TypeOf(obj), err))
            return obj
         }
         // When we are transformed to a table, use the table options as the state for whether we
         // should print headers - on watch, we only want to print table headers on the first object
         // and omit them on subsequent events.
         if tableOptions, ok := options.(*metav1.TableOptions); ok {
            tableOptions.NoHeaders = true
         }
         return result
      },


      TimeoutFactory: &realTimeoutFactory{timeout},
   }


   server.ServeHTTP(w, req)
}
```

创建了一个watchServer，并执行watchServer的ServeHTTP方法，看一下ServeHTTP的实现

staging/src/k8s.io/apiserver/pkg/endpoints/handlers/watch.go

```go
func (s *WatchServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
   kind := s.Scope.Kind


   if wsstream.IsWebSocketRequest(req) {
      w.Header().Set("Content-Type", s.MediaType)
      websocket.Handler(s.HandleWS).ServeHTTP(w, req)
      return
   }


   flusher, ok := w.(http.Flusher)
   if !ok {
      err := fmt.Errorf("unable to start watch - can't get http.Flusher: %#v", w)
      utilruntime.HandleError(err)
      s.Scope.err(errors.NewInternalError(err), w, req)
      return
   }


   framer := s.Framer.NewFrameWriter(w)
   if framer == nil {
      // programmer error
      err := fmt.Errorf("no stream framing support is available for media type %q", s.MediaType)
      utilruntime.HandleError(err)
      s.Scope.err(errors.NewBadRequest(err.Error()), w, req)
      return
   }


   var e streaming.Encoder
   var memoryAllocator runtime.MemoryAllocator


   if encoder, supportsAllocator := s.Encoder.(runtime.EncoderWithAllocator); supportsAllocator {
      memoryAllocator = runtime.AllocatorPool.Get().(*runtime.Allocator)
      defer runtime.AllocatorPool.Put(memoryAllocator)
      e = streaming.NewEncoderWithAllocator(framer, encoder, memoryAllocator)
   } else {
      e = streaming.NewEncoder(framer, s.Encoder)
   }


   // ensure the connection times out
   timeoutCh, cleanup := s.TimeoutFactory.TimeoutCh()
   defer cleanup()


   // begin the stream
   w.Header().Set("Content-Type", s.MediaType)
   w.Header().Set("Transfer-Encoding", "chunked")
   w.WriteHeader(http.StatusOK)
   flusher.Flush()


   var unknown runtime.Unknown
   internalEvent := &metav1.InternalEvent{}
   outEvent := &metav1.WatchEvent{}
   buf := &bytes.Buffer{}
   ch := s.Watching.ResultChan()
   done := req.Context().Done()


   embeddedEncodeFn := s.EmbeddedEncoder.Encode
   if encoder, supportsAllocator := s.EmbeddedEncoder.(runtime.EncoderWithAllocator); supportsAllocator {
      if memoryAllocator == nil {
         // don't put the allocator inside the embeddedEncodeFn as that would allocate memory on every call.
         // instead, we allocate the buffer for the entire watch session and release it when we close the connection.
         memoryAllocator = runtime.AllocatorPool.Get().(*runtime.Allocator)
         defer runtime.AllocatorPool.Put(memoryAllocator)
      }
      embeddedEncodeFn = func(obj runtime.Object, w io.Writer) error {
         return encoder.EncodeWithAllocator(obj, w, memoryAllocator)
      }
   }


   for {
      select {
      case <-done:
         return
      case <-timeoutCh:
         return
      case event, ok := <-ch:
         if !ok {
            // End of results.
            return
         }
         metrics.WatchEvents.WithContext(req.Context()).WithLabelValues(kind.Group, kind.Version, kind.Kind).Inc()


         obj := s.Fixup(event.Object)
         if err := embeddedEncodeFn(obj, buf); err != nil {
            // unexpected error
            utilruntime.HandleError(fmt.Errorf("unable to encode watch object %T: %v", obj, err))
            return
         }


         // ContentType is not required here because we are defaulting to the serializer
         // type
         unknown.Raw = buf.Bytes()
         event.Object = &unknown
         metrics.WatchEventsSizes.WithContext(req.Context()).WithLabelValues(kind.Group, kind.Version, kind.Kind).Observe(float64(len(unknown.Raw)))


         *outEvent = metav1.WatchEvent{}


         // create the external type directly and encode it.  Clients will only recognize the serialization we provide.
         // The internal event is being reused, not reallocated so its just a few extra assignments to do it this way
         // and we get the benefit of using conversion functions which already have to stay in sync
         *internalEvent = metav1.InternalEvent(event)
         err := metav1.Convert_v1_InternalEvent_To_v1_WatchEvent(internalEvent, outEvent, nil)
         if err != nil {
            utilruntime.HandleError(fmt.Errorf("unable to convert watch object: %v", err))
            // client disconnect.
            return
         }
         if err := e.Encode(outEvent); err != nil {
            utilruntime.HandleError(fmt.Errorf("unable to encode watch object %T: %v (%#v)", outEvent, err, e))
            // client disconnect.
            return
         }
         if len(ch) == 0 {
            flusher.Flush()
         }


         buf.Reset()
      }
   }
}
```

可以看到，这里主要就是处理长连接发送给客户端的事件，读取watcher的resultChan里的事件，持续不断的放到http response的流当中，如果客户端发起的是websocket请求，则直接处理watcher的resultChan里的事件，如果是正常的http请求则需要修改请求头建立http 1.1 的长连接。

上面说到，客户端发起watch请求时，apiserver实际上调用的是cacher的Watch方法，下面看一下Watch方法

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go

```go

func (c *Cacher) Watch(ctx context.Context, key string, opts storage.ListOptions) (watch.Interface, error) {
   ...
   watcher := newCacheWatcher(chanSize, filterWithAttrsFunction(key, pred), emptyFunc, c.versioner, deadline, pred.AllowWatchBookmarks, c.objectType, identifier)
   ...
   cacheInterval, err := c.watchCache.getAllEventsSinceLocked(watchRV)
   if err != nil {
      // To match the uncached watch implementation, once we have passed authn/authz/admission,
      // and successfully parsed a resource version, other errors must fail with a watch event of type ERROR,
      // rather than a directly returned error.
      return newErrWatcher(err), nil
   }


   func() {
      c.Lock()
      defer c.Unlock()
      // Update watcher.forget function once we can compute it.
      watcher.forget = forgetWatcher(c, c.watcherIdx, triggerValue, triggerSupported)
      c.watchers.addWatcher(watcher, c.watcherIdx, triggerValue, triggerSupported)


      // Add it to the queue only when the client support watch bookmarks.
      if watcher.allowWatchBookmarks {
         c.bookmarkWatchers.addWatcher(watcher)
      }
      c.watcherIdx++
   }()


   go watcher.processInterval(ctx, cacheInterval, watchRV)
   return watcher, nil
}
```

可以看到，当客户端发起watch请求时，apiserver调用cacher的watch方法的时候创建了CacheWatcher，因此客户端的watch请求和cachWatcher是一一对应的。`cacheInterval, err := c.watchCache.getAllEventsSinceLocked(watchRV)` 是指根据客户端传过来的resourceVersion 获取watchCache滑动窗口里大于当前resourceVersion的事件，并发送给后续的协程`go watcher.processInterval(ctx, cacheInterval, watchRV)` 处理，防止客户端的watch连接断开可能导致的事件丢失。`go watcher.processInterval(ctx, cacheInterval, watchRV)` 协程中将首次watch时滑动窗口中的事件和后续watch input通道中收到的事件放到cacheWatcher的resultChan里。代码如下

staging/src/k8s.io/apiserver/pkg/storage/cacher/cacher.go

```go
func (c *cacheWatcher) processInterval(ctx context.Context, cacheInterval *watchCacheInterval, resourceVersion uint64) {
   ...
   initEventCount := 0
   /* 首次watch 获取cacheInterval 的事件并发送到resultChan*/
   for {
      event, err := cacheInterval.Next()
      if err != nil {
         klog.Warningf("couldn't retrieve watch event to serve: %#v", err)
         return
      }
      if event == nil {
         break
      }
      c.sendWatchCacheEvent(event)
      resourceVersion = event.ResourceVersion
      initEventCount++
   }
   /* 后续建立watch连接之后，将input通道中的事件发送到resultChan*/
   c.process(ctx, resourceVersion)
}
```

至此，k8s的apiserver对etcd的list 和watch 以及对客户端的list watch 处理逻辑完成了闭环，我们可以用一张图表示

![Image](https://mmbiz.qpic.cn/mmbiz_png/DmBLZYMe830raY0HxoLyd9AmBZMFhJUuPbZvVuOzGQONrpsg2ufK6h0fIicaNeUqIQNMFIPIHB6xtYib8MeBExdQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)



参考文档：

https://github.com/kubernetes/kubernetes

点击【阅读原文】到达项目仓库。