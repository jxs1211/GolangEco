Original 邹俊豪 [gopher云原生](javascript:void(0);) *2023-06-13 20:33* *Posted on 广东*

收录于合集

\#k8s源码阅读6个

\#kubernetes7个

\#云原生17个

\#Go17个

## 前情回顾

### 卷一：kube-apiserver

[第1回：kube-apiserver 启动及前期调试准备](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488266&idx=1&sn=07226e3e82c90aeaf6f0d10782768a8e&scene=21#wechat_redirect)

[第2回：kube-apiserver 三种 HTTP Server 的初始化](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488299&idx=1&sn=0c9520ae2f183d0fd2ae2e3785ded266&scene=21#wechat_redirect)

[第3回：KubeAPIServer 的路由注册](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488360&idx=1&sn=3da5a3d98acee78fabc9cdb62fbe121c&scene=21#wechat_redirect)

[第4回：KubeAPIServer 的存储接口实现](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488379&idx=1&sn=e99a65d501e4e5bbc79c0c2b58898693&scene=21#wechat_redirect)

[第5回：kube-apiserver 启动 HTTP Server](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488472&idx=1&sn=7078b9b0927d918b029b857bc71be558&scene=21#wechat_redirect)

## 正文

在第 3 回核心 API 的路由注册 `InstallLegacyAPI` 过程中，最后会执行 `registerResourceHandlers` 方法：

```
// k8s.io/apiserver/pkg/endpoints/installer.go

// path 是资源的请求路径（不包含前缀和版本），例如 pods 、pods/attach 、pods/log 等
// storage 是 path 所对应的 RESTStorage 实现
func (a *APIInstaller) registerResourceHandlers(path string, storage rest.Storage, ws *restful.WebService) (*metav1.APIResource, *storageversion.ResourceInfo, error) {
 // ...

 // ...
 // 判断 storage 是否实现了 rest.Lister 接口，即支持 LIST 请求
 lister, isLister := storage.(rest.Lister)
 // 判断 storage 是否实现了 rest.Getter 接口，即支持 GET 请求
 getter, isGetter := storage.(rest.Getter)
 // ...

 switch {
 case !namespaceScoped:
  // ...
 default:
  // ...
  // 如果资源支持 LIST 请求，则添加到 actions
  actions = appendIf(actions, action{"LIST", resourcePath, resourceParams, namer, false}, isLister)
  // ...
  // 如果资源支持 GET 请求，则添加到 actions
  actions = appendIf(actions, action{"GET", itemPath, nameParams, namer, false}, isGetter)
  // ...
 }

 // 遍历 actions ，为资源支持的请求类型，设置 handler 并进行路由注册
 for _, action := range actions {
  // ...

  switch action.Verb {
  case "GET": // Get a resource.

   // 初始化资源的 GEI 请求的 handler
   var handler restful.RouteFunction
   if isGetterWithOptions {
    handler = restfulGetResourceWithOptions(getterWithOptions, reqScope, isSubresource)
   } else {
    handler = restfulGetResource(getter, reqScope)
   }

   // ...
   // 路由注册，绑定 handler
   route := ws.GET(action.Path).To(handler).
    Doc(doc).
    Param(ws.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
    Operation("read"+namespaced+kind+strings.Title(subresource)+operationSuffix).
    Produces(append(storageMeta.ProducesMIMETypes(action.Verb), mediaTypes...)...).
    Returns(http.StatusOK, "OK", producedObject).
    Writes(producedObject)
   if isGetterWithOptions {
    if err := AddObjectParams(ws, route, versionedGetOptions); err != nil {
     return nil, nil, err
    }
   }
   addParams(route, action.Params)
   routes = append(routes, route)
  case "LIST": // List all resources of a kind.
   // ...
   // 初始化资源的 LIST 请求的 handler
   handler := metrics.InstrumentRouteFunc(action.Verb, group, version, resource, subresource, requestScope, metrics.APIServerComponent, deprecated, removedRelease, restfulListResource(lister, watcher, reqScope, false, a.minRequestTimeout))
   handler = utilwarning.AddWarningsHandler(handler, warnings)
   // 路由注册，绑定 handler
   route := ws.GET(action.Path).To(handler).
    Doc(doc).
    Param(ws.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
    Operation("list"+namespaced+kind+strings.Title(subresource)+operationSuffix).
    Produces(append(storageMeta.ProducesMIMETypes(action.Verb), allMediaTypes...)...).
    Returns(http.StatusOK, "OK", versionedList).
    Writes(versionedList)
   if err := AddObjectParams(ws, route, versionedListOptions); err != nil {
    return nil, nil, err
   }
   switch {
   case isLister && isWatcher:
    doc := "list or watch objects of kind " + kind
    if isSubresource {
     doc = "list or watch " + subresource + " of objects of kind " + kind
    }
    route.Doc(doc)
   case isWatcher:
    doc := "watch objects of kind " + kind
    if isSubresource {
     doc = "watch " + subresource + "of objects of kind " + kind
    }
    route.Doc(doc)
   }
   addParams(route, action.Params)
   routes = append(routes, route)
  case "PUT": // Update a resource.
   // 同上
  case "PATCH": // Partially update a resource
   // 同上
  case "POST": // Create a resource.
   // 同上
  case "DELETE": // Delete a resource.
   // 同上
  case "DELETECOLLECTION":
   // 同上
  // deprecated in 1.11
  case "WATCH": // Watch a resource.
   // 同上
  // deprecated in 1.11
  case "WATCHLIST": // Watch all resources of a kind.
   // 同上
  case "CONNECT":
   // 同上
  default:
   return nil, nil, fmt.Errorf("unrecognized action verb: %s", action.Verb)
  }
  // ...
 }

 return &apiResource, resourceInfo, nil
}
```

其中资源的 RESTStorage 实现即 `storage` 会进行类型断言来判断是否实现了 `rest.Lister` 、`rest.Getter` 等接口：

```
// k8s.io/apiserver/pkg/registry/rest/rest.go

// Lister is an object that can retrieve resources that match the provided field and label criteria.
type Lister interface {
 // NewList returns an empty object that can be used with the List call.
 // This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
 NewList() runtime.Object
 // List selects resources in the storage which match to the selector. 'options' can be nil.
 List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error)
 // TableConvertor ensures all list implementers also implement table conversion
 TableConvertor
}

// Getter is an object that can retrieve a named RESTful resource.
type Getter interface {
 // Get finds a resource in the storage by name and returns it.
 // Although it can return an arbitrary error value, IsNotFound(err) is true for the
 // returned error value err when the specified resource is not found.
 Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error)
}
```

以 LIST 请求为例（其它请求类似），当对一个资源发起 LIST 请求时，会执行 `metrics.InstrumentRouteFunc` 方法：

```
// k8s.io/apiserver/pkg/endpoints/metrics/metrics.go

func InstrumentRouteFunc(verb, group, version, resource, subresource, scope, component string, deprecated bool, removedRelease string, routeFunc restful.RouteFunction) restful.RouteFunction {
 return restful.RouteFunction(func(req *restful.Request, response *restful.Response) {

  requestReceivedTimestamp, ok := request.ReceivedTimestampFrom(req.Request.Context())
  if !ok {
   requestReceivedTimestamp = time.Now()
  }

  // 对 response.ResponseWriter 重新包装，以便在后续的处理中对响应进行更加灵活的控制
  delegate := &ResponseWriterDelegator{ResponseWriter: response.ResponseWriter}

  rw := responsewriter.WrapForHTTP1Or2(delegate)
  response.ResponseWriter = rw

  // 执行 routeFunc 处理请求
  routeFunc(req, response)

  // 对网络请求进行监视和跟踪，不展开
  MonitorRequest(req.Request, verb, group, version, resource, subresource, scope, component, deprecated, removedRelease, delegate.Status(), delegate.ContentLength(), time.Since(requestReceivedTimestamp))
 })
}
```

跳到 `routeFunc` 方法，即传入的 `restfulListResource(lister, watcher, reqScope, false, a.minRequestTimeout)` 方法，其中 `lister` 参数就是 `storage` 经过类型断言得到的 `rest.Lister` 接口的实现：

```
// k8s.io/apiserver/pkg/endpoints/installer.go

func restfulListResource(r rest.Lister, rw rest.Watcher, scope handlers.RequestScope, forceWatch bool, minRequestTimeout time.Duration) restful.RouteFunction {
 return func(req *restful.Request, res *restful.Response) {
  // 跳转
  handlers.ListResource(r, rw, &scope, forceWatch, minRequestTimeout)(res.ResponseWriter, req.Request)
 }
}

// k8s.io/apiserver/pkg/endpoints/handlers/get.go

func ListResource(r rest.Lister, rw rest.Watcher, scope *RequestScope, forceWatch bool, minRequestTimeout time.Duration) http.HandlerFunc {
 return func(w http.ResponseWriter, req *http.Request) {
  // ... 一些请求参数的配置校验（忽略）

  if opts.Watch || forceWatch {
   // ... watch 功能（忽略）
   return
  }

  // Log only long List requests (ignore Watch).
  defer span.End(500 * time.Millisecond)
  span.AddEvent("About to List from storage")

  // 调用 rest.Lister 接口的 List 方法获取数据
  result, err := r.List(ctx, &opts)
  if err != nil {
   scope.err(err, w, req)
   return
  }
  span.AddEvent("Listing from storage done")
  defer span.AddEvent("Writing http response done", attribute.Int("count", meta.LenList(result)))
  transformResponseObject(ctx, scope, req, w, http.StatusOK, outputMediaType, result)
 }
}
```

可以看到 handler 处理流程的最后一步是直接调用 `rest.Lister` 接口的 `List` 方法来获取数据的。

到这里已经没法继续跟下去了，得回过头看看 RESTStorage 的实现过程中是如何实现 `rest.Lister` 接口的，以 `pods` path 为例，来到第 4 回创建 RESTStorage 的 `NewLegacyRESTStorage` 方法：

```
// pkg/registry/core/rest/storage_core.go

func (c LegacyRESTStorageProvider) NewLegacyRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (LegacyRESTStorage, genericapiserver.APIGroupInfo, error) {
 apiGroupInfo := genericapiserver.APIGroupInfo{
  // ...
  VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
 }
 // ...

 // Pod 资源的 RESTStorage 初始化
 podStorage, err := podstore.NewStorage(
  restOptionsGetter,
  nodeStorage.KubeletConnectionInfo,
  c.ProxyTransport,
  podDisruptionClient,
 )
 if err != nil {
  return LegacyRESTStorage{}, genericapiserver.APIGroupInfo{}, err
 }
 // ...

 storage := map[string]rest.Storage{}
 // 添加 Pod 资源到 storage
 if resource := "pods"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
  // 这个 podStorage.Pod 就是 path 为 pods 的 rest.Storage
  // 后续经过一系列的调用就会传递到路由注册的 registerResourceHandlers 方法
  storage[resource] = podStorage.Pod
  storage[resource+"/attach"] = podStorage.Attach
  storage[resource+"/status"] = podStorage.Status
  storage[resource+"/log"] = podStorage.Log
  storage[resource+"/exec"] = podStorage.Exec
  storage[resource+"/portforward"] = podStorage.PortForward
  storage[resource+"/proxy"] = podStorage.Proxy
  storage[resource+"/binding"] = podStorage.Binding
  if podStorage.Eviction != nil {
   storage[resource+"/eviction"] = podStorage.Eviction
  }
  storage[resource+"/ephemeralcontainers"] = podStorage.EphemeralContainers

 }
 // ... 添加其它资源到 storage

 if len(storage) > 0 {
  // 将 storage 存到 apiGroupInfo.VersionedResourcesStorageMap
  apiGroupInfo.VersionedResourcesStorageMap["v1"] = storage
 }

 // 返回 apiGroupInfo
 return restStorage, apiGroupInfo, nil
}
```

`podStorage.Pod` 的实现就是我们要找的，看到 `podstore.NewStorage` 方法：

```
// pkg/registry/core/pod/storage/storage.go

func NewStorage(optsGetter generic.RESTOptionsGetter, k client.ConnectionInfoGetter, proxyTransport http.RoundTripper, podDisruptionBudgetClient policyclient.PodDisruptionBudgetsGetter) (PodStorage, error) {

 store := &genericregistry.Store{
  NewFunc:                   func() runtime.Object { return &api.Pod{} },
  NewListFunc:               func() runtime.Object { return &api.PodList{} },
  PredicateFunc:             registrypod.MatchPod,
  DefaultQualifiedResource:  api.Resource("pods"),
  SingularQualifiedResource: api.Resource("pod"),

  CreateStrategy:      registrypod.Strategy,
  UpdateStrategy:      registrypod.Strategy,
  DeleteStrategy:      registrypod.Strategy,
  ResetFieldsStrategy: registrypod.Strategy,
  ReturnDeletedObject: true,

  TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
 }
 // ...
 return PodStorage{
  // Pod 初始化为 REST 对象
  Pod:                 &REST{store, proxyTransport},
  Binding:             &BindingREST{store: store},
  LegacyBinding:       &LegacyBindingREST{bindingREST},
  Eviction:            newEvictionStorage(&statusStore, podDisruptionBudgetClient),
  Status:              &StatusREST{store: &statusStore},
  EphemeralContainers: &EphemeralContainersREST{store: &ephemeralContainersStore},
  Log:                 &podrest.LogREST{Store: store, KubeletConn: k},
  Proxy:               &podrest.ProxyREST{Store: store, ProxyTransport: proxyTransport},
  Exec:                &podrest.ExecREST{Store: store, KubeletConn: k},
  Attach:              &podrest.AttachREST{Store: store, KubeletConn: k},
  PortForward:         &podrest.PortForwardREST{Store: store, KubeletConn: k},
 }, nil
}
```

`REST` 对象就是 `pods` path 对应的 RESTStorage 实现，而其中内嵌的 `genericregistry.Store` 对象实现了 `rest.Lister` 接口：

```
// pkg/registry/core/pod/storage/storage.go

type REST struct {
 *genericregistry.Store
 proxyTransport http.RoundTripper
}

// k8s.io/apiserver/pkg/registry/generic/registry/store.go

// REST 内嵌的 Store
// 实际是 Store 实现了 rest.Lister ，相当于 REST 实现了
type Store struct {
 // ...
}

func (e *Store) NewList() runtime.Object {
 return e.NewListFunc()
}

func (e *Store) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
 // ...
}
```

也就是说 path 为 `pods` 的 LIST 请求的 handler 处理流程会来到 `Store.List` 方法：

```
// k8s.io/apiserver/pkg/registry/generic/registry/store.go

func (e *Store) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
 // 标签选择器
 label := labels.Everything()
 if options != nil && options.LabelSelector != nil {
  label = options.LabelSelector
 }
 // 字段选择器
 field := fields.Everything()
 if options != nil && options.FieldSelector != nil {
  field = options.FieldSelector
 }
 // 调用 ListPredicate 获取数据
 out, err := e.ListPredicate(ctx, e.PredicateFunc(label, field), options)
 if err != nil {
  return nil, err
 }
 if e.Decorator != nil {
  e.Decorator(out)
 }
 return out, nil
}
```

继续看到 `ListPredicate` 方法：

```
// k8s.io/apiserver/pkg/registry/generic/registry/store.go

func (e *Store) ListPredicate(ctx context.Context, p storage.SelectionPredicate, options *metainternalversion.ListOptions) (runtime.Object, error) {
 if options == nil {
  // 调用 etcd 列表资源的默认参数选项
  options = &metainternalversion.ListOptions{ResourceVersion: ""}
 }
 p.Limit = options.Limit
 p.Continue = options.Continue

 // 调用 NewListFunc 获取具体的资源类型，这里是 &api.PodList{}
 // 返回结果会存到这里
 list := e.NewListFunc()
 qualifiedResource := e.qualifiedResourceFromContext(ctx)
 storageOpts := storage.ListOptions{
  ResourceVersion:      options.ResourceVersion,
  ResourceVersionMatch: options.ResourceVersionMatch,
  Predicate:            p,
  Recursive:            true,
 }

 // ...
 // 如果请求指定了 metadata.name ，则直接获取单个 object ，无需对全量数据做过滤
 if name, ok := p.MatchesSingle(); ok {
  if key, err := e.KeyFunc(ctx, name); err == nil {
   storageOpts.Recursive = false
   err := e.Storage.GetList(ctx, key, storageOpts, list)
   return list, storeerr.InterpretListError(err, qualifiedResource)
  }
  // 如果不行，则跳过优化
 }

 // 调用底层 storage.Interface 接口的 GetList 方法，查询全量数据过滤后写入到 list
 err := e.Storage.GetList(ctx, e.KeyRootFunc(ctx), storageOpts, list)
 return list, storeerr.InterpretListError(err, qualifiedResource)
}
```

`storage.Interface` 接口的实现有两种，一种是带缓存的 `cacher.Cacher` 需要指定 `--watch-cache-sizes` 参数开启，另一种是默认的不带缓存的 `etcd3.store` ，在第 4 回的 `GetRESTOptions` 方法中进行设置：

```
// k8s.io/apiserver/pkg/server/options/etcd.go

func (f *StorageFactoryRestOptionsFactory) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
 storageConfig, err := f.StorageFactory.NewConfig(resource)
 if err != nil {
  return generic.RESTOptions{}, fmt.Errorf("unable to find storage destination for %v, due to %v", resource, err.Error())
 }

 ret := generic.RESTOptions{
  StorageConfig:             storageConfig,
  Decorator:                 generic.UndecoratedStorage,
  DeleteCollectionWorkers:   f.Options.DeleteCollectionWorkers,
  EnableGarbageCollection:   f.Options.EnableGarbageCollection,
  ResourcePrefix:            f.StorageFactory.ResourcePrefix(resource),
  CountMetricPollPeriod:     f.Options.StorageConfig.CountMetricPollPeriod,
  StorageObjectCountTracker: f.Options.StorageConfig.StorageObjectCountTracker,
 }

 if f.Options.EnableWatchCache {
  sizes, err := ParseWatchCacheSizes(f.Options.WatchCacheSizes)
  if err != nil {
   return generic.RESTOptions{}, err
  }
  size, ok := sizes[resource]
  if ok && size > 0 {
   klog.Warningf("Dropping watch-cache-size for %v - watchCache size is now dynamic", resource)
  }
  if ok && size <= 0 {
   klog.V(3).InfoS("Not using watch cache", "resource", resource)
   // 不使用 cache 的 Storage 实现
   ret.Decorator = generic.UndecoratedStorage
  } else {
   klog.V(3).InfoS("Using watch cache", "resource", resource)
   // 使用 cache 的 Storage 实现
   ret.Decorator = genericregistry.StorageWithCacher()
  }
 }

 return ret, nil
}
```

来看默认的不使用 cache 的 Storage 实现方案，`generic.UndecoratedStorage` 的跳转过程在第 4 回讲过了，这里直接看到最终的 `newETCD3Storage` 方法：

```
// k8s.io/apiserver/pkg/storage/storagebackend/factory/etcd3.go

func newETCD3Storage(c storagebackend.ConfigForResource, newFunc func() runtime.Object) (storage.Interface, DestroyFunc, error) {
 // etcd v3 客户端
 client, err := newETCD3Client(c.Transport)
 if err != nil {
  stopCompactor()
  return nil, nil, err
 }
 // 返回 storage.Interface 实现
 return etcd3.New(client, c.Codec, newFunc, c.Prefix, c.GroupResource, transformer, c.Paging, c.LeaseManagerConfig), destroyFunc, nil
}

// k8s.io/apiserver/pkg/storage/etcd3/store.go

// store 就是 storage.Interface 接口的实现
type store struct {
 client              *clientv3.Client
 codec               runtime.Codec
 versioner           storage.Versioner
 transformer         value.Transformer
 pathPrefix          string
 groupResource       schema.GroupResource
 groupResourceString string
 watcher             *watcher
 pagingEnabled       bool
 leaseManager        *leaseManager
}

func New(c *clientv3.Client, codec runtime.Codec, newFunc func() runtime.Object, prefix string, groupResource schema.GroupResource, transformer value.Transformer, pagingEnabled bool, leaseManagerConfig LeaseManagerConfig) storage.Interface {
 return newStore(c, codec, newFunc, prefix, groupResource, transformer, pagingEnabled, leaseManagerConfig)
}

func newStore(c *clientv3.Client, codec runtime.Codec, newFunc func() runtime.Object, prefix string, groupResource schema.GroupResource, transformer value.Transformer, pagingEnabled bool, leaseManagerConfig LeaseManagerConfig) *store {
 // ...
 result := &store{
  client:              c,
  codec:               codec,
  versioner:           versioner,
  transformer:         transformer,
  pagingEnabled:       pagingEnabled,
  pathPrefix:          pathPrefix,
  groupResource:       groupResource,
  groupResourceString: groupResource.String(),
  watcher:             newWatcher(c, codec, groupResource, newFunc, versioner),
  leaseManager:        newDefaultLeaseManager(c, leaseManagerConfig),
 }
 return result
}

// handler 的处理流程的终点站
func (s *store) GetList(ctx context.Context, key string, opts storage.ListOptions, listObj runtime.Object) error {
 // ...
}
```

也就是说最终 handler 的处理流程是 `GetList` 方法：

```
// k8s.io/apiserver/pkg/storage/etcd3/store.go

func (s *store) GetList(ctx context.Context, key string, opts storage.ListOptions, listObj runtime.Object) error {
 // 将 key 解析成对应 etcd 的 preparedKey ，例如 /registry/pods/default
 preparedKey, err := s.prepareKey(key)
 if err != nil {
  return err
 }
 // ...

 // 转换 listObj 为指针类型 listPtr
 listPtr, err := meta.GetItemsPtr(listObj)
 if err != nil {
  return err
 }
 // 将 listPtr 转化为 reflect.Value 类型，后续直接通过 v 对象操作 listObj 的值
 v, err := conversion.EnforcePtr(listPtr)
 if err != nil || v.Kind() != reflect.Slice {
  return fmt.Errorf("need ptr to slice: %v", err)
 }

 // ... 一些查询数据的 options 配置

 for {
  startTime := time.Now()
  // 根据 preparedKey 和 options 调用 go.etcd.io/etcd/client/v3 库从 etcd 中查询数据
  getResp, err = s.client.KV.Get(ctx, preparedKey, options...)

  // ...

  // 对数据进行筛选，只选择符合条件的项目，然后将它们放入到 listObj 中保存
  for i, kv := range getResp.Kvs {
   if paging && int64(v.Len()) >= pred.Limit {
    hasMore = true
    break
   }
   lastKey = kv.Key

   data, _, err := s.transformer.TransformFromStorage(ctx, kv.Value, authenticatedDataString(kv.Key))
   if err != nil {
    return storage.NewInternalErrorf("unable to transform key %q: %v", kv.Key, err)
   }

   // 筛选后调用 v.Set 保存到 listObj 中
   if err := appendListItem(v, data, uint64(kv.ModRevision), pred, s.codec, s.versioner, newItemFunc); err != nil {
    recordDecodeError(s.groupResourceString, string(kv.Key))
    return err
   }
   numEvald++

   // 清掉数据，减少内存使用
   getResp.Kvs[i] = nil
  }

  // ... 一些退出循环的判断
 }

 // ...
}
```

到这里，handler 的处理流程就算走完了。

上个调试来结束本回，打个断点：

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmfpdUhg60Dh0sbFTD9licKMQ4ic3gwlIS5EFvme1ibljHvEaeTwBJiax8ZO4tgdicia0A57TFxDHRkdD08g/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

执行 `kubectl get pods` 命令：

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmfpdUhg60Dh0sbFTD9licKMQup8zgL0DjiczVdWQrJBPda4QLXavHwR4NhUZFTQr7tScvTuW27m1YgQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmfpdUhg60Dh0sbFTD9licKMQUGVshmk8zmFRrl8AUbGejcs2AZHM19O0wA6aL7sMv8kIiazUhdasDbg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)