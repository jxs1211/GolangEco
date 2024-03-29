# 第5回：kube-apiserver 启动 HTTP Server

Original 邹俊豪 [gopher云原生](javascript:void(0);) *2023-06-12 16:00* *Posted on 广东*

收录于合集

\#k8s源码阅读5个

\#kubernetes6个

\#云原生16个

\#Go16个

## 前情回顾

### 卷一：kube-apiserver

[第1回：kube-apiserver 启动及前期调试准备](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488266&idx=1&sn=07226e3e82c90aeaf6f0d10782768a8e&scene=21#wechat_redirect)

[第2回：kube-apiserver 三种 HTTP Server 的初始化](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488299&idx=1&sn=0c9520ae2f183d0fd2ae2e3785ded266&scene=21#wechat_redirect)

[第3回：KubeAPIServer 的路由注册](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488360&idx=1&sn=3da5a3d98acee78fabc9cdb62fbe121c&scene=21#wechat_redirect)

[第4回：KubeAPIServer 的存储接口实现](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488379&idx=1&sn=e99a65d501e4e5bbc79c0c2b58898693&scene=21#wechat_redirect)

## 正文

前面几回都是围绕 `CreateServerChain` 创建服务调用链在讲，主要涉及了以下相关调用：

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmfZhU04Nfm5TDRBvGib9VyUHOxPtnRVXl83dLNwylIO1eCc1jJwhWiboKicwy34jLz5PpdibiahdQYjY4g/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

```
// cmd/kube-apiserver/app/server.go

func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {
 // 创建服务调用链，返回的 server 即 AggregatorServer
 server, err := CreateServerChain(completeOptions)

 // 服务启动前的准备工作
 prepared, err := server.PrepareRun()

 // 服务启动
 return prepared.Run(stopCh)
}
```

完成创建服务调用链的一系列操作之后，来到 `PrepareRun` 方法：

```
// k8s.io/kube-aggregator/pkg/apiserver/apiserver.go

func (s *APIAggregator) PrepareRun() (preparedAPIAggregator, error) {
 // ...

 // 调用 GenericAPIServer （实现了 DelegationTarget 接口）的 PrepareRun 方法
 prepared := s.GenericAPIServer.PrepareRun()

 // prepared 是 runnable 的实现
 return preparedAPIAggregator{APIAggregator: s, runnable: prepared}, nil
}
```

`GenericAPIServer` 的 `PrepareRun` 方法实际在第 2 回就已经讲过了：

```
// k8s.io/apiserver/pkg/server/genericapiserver.go

// 递归调用委托对象的 PrepareRun 方法，直到最后一个
func (s *GenericAPIServer) PrepareRun() preparedGenericAPIServer {
 // 调用下一个委托对象的 PrepareRun 方法
 s.delegationTarget.PrepareRun()

 // OpenAPI 路由的注册
 if s.openAPIConfig != nil && !s.skipOpenAPIInstallation {
  s.OpenAPIVersionedService, s.StaticOpenAPISpec = routes.OpenAPI{
   Config: s.openAPIConfig,
  }.InstallV2(s.Handler.GoRestfulContainer, s.Handler.NonGoRestfulMux)
 }

 if s.openAPIV3Config != nil && !s.skipOpenAPIInstallation {
  if utilfeature.DefaultFeatureGate.Enabled(features.OpenAPIV3) {
   s.OpenAPIV3VersionedService = routes.OpenAPI{
    Config: s.openAPIV3Config,
   }.InstallV3(s.Handler.GoRestfulContainer, s.Handler.NonGoRestfulMux)
  }
 }

 // 健康检查
 s.installHealthz()
 // 存活检查
 s.installLivez()

 // as soon as shutdown is initiated, readiness should start failing
 readinessStopCh := s.lifecycleSignals.ShutdownInitiated.Signaled()
 err := s.addReadyzShutdownCheck(readinessStopCh)
 if err != nil {
  klog.Errorf("Failed to install readyz shutdown check %s", err)
 }
 // 启动准备就绪检查
 s.installReadyz()

 return preparedGenericAPIServer{s}
}
```

总结下来就是：`server.PrepareRun()` 会依次地对 `APIExtensionsServer` 、`KubeAPIServer` 、`AggregatorServer` 进行 OpenAPI 路由的注册、健康检查、存活检查和启动准备就绪检查等工作，以便最终 **kube-apiserver** 能够顺利地运行。

准备工作完成后，就开始真正的服务启动：

```
// k8s.io/kube-aggregator/pkg/apiserver/apiserver.go

type runnable interface {
 Run(stopCh <-chan struct{}) error
}

func (s preparedAPIAggregator) Run(stopCh <-chan struct{}) error {
 return s.runnable.Run(stopCh)
}
```

这里 runnable 对应的实现很容易看出来，就是 `GenericAPIServer` 的 `PrepareRun` 方法所返回的 `preparedGenericAPIServer` 对象：

```
// k8s.io/apiserver/pkg/server/genericapiserver.go

// 实现了 runnable 接口
type preparedGenericAPIServer struct {
 *GenericAPIServer
}

// 启动 kube-apiserver 服务
func (s preparedGenericAPIServer) Run(stopCh <-chan struct{}) error {
 // ...

 // 非阻塞的方式启动服务
 stoppedCh, listenerStoppedCh, err := s.NonBlockingRun(stopHttpServerCh, shutdownTimeout)
 if err != nil {
  return err
 }

 //...
 // 阻塞等待服务退出信号
 <-stopCh
 // 服务退出后的一些收尾工作
 // ...

 // 等待服务器的优雅关闭
 <-listenerStoppedCh
 <-stoppedCh

 klog.V(1).Info("[graceful-termination] apiserver is exiting")
 return nil
}
```

继续来到 `NonBlockingRun` 方法：

```
// k8s.io/apiserver/pkg/server/genericapiserver.go

func (s preparedGenericAPIServer) NonBlockingRun(stopCh <-chan struct{}, shutdownTimeout time.Duration) (<-chan struct{}, <-chan struct{}, error) {
 internalStopCh := make(chan struct{})
 var stoppedCh <-chan struct{}
 var listenerStoppedCh <-chan struct{}
 if s.SecureServingInfo != nil && s.Handler != nil {
  var err error
  // 启动 HTTP Server
  stoppedCh, listenerStoppedCh, err = s.SecureServingInfo.Serve(s.Handler, shutdownTimeout, internalStopCh)
  if err != nil {
   close(internalStopCh)
   return nil, nil, err
  }
 }

 // ...

 return stoppedCh, listenerStoppedCh, nil
}
```

kube-apiserver 直接使用的 Go 标准库 `net/http` 和 `golang.org/x/net/http2` 来开启 HTTP/HTTP2 服务，代码比较简单，直接贴出：

```
// k8s.io/apiserver/pkg/server/secure_serving.go

// Serve runs the secure http server. It fails only if certificates cannot be loaded or the initial listen call fails.
// The actual server loop (stoppable by closing stopCh) runs in a go routine, i.e. Serve does not block.
// It returns a stoppedCh that is closed when all non-hijacked active requests have been processed.
// It returns a listenerStoppedCh that is closed when the underlying http Server has stopped listening.
func (s *SecureServingInfo) Serve(handler http.Handler, shutdownTimeout time.Duration, stopCh <-chan struct{}) (<-chan struct{}, <-chan struct{}, error) {
 if s.Listener == nil {
  return nil, nil, fmt.Errorf("listener must not be nil")
 }

 tlsConfig, err := s.tlsConfig(stopCh)
 if err != nil {
  return nil, nil, err
 }

 secureServer := &http.Server{
  Addr:           s.Listener.Addr().String(),
  Handler:        handler,
  MaxHeaderBytes: 1 << 20,
  TLSConfig:      tlsConfig,

  IdleTimeout:       90 * time.Second, // matches http.DefaultTransport keep-alive timeout
  ReadHeaderTimeout: 32 * time.Second, // just shy of requestTimeoutUpperBound
 }

 // At least 99% of serialized resources in surveyed clusters were smaller than 256kb.
 // This should be big enough to accommodate most API POST requests in a single frame,
 // and small enough to allow a per connection buffer of this size multiplied by `MaxConcurrentStreams`.
 const resourceBody99Percentile = 256 * 1024

 http2Options := &http2.Server{
  IdleTimeout: 90 * time.Second, // matches http.DefaultTransport keep-alive timeout
 }

 // shrink the per-stream buffer and max framesize from the 1MB default while still accommodating most API POST requests in a single frame
 http2Options.MaxUploadBufferPerStream = resourceBody99Percentile
 http2Options.MaxReadFrameSize = resourceBody99Percentile

 // use the overridden concurrent streams setting or make the default of 250 explicit so we can size MaxUploadBufferPerConnection appropriately
 if s.HTTP2MaxStreamsPerConnection > 0 {
  http2Options.MaxConcurrentStreams = uint32(s.HTTP2MaxStreamsPerConnection)
 } else {
  http2Options.MaxConcurrentStreams = 250
 }

 // increase the connection buffer size from the 1MB default to handle the specified number of concurrent streams
 http2Options.MaxUploadBufferPerConnection = http2Options.MaxUploadBufferPerStream * int32(http2Options.MaxConcurrentStreams)

 if !s.DisableHTTP2 {
  // apply settings to the server
  if err := http2.ConfigureServer(secureServer, http2Options); err != nil {
   return nil, nil, fmt.Errorf("error configuring http2: %v", err)
  }
 }

 // use tlsHandshakeErrorWriter to handle messages of tls handshake error
 tlsErrorWriter := &tlsHandshakeErrorWriter{os.Stderr}
 tlsErrorLogger := log.New(tlsErrorWriter, "", 0)
 secureServer.ErrorLog = tlsErrorLogger

 klog.Infof("Serving securely on %s", secureServer.Addr)
 return RunServer(secureServer, s.Listener, shutdownTimeout, stopCh)
}

// RunServer spawns a go-routine continuously serving until the stopCh is
// closed.
// It returns a stoppedCh that is closed when all non-hijacked active requests
// have been processed.
// This function does not block
// TODO: make private when insecure serving is gone from the kube-apiserver
func RunServer(
 server *http.Server,
 ln net.Listener,
 shutDownTimeout time.Duration,
 stopCh <-chan struct{},
) (<-chan struct{}, <-chan struct{}, error) {
 if ln == nil {
  return nil, nil, fmt.Errorf("listener must not be nil")
 }

 // Shutdown server gracefully.
 serverShutdownCh, listenerStoppedCh := make(chan struct{}), make(chan struct{})
 go func() {
  defer close(serverShutdownCh)
  <-stopCh
  ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
  server.Shutdown(ctx)
  cancel()
 }()

 go func() {
  defer utilruntime.HandleCrash()
  defer close(listenerStoppedCh)

  var listener net.Listener
  listener = tcpKeepAliveListener{ln}
  if server.TLSConfig != nil {
   listener = tls.NewListener(listener, server.TLSConfig)
  }

  err := server.Serve(listener)

  msg := fmt.Sprintf("Stopped listening on %s", ln.Addr().String())
  select {
  case <-stopCh:
   klog.Info(msg)
  default:
   panic(fmt.Sprintf("%s due to error: %v", msg, err))
  }
 }()

 return serverShutdownCh, listenerStoppedCh, nil
}
```

到此，HTTP Server 的启动就完成了。