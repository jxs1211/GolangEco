# kube-apiserver 三种 HTTP Server 的初始化

[k8s技术圈](javascript:void(0);) *2023-06-07 18:53* *Posted on 四川*

[第2回：kube-apiserver 三种 HTTP Server 的初始化 (qq.com)](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488299&idx=1&sn=0c9520ae2f183d0fd2ae2e3785ded266&chksm=ce56a136f921282077a500816a607884944740dca7b93775c2044af462025248a9d22527a3c5&cur_album_id=2958341226519298049&scene=190#rd)

The following article is from gopher云原生 Author 邹俊豪

[![img](http://wx.qlogo.cn/mmhead/Q3auHgzwzM62KB2Ce0NLjthrSTtH4JPKvYFRNkibGotXyOvNgIEIiavw/0)**gopher云原生**.技术log](https://mp.weixin.qq.com/s?__biz=MzU4MjQ0MTU4Ng==&mid=2247508255&idx=1&sn=027bbf95145b66b92318d4e25309a9f9&chksm=fdbaae02cacd271400ac391a2a5f4348e9b1d0bce14c6b5a908523bed1686eea1f724f7aa1b4&mpshare=1&scene=1&srcid=0607oaDoapmbRF9JYWiFVkfs&sharer_sharetime=1686149093683&sharer_shareid=85cfd1be9a9cdea1f202dd1f395a2697&key=0af7c27fedbcdd2b32b28af508d0f4b2dfddf8891efe151a462485be8bacfb0d8b3195e0a1e57c3418d0c3ffd6c78be5c9de288a24bf5078174860fba127ed2d0f4010eab2311e9ea3cf2d336ff8cec2f38065379774c15fc161281c857e2296d19b380a2be461f6bee9a9316eb2ceb9a7bcd90e1c78f6c3484d7dd1f7be1188&ascene=1&uin=MjMxOTI2NTEwMA%3D%3D&devicetype=Windows+10+x64&version=63090549&lang=en&countrycode=CN&exportkey=n_ChQIAhIQuWAm%2F5HuUXPu9e8qKzCc0xLiAQIE97dBBAEAAAAAACNPEnWs1tIAAAAOpnltbLcz9gKNyK89dVj00gJGuCnFd8NSw4z1HGWfhbo3rs6SFCc23il5a3sOc6c9HKtvobiNLOBnR66kSFI86ICV2uds3ls4J9Hiqf32fPVv%2BOwe5XA9jNnUZttGHppOTWEmSpNSqchOcfPgp3ixtM7BOwypmH%2FhIQ1ewdECnx7MgPp7HHuxXlQ6X2%2FmJO%2Fjx%2FgRPaDKwTg5trvR15MfQtNpQihICOm1IsOJSXtTCjX%2BiRW2%2B6pdX5hUyJA%2Bek9CEW1kgwNjLyMGCSQ%3D&acctmode=0&pass_ticket=ogWAHryH4Ib%2F%2FKjUJundtTKOMgSfnSCx4%2FbqPaD9%2BtPobA6zyUQcCQBjBPSPAnlE&wx_header=1#)

## 前情回顾

### 卷一：kube-apiserver

[第1回：kube-apiserver 启动及前期调试准备](http://mp.weixin.qq.com/s?__biz=MzU4MjQ0MTU4Ng==&mid=2247508242&idx=1&sn=721fa86730a870b1b4a6e654ced8ac50&chksm=fdbaae0fcacd271969115da2973ec7eb239a4421dcee3eae487ecfa17b904431e0f6a2f59499&scene=21#wechat_redirect)

## 正文

接上回，apiserver 对启动参数进行合法性校验通过后，就会调用 `Run()` 启动函数，并传递经过验证的选项配置 `completeOptions` 以及一个停止信号的通道 `stopCh` ，函数的定义如下：

```
func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {
 // To help debugging, immediately log version
 klog.Infof("Version: %+v", version.Get())

 klog.InfoS("Golang settings", "GOGC", os.Getenv("GOGC"), "GOMAXPROCS", os.Getenv("GOMAXPROCS"), "GOTRACEBACK", os.Getenv("GOTRACEBACK"))

 // 1、创建服务调用链
 server, err := CreateServerChain(completeOptions)
 if err != nil {
  return err
 }

 // 2、进行服务启动前的准备工作
 prepared, err := server.PrepareRun()
 if err != nil {
  return err
 }

 // 3、服务启动
 return prepared.Run(stopCh)
}
```

略过日志打印，可以看到整个启动函数可以分为 3 个步骤：

- `CreateServerChain` ：创建服务调用链。该函数负责创建各种不同 API Server 的配置并初始化，最后构建出完整的 API Server 链式结构
- `PrepareRun`：服务启动前的准备工作。该函数负责进行健康检查、存活检查和 OpenAPI 路由的注册工作，以便 apiserver 能够顺利地运行
- `Run`：服务启动。该函数启动 HTTP Server 实例并开始监听和处理来自客户端的请求

首先看服务调用链的创建，在这里会根据不同功能进行解耦，创建出三个不同的 API Server ：

- `AggregatorServer`：API 聚合服务。用于实现 **Kubernetes API 聚合层**[1] 的功能，当 AggregatorServer 接收到请求之后，如果发现对应的是一个 APIService 的请求，则会直接转发到对应的服务上（自行编写和部署的 API 服务器），否则则**委托**给 KubeAPIServer 进行处理
- `KubeAPIServer`：API 核心服务。实现认证、鉴权以及所有 Kubernetes 内置资源的 REST API 接口（诸如 Pod 和 Service 等资源的接口），如果请求未能找到对应的处理，则**委托**给 APIExtensionsServer 进行处理
- `APIExtensionsServer`：API 扩展服务。处理 CustomResourceDefinitions（CRD）和 Custom Resource（CR）的 REST 请求（自定义资源的接口），如果请求仍不能被处理则**委托**给 404 Handler 处理

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmfCdWVenY8LouytvOPXd6FVnNrdDvsuu0DHkvBsjoMtbthbUtZvfZbAiaeBDMMUewEp8viaMmpHzOQQ/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

```
func CreateServerChain(completedOptions completedServerRunOptions) (*aggregatorapiserver.APIAggregator, error) {
 // 为 KubeAPIServer 创建配置
 kubeAPIServerConfig, serviceResolver, pluginInitializer, err := CreateKubeAPIServerConfig(completedOptions)
 if err != nil {
  return nil, err
 }

 // 为 APIExtensionsServer 创建配置
 apiExtensionsConfig, err := createAPIExtensionsConfig(*kubeAPIServerConfig.GenericConfig, kubeAPIServerConfig.ExtraConfig.VersionedInformers, pluginInitializer, completedOptions.ServerRunOptions, completedOptions.MasterCount,
  serviceResolver, webhook.NewDefaultAuthenticationInfoResolverWrapper(kubeAPIServerConfig.ExtraConfig.ProxyTransport, kubeAPIServerConfig.GenericConfig.EgressSelector, kubeAPIServerConfig.GenericConfig.LoopbackClientConfig, kubeAPIServerConfig.GenericConfig.TracerProvider))
 if err != nil {
  return nil, err
 }

 // 1、初始化 APIExtensionsServer
 notFoundHandler := notfoundhandler.New(kubeAPIServerConfig.GenericConfig.Serializer, genericapifilters.NoMuxAndDiscoveryIncompleteKey)
 apiExtensionsServer, err := createAPIExtensionsServer(apiExtensionsConfig, genericapiserver.NewEmptyDelegateWithCustomHandler(notFoundHandler))
 if err != nil {
  return nil, err
 }

 // 2、初始化 KubeAPIServer
 kubeAPIServer, err := CreateKubeAPIServer(kubeAPIServerConfig, apiExtensionsServer.GenericAPIServer)
 if err != nil {
  return nil, err
 }

 // 为 AggregatorServer 创建配置
 aggregatorConfig, err := createAggregatorConfig(*kubeAPIServerConfig.GenericConfig, completedOptions.ServerRunOptions, kubeAPIServerConfig.ExtraConfig.VersionedInformers, serviceResolver, kubeAPIServerConfig.ExtraConfig.ProxyTransport, pluginInitializer)
 if err != nil {
  return nil, err
 }
 // 3、初始化 AggregatorServer
 aggregatorServer, err := createAggregatorServer(aggregatorConfig, kubeAPIServer.GenericAPIServer, apiExtensionsServer.Informers)
 if err != nil {
  // we don't need special handling for innerStopCh because the aggregator server doesn't create any go routines
  return nil, err
 }

 // 返回的是最后的 AggregatorServer
 return aggregatorServer, nil
}
```

这三个服务通过**委托模式**连接在一起，形成了一个链式结构，函数最后返回的 `AggregatorServer` 服务为头结点。把这个逻辑搞懂。

这三个服务的类型

- `APIExtensionsServer` ：`*apiextensionsapiserver.CustomResourceDefinitions`
- `KubeAPIServer` ：`*controlplane.Instance`
- `AggregatorServer` ：`*aggregatorapiserver.APIAggregator`

```
// APIExtensionsServer 类型
type CustomResourceDefinitions struct {
 GenericAPIServer *genericapiserver.GenericAPIServer

 // ...
}

// KubeAPIServer 类型
type Instance struct {
 GenericAPIServer *genericapiserver.GenericAPIServer

 // ...
}

// AggregatorServer 类型
type APIAggregator struct {
 GenericAPIServer *genericapiserver.GenericAPIServer

 // ...
}
```

都有一个共同点，包含了 `GenericAPIServer` 成员，而该成员实现了 `DelegationTarget` 接口：

```
type DelegationTarget interface {
 // ...

 // 获取委托链中的下一个委托对象
 NextDelegate() DelegationTarget

 // 执行 API Server 启动前的准备工作
 PrepareRun() preparedGenericAPIServer

 // ...
}

// 实现了 DelegationTarget 接口
type GenericAPIServer struct {
  // ...
 // delegationTarget是链中的下一个委托对象
 delegationTarget DelegationTarget
  // ...
}

// 实现 NextDelegate 方法
func (s *GenericAPIServer) NextDelegate() DelegationTarget {
 return s.delegationTarget
}

// 实现 PrepareRun 方法，会递归调用委托对象的 PrepareRun 方法，直到最后一个
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

基于委托模式，重新看 CreateServerChain 函数，从尾节点开始依次创建 API Server 委托对象：

```
// 0、初始化 404 Handler Server
notFoundHandler := notfoundhandler.New(kubeAPIServerConfig.GenericConfig.Serializer, genericapifilters.NoMuxAndDiscoveryIncompleteKey)

// 1、初始化 APIExtensionsServer ，传入 404 Handler Server 作为下一个委托对象
apiExtensionsServer, err := createAPIExtensionsServer(apiExtensionsConfig, genericapiserver.NewEmptyDelegateWithCustomHandler(notFoundHandler))

// 2、初始化 KubeAPIServer ，传入 APIExtensionsServer 作为下一个委托对象
kubeAPIServer, err := CreateKubeAPIServer(kubeAPIServerConfig, apiExtensionsServer.GenericAPIServer)

// 3、初始化 AggregatorServer ，传入 KubeAPIServer 作为下一个委托对象
aggregatorServer, err := createAggregatorServer(aggregatorConfig, kubeAPIServer.GenericAPIServer, apiExtensionsServer.Informers)

// 4、返回 AggregatorServer
return aggregatorServer, nil
```

先看 createAPIExtensionsServer ：

```
func createAPIExtensionsServer(apiextensionsConfig *apiextensionsapiserver.Config, delegateAPIServer genericapiserver.DelegationTarget) (*apiextensionsapiserver.CustomResourceDefinitions, error) {
 return apiextensionsConfig.Complete().New(delegateAPIServer)
}

func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*CustomResourceDefinitions, error) {
 // 创建 APIExtensionsServer 委托对象，并传入所指向的下一个委托对象，这里是 404 Handler Server
  genericServer, err := c.GenericConfig.New("apiextensions-apiserver", delegationTarget)

  // ...
}
```

再看 CreateKubeAPIServer：

```
func CreateKubeAPIServer(kubeAPIServerConfig *controlplane.Config, delegateAPIServer genericapiserver.DelegationTarget) (*controlplane.Instance, error) {
 return kubeAPIServerConfig.Complete().New(delegateAPIServer)
}

func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*Instance, error) {
 // ...

 // 创建 KubeAPIServer 委托对象，并传入所指向的下一个委托对象，这里是 APIExtensionsServer
 s, err := c.GenericConfig.New("kube-apiserver", delegationTarget)

  // ...
}
```

最后看 createAggregatorServer：

```
func createAggregatorServer(aggregatorConfig *aggregatorapiserver.Config, delegateAPIServer genericapiserver.DelegationTarget, apiExtensionInformers apiextensionsinformers.SharedInformerFactory) (*aggregatorapiserver.APIAggregator, error) {
 aggregatorServer, err := aggregatorConfig.Complete().NewWithDelegate(delegateAPIServer)

 // ...
}

func (c completedConfig) NewWithDelegate(delegationTarget genericapiserver.DelegationTarget) (*APIAggregator, error) {
 // 创建 AggregatorServer 委托对象，并传入所指向的下一个委托对象，这里是 KubeAPIServer
 genericServer, err := c.GenericConfig.New("kube-aggregator", delegationTarget)

 // ...
}
```

可以看到三个服务的初始化过程都是一样，调用 `c.GenericConfig.New("server name", delegationTarget)` 方法：

```
func (c completedConfig) New(name string, delegationTarget DelegationTarget) (*GenericAPIServer, error) {
 // ...

  // 创建 API Server 的 Handler 处理器
 apiServerHandler := NewAPIServerHandler(name, c.Serializer, handlerChainBuilder, delegationTarget.UnprotectedHandler())

 s := &GenericAPIServer{
  // 保存下一个委托对象
  delegationTarget:               delegationTarget,
  // 保存当前委托对象 API Server 的 Handler 处理器
  Handler:                        apiServerHandler,
  // ...
 }

 // ...
 return s, nil
}
```

在继续看 NewAPIServerHandler 函数之前，先了解一下 **github.com/emicklei/go-restful**[2] ，因为 apiserver 就是使用这个库实现的 RESTful API 。

直接用一个小例子搞懂 apiserver 中 go-restful 的使用：

```
package main

import (
 "net/http"

 "github.com/emicklei/go-restful/v3"

 "k8s.io/apiserver/pkg/server"
 "k8s.io/kubernetes/pkg/api/legacyscheme"
)

func main() {
 // 创建 API Server 的 Handler 处理器
 handler := server.NewAPIServerHandler(
  "test-server",
  legacyscheme.Codecs,
  func(apiHandler http.Handler) http.Handler {
   return apiHandler
  },
  nil)

 // 注册路由
 testApisV1 := new(restful.WebService).Path("/apis/test/v1")
 {
  testApisV1.Route(testApisV1.GET("hello").To(
   func(req *restful.Request, resp *restful.Response) {
    resp.WriteAsJson(map[string]interface{}{"k": "v"})
   },
  )).Doc("hello endpoint")
 }

 // 路由添加到 GoRestfulContainer
 handler.GoRestfulContainer.Add(testApisV1)

 // 启动监听服务
 panic(http.ListenAndServe(":8080", handler))
}

// $ curl 127.0.0.1:8080/apis/test/v1/hello
// {
// "k": "v"
// }
```

继续看源码，可以看到 apiserver 实际只是对 go-restful 进行了一些简单的封装，使用了其中的一些基本方法。不过因为 go-restful 对于一些 API 有兼容性问题，因此引入了 Director 机制来选择是使用 GoRestfulContainer 还是 NonGoRestfulMux ：

```
// 处理顺序：FullHandlerChain -> Director（选择） -> {GoRestfulContainer,NonGoRestfulMux}
type APIServerHandler struct {
 // 完整的处理程序链，包含了所有的中间件和处理程序
 FullHandlerChain http.Handler

 // 注册和管理 go-restful 的路由和处理程序
 GoRestfulContainer *restful.Container
 // 注册和管理非 go-restful 的路由和处理程序
 NonGoRestfulMux *mux.PathRecorderMux

 // 根据已注册的 web 服务检查来选择使用哪个处理程序（gorestful 或非 gorestful）
 Director http.Handler
}

func NewAPIServerHandler(name string, s runtime.NegotiatedSerializer, handlerChainBuilder HandlerChainBuilderFn, notFoundHandler http.Handler) *APIServerHandler {
 // 非 go-restful ，以下称 mux 的初始化
 nonGoRestfulMux := mux.NewPathRecorderMux(name)
 if notFoundHandler != nil {
  // 自定义 404 处理器
  nonGoRestfulMux.NotFoundHandler(notFoundHandler)
 }

 // go-restful 的初始化
 gorestfulContainer := restful.NewContainer()
 gorestfulContainer.ServeMux = http.NewServeMux()
 gorestfulContainer.Router(restful.CurlyRouter{}) // e.g. for proxy/{kind}/{name}/{*}
 gorestfulContainer.RecoverHandler(func(panicReason interface{}, httpWriter http.ResponseWriter) {
  logStackOnRecover(s, panicReason, httpWriter)
 })
 gorestfulContainer.ServiceErrorHandler(func(serviceErr restful.ServiceError, request *restful.Request, response *restful.Response) {
  serviceErrorHandler(s, serviceErr, request, response)
 })

 // 声明 director ，后续使用 director 来决定是使用 mux 还是 go-restful
 director := director{
  name:               name,
  goRestfulContainer: gorestfulContainer,
  nonGoRestfulMux:    nonGoRestfulMux,
 }

 // 返回 APIServerHandler
 return &APIServerHandler{
  // 在 director 基础上增加中间件，得到完整的处理程序链
  FullHandlerChain:   handlerChainBuilder(director),
  GoRestfulContainer: gorestfulContainer,
  NonGoRestfulMux:    nonGoRestfulMux,
  Director:           director,
 }
}
```

至此，三个服务的初始化暂告一段落。后续将继续完成路由注册，服务启动等步骤。

### 参考资料

[1]Kubernetes API 聚合层: https://kubernetes.io/zh-cn/docs/concepts/extend-kubernetes/api-extension/apiserver-aggregation/

[2]github.com/emicklei/go-restful: http://github.com/emicklei/go-restful

