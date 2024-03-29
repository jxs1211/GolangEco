# [第7回：kube-apiserver API 认证和鉴权]([第7回：kube-apiserver API 认证和鉴权 (qq.com)](https://mp.weixin.qq.com/s/H61V0hmWZTt7tFDL-vXG2Q))

Original 邹俊豪 [gopher云原生](javascript:void(0);) *2023-06-19 12:12* *Posted on 广东*

收录于合集

\#k8s源码阅读7个

\#kubernetes8个

\#云原生18个

\#Go18个

## 前情回顾

### 卷一：kube-apiserver

[第1回：kube-apiserver 启动及前期调试准备](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488266&idx=1&sn=07226e3e82c90aeaf6f0d10782768a8e&scene=21#wechat_redirect)

[第2回：kube-apiserver 三种 HTTP Server 的初始化](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488299&idx=1&sn=0c9520ae2f183d0fd2ae2e3785ded266&scene=21#wechat_redirect)

[第3回：KubeAPIServer 的路由注册](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488360&idx=1&sn=3da5a3d98acee78fabc9cdb62fbe121c&scene=21#wechat_redirect)

[第4回：KubeAPIServer 的存储接口实现](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488379&idx=1&sn=e99a65d501e4e5bbc79c0c2b58898693&scene=21#wechat_redirect)

[第5回：kube-apiserver 启动 HTTP Server](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488472&idx=1&sn=7078b9b0927d918b029b857bc71be558&scene=21#wechat_redirect)

[第6回：kube-apiserver handler 处理流程](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488494&idx=1&sn=1f0a6a3579f359a5e29cafcea1ed077e&scene=21#wechat_redirect)

## 正文

用户使用 kubectl、客户端库或构造 REST 请求来访问 kube-apiserver 的 API 时，会先经过 **认证**（Authentication）、**鉴权**（Authorization）、**准入控制**（Admission Controllers）三个阶段后才进入到第 6 回的 handler 处理流程中操作资源对象。

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmcRhMPmyo6VpT21cs3nhSaOxMBGtib0QLpn2oh3oh5qxKOpM64L7piaiagxJyOYV4ib3Pk9TGHu31Ar5w/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

- 认证：针对请求的认证，确认是否有访问集群的权限
- 鉴权：针对资源的授权，确认是否有对资源操作的权限
- 准入控制：针对请求的资源内容进行校验、修改、拒绝（之前讲过的 Sidecar 注入原理）

本回先看认证和鉴权阶段，如果熟悉 gin 框架的话，其实现原理和 gin 中间件是一样的，本质只是一个个的 handler ，在整个 HandlerChain 中，位于业务逻辑 handler 的前面优先处理。

在第 2 回讲到，不论是 AggregatorServer 、KubeAPIServer 还是 APIExtensionsServer ，这三个服务都是调用的 `NewAPIServerHandler` 方法进行 API Server 的 Handler 处理器创建：

```
// k8s.io/apiserver/pkg/server/handler.go

func NewAPIServerHandler(name string, s runtime.NegotiatedSerializer, handlerChainBuilder HandlerChainBuilderFn, notFoundHandler http.Handler) *APIServerHandler {
 // ...

 // FullHandlerChain 是在 director 基础上增加中间件，得到完整的处理程序链
 // 处理顺序：FullHandlerChain -> Director（选择） -> {GoRestfulContainer,NonGoRestfulMux}
 return &APIServerHandler{
  FullHandlerChain:   handlerChainBuilder(director),
  GoRestfulContainer: gorestfulContainer,
  NonGoRestfulMux:    nonGoRestfulMux,
  Director:           director,
 }
}
```

其中 FullHandlerChain 是通过参数传递的 `handlerChainBuilder` 方法在 director 选择器的基础上增加一系列的中间件（包括认证和鉴权 handler ），从而得到的完整的程序处理链，看到 `handlerChainBuilder` 的定义：

```
// k8s.io/apiserver/pkg/server/config.go

func (c completedConfig) New(name string, delegationTarget DelegationTarget) (*GenericAPIServer, error) {
 // ...

 handlerChainBuilder := func(handler http.Handler) http.Handler {
  // 调用的 completedConfig.BuildHandlerChainFunc 方法
  return c.BuildHandlerChainFunc(handler, c.Config)
 }

 // ...
 // 创建 API Server 的 Handler 处理器，传入 handlerChainBuilder
 apiServerHandler := NewAPIServerHandler(name, c.Serializer, handlerChainBuilder, delegationTarget.UnprotectedHandler())

 // ...
}
```

`BuildHandlerChainFunc` 方法很容易找到，不啰嗦了，就在创建服务调用链中的创建配置的地方：

```
// cmd/kube-apiserver/app/server.go
func CreateServerChain(completedOptions completedServerRunOptions) (*aggregatorapiserver.APIAggregator, error) {
 // 为 KubeAPIServer 创建配置
 kubeAPIServerConfig, serviceResolver, pluginInitializer, err := CreateKubeAPIServerConfig(completedOptions)
 if err != nil {
  return nil, err
 }

 // ...
}

// cmd/kube-apiserver/app/server.go
func CreateKubeAPIServerConfig(s completedServerRunOptions) (
 *controlplane.Config,
 aggregatorapiserver.ServiceResolver,
 []admission.PluginInitializer,
 error,
) {
 // ...
 // 跳到 buildGenericConfig
 genericConfig, versionedInformers, serviceResolver, pluginInitializers, admissionPostStartHook, storageFactory, err := buildGenericConfig(s.ServerRunOptions, proxyTransport)
 // ...
}

// cmd/kube-apiserver/app/server.go
func buildGenericConfig(
 s *options.ServerRunOptions,
 proxyTransport *http.Transport,
) (
 genericConfig *genericapiserver.Config,
 versionedInformers clientgoinformers.SharedInformerFactory,
 serviceResolver aggregatorapiserver.ServiceResolver,
 pluginInitializers []admission.PluginInitializer,
 admissionPostStartHook genericapiserver.PostStartHookFunc,
 storageFactory *serverstorage.DefaultStorageFactory,
 lastErr error,
) {
 // 初始化配置
 genericConfig = genericapiserver.NewConfig(legacyscheme.Codecs)
 // ...
}

// k8s.io/apiserver/pkg/server/config.go
func NewConfig(codecs serializer.CodecFactory) *Config {
 // ...

 return &Config{
  // BuildHandlerChainFunc 使用的是 DefaultBuildHandlerChain
  BuildHandlerChainFunc:          DefaultBuildHandlerChain,
  // ...
 }
}
```

所以 `handlerChainBuilder` 方法的实现就是 `DefaultBuildHandlerChain` 方法：

```
// k8s.io/apiserver/pkg/server/config.go

// 注册各种中间件，这里中间件的实际执行顺序和声明顺序是相反的，先认证，再鉴权
func DefaultBuildHandlerChain(apiHandler http.Handler, c *Config) http.Handler {

 // 鉴权中间件，传入 c.Authorization.Authorizer 鉴权器
 handler = genericapifilters.WithAuthorization(handler, c.Authorization.Authorizer, c.Serializer)

 // 认证中间件，传入 c.Authentication.Authenticator 认证器
 handler = genericapifilters.WithAuthentication(handler, c.Authentication.Authenticator, failedHandler, c.Authentication.APIAudiences, c.Authentication.RequestHeaderConfig)

 // 跨域中间件
 handler = genericfilters.WithCORS(handler, c.CorsAllowedOriginList, nil, nil, nil, "true")
 // Panic Recovery 中间件
 handler = genericfilters.WithPanicRecovery(handler, c.RequestInfoResolver)

 // ... 等等还有很多的中间件

 return handler
}
```

`WithAuthorization` 和 `WithAuthentication` 方法很简单，就是普通的中间件处理逻辑的写法，在其中调用认证器进行认证和鉴权器进行鉴权：

```
// k8s.io/apiserver/pkg/endpoints/filters/authorization.go
func WithAuthorization(handler http.Handler, a authorizer.Authorizer, s runtime.NegotiatedSerializer) http.Handler {
 if a == nil {
  klog.Warning("Authorization is disabled")
  return handler
 }
 return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
  ctx := req.Context()

  // 调用鉴权器进行鉴权操作
  authorized, reason, err := a.Authorize(ctx, attributes)
  // 返回 DecisionAllow 代表鉴权通过，继续前往下一个 handler
  if authorized == authorizer.DecisionAllow {
   audit.AddAuditAnnotations(ctx,
    decisionAnnotationKey, decisionAllow,
    reasonAnnotationKey, reason)
   handler.ServeHTTP(w, req)
   return
  }
  if err != nil {
   audit.AddAuditAnnotation(ctx, reasonAnnotationKey, reasonError)
   responsewriters.InternalError(w, req, err)
   return
  }
  // 鉴权失败，返回 403
  klog.V(4).InfoS("Forbidden", "URI", req.RequestURI, "reason", reason)
  audit.AddAuditAnnotations(ctx,
   decisionAnnotationKey, decisionForbid,
   reasonAnnotationKey, reason)
  responsewriters.Forbidden(ctx, attributes, w, req, reason, s)
 })
}

// k8s.io/apiserver/pkg/endpoints/filters/authentication.go
func WithAuthentication(handler http.Handler, auth authenticator.Request, failed http.Handler, apiAuds authenticator.Audiences, requestHeaderConfig *authenticatorfactory.RequestHeaderConfig) http.Handler {
 return withAuthentication(handler, auth, failed, apiAuds, requestHeaderConfig, recordAuthMetrics)
}

func withAuthentication(handler http.Handler, auth authenticator.Request, failed http.Handler, apiAuds authenticator.Audiences, requestHeaderConfig *authenticatorfactory.RequestHeaderConfig, metrics recordMetrics) http.Handler {
 if auth == nil {
  klog.Warning("Authentication is disabled")
  return handler
 }
 standardRequestHeaderConfig := &authenticatorfactory.RequestHeaderConfig{
  UsernameHeaders:     headerrequest.StaticStringSlice{"X-Remote-User"},
  GroupHeaders:        headerrequest.StaticStringSlice{"X-Remote-Group"},
  ExtraHeaderPrefixes: headerrequest.StaticStringSlice{"X-Remote-Extra-"},
 }
 return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
  authenticationStart := time.Now()

  if len(apiAuds) > 0 {
   req = req.WithContext(authenticator.WithAudiences(req.Context(), apiAuds))
  }
  // 调用认证器进行认证操作
  resp, ok, err := auth.AuthenticateRequest(req)
  authenticationFinish := time.Now()
  defer func() {
   metrics(req.Context(), resp, ok, err, apiAuds, authenticationStart, authenticationFinish)
  }()
  // 认证失败
  if err != nil || !ok {
   if err != nil {
    klog.ErrorS(err, "Unable to authenticate the request")
   }
   failed.ServeHTTP(w, req)
   return
  }

  if !audiencesAreAcceptable(apiAuds, resp.Audiences) {
   err = fmt.Errorf("unable to match the audience: %v , accepted: %v", resp.Audiences, apiAuds)
   klog.Error(err)
   failed.ServeHTTP(w, req)
   return
  }

  // 认证成功，前往下一个 handler
  // authorization header is not required anymore in case of a successful authentication.
  req.Header.Del("Authorization")

  // delete standard front proxy headers
  headerrequest.ClearAuthenticationHeaders(
   req.Header,
   standardRequestHeaderConfig.UsernameHeaders,
   standardRequestHeaderConfig.GroupHeaders,
   standardRequestHeaderConfig.ExtraHeaderPrefixes,
  )

  // also delete any custom front proxy headers
  if requestHeaderConfig != nil {
   headerrequest.ClearAuthenticationHeaders(
    req.Header,
    requestHeaderConfig.UsernameHeaders,
    requestHeaderConfig.GroupHeaders,
    requestHeaderConfig.ExtraHeaderPrefixes,
   )
  }

  req = req.WithContext(genericapirequest.WithUser(req.Context(), resp.User))
  handler.ServeHTTP(w, req)
 })
}
```

直接看其核心调用的认证器和鉴权器，回到 `buildGenericConfig` 方法：

```
// cmd/kube-apiserver/app/server.go

func buildGenericConfig(
 s *options.ServerRunOptions,
 proxyTransport *http.Transport,
) (
 genericConfig *genericapiserver.Config,
 versionedInformers clientgoinformers.SharedInformerFactory,
 serviceResolver aggregatorapiserver.ServiceResolver,
 pluginInitializers []admission.PluginInitializer,
 admissionPostStartHook genericapiserver.PostStartHookFunc,
 storageFactory *serverstorage.DefaultStorageFactory,
 lastErr error,
) {
 // 初始化配置
 genericConfig = genericapiserver.NewConfig(legacyscheme.Codecs)

 // ...

 // 认证器的初始化
 if lastErr = s.Authentication.ApplyTo(&genericConfig.Authentication, genericConfig.SecureServing, genericConfig.EgressSelector, genericConfig.OpenAPIConfig, genericConfig.OpenAPIV3Config, clientgoExternalClient, versionedInformers); lastErr != nil {
  return
 }

 // 鉴权器的初始化
 genericConfig.Authorization.Authorizer, genericConfig.RuleResolver, err = BuildAuthorizer(s, genericConfig.EgressSelector, versionedInformers)
 if err != nil {
  lastErr = fmt.Errorf("invalid authorization config: %v", err)
  return
 }
 if !sets.NewString(s.Authorization.Modes...).Has(modes.ModeRBAC) {
  genericConfig.DisabledPostStartHooks.Insert(rbacrest.PostStartHookName)
 }

 // ...

 return
}
```

先看认证器的初始化：

```
// pkg/kubeapiserver/options/authentication.go
func (o *BuiltInAuthenticationOptions) ApplyTo(authInfo *genericapiserver.AuthenticationInfo, secureServing *genericapiserver.SecureServingInfo, egressSelector *egressselector.EgressSelector, openAPIConfig *openapicommon.Config, openAPIV3Config *openapicommon.Config, extclient kubernetes.Interface, versionedInformer informers.SharedInformerFactory) error {
 // ...

 // 认证器的初始化
 authInfo.Authenticator, openAPIConfig.SecurityDefinitions, err = authenticatorConfig.New()

 // ...
 return nil
}

// pkg/kubeapiserver/authenticator/config.go
func (config Config) New() (authenticator.Request, *spec.SecurityDefinitions, error) {
 // 用来存放一系列的认证器
 var authenticators []authenticator.Request
 // 用来存放 Token 相关的认证器，最后会组合成一个 Token 认证器一起放入 authenticators
 var tokenAuthenticators []authenticator.Token
 securityDefinitions := spec.SecurityDefinitions{}

 // RequestHeader 认证器
 if config.RequestHeaderConfig != nil {
  requestHeaderAuthenticator := headerrequest.NewDynamicVerifyOptionsSecure(
   config.RequestHeaderConfig.CAContentProvider.VerifyOptions,
   config.RequestHeaderConfig.AllowedClientNames,
   config.RequestHeaderConfig.UsernameHeaders,
   config.RequestHeaderConfig.GroupHeaders,
   config.RequestHeaderConfig.ExtraHeaderPrefixes,
  )
  authenticators = append(authenticators, authenticator.WrapAudienceAgnosticRequest(config.APIAudiences, requestHeaderAuthenticator))
 }

 // ClientCA 认证器
 if config.ClientCAContentProvider != nil {
  certAuth := x509.NewDynamic(config.ClientCAContentProvider.VerifyOptions, x509.CommonNameUserConversion)
  authenticators = append(authenticators, certAuth)
 }

 // TokenAuth 认证器
 if len(config.TokenAuthFile) > 0 {
  tokenAuth, err := newAuthenticatorFromTokenFile(config.TokenAuthFile)
  if err != nil {
   return nil, nil, err
  }
  tokenAuthenticators = append(tokenAuthenticators, authenticator.WrapAudienceAgnosticToken(config.APIAudiences, tokenAuth))
 }
 // ServiceAccountAuth 认证器
 if len(config.ServiceAccountKeyFiles) > 0 {
  serviceAccountAuth, err := newLegacyServiceAccountAuthenticator(config.ServiceAccountKeyFiles, config.ServiceAccountLookup, config.APIAudiences, config.ServiceAccountTokenGetter, config.SecretsWriter)
  if err != nil {
   return nil, nil, err
  }
  tokenAuthenticators = append(tokenAuthenticators, serviceAccountAuth)
 }
 if len(config.ServiceAccountIssuers) > 0 {
  serviceAccountAuth, err := newServiceAccountAuthenticator(config.ServiceAccountIssuers, config.ServiceAccountKeyFiles, config.APIAudiences, config.ServiceAccountTokenGetter)
  if err != nil {
   return nil, nil, err
  }
  tokenAuthenticators = append(tokenAuthenticators, serviceAccountAuth)
 }
 // BootstrapToken 认证器
 if config.BootstrapToken {
  if config.BootstrapTokenAuthenticator != nil {
   // TODO: This can sometimes be nil because of
   tokenAuthenticators = append(tokenAuthenticators, authenticator.WrapAudienceAgnosticToken(config.APIAudiences, config.BootstrapTokenAuthenticator))
  }
 }
 // OIDC 认证器
 if len(config.OIDCIssuerURL) > 0 && len(config.OIDCClientID) > 0 {
  // TODO(enj): wire up the Notifier and ControllerRunner bits when OIDC supports CA reload
  var oidcCAContent oidc.CAContentProvider
  if len(config.OIDCCAFile) != 0 {
   var oidcCAErr error
   oidcCAContent, oidcCAErr = staticCAContentProviderFromFile("oidc-authenticator", config.OIDCCAFile)
   if oidcCAErr != nil {
    return nil, nil, oidcCAErr
   }
  }

  oidcAuth, err := newAuthenticatorFromOIDCIssuerURL(oidc.Options{
   IssuerURL:            config.OIDCIssuerURL,
   ClientID:             config.OIDCClientID,
   CAContentProvider:    oidcCAContent,
   UsernameClaim:        config.OIDCUsernameClaim,
   UsernamePrefix:       config.OIDCUsernamePrefix,
   GroupsClaim:          config.OIDCGroupsClaim,
   GroupsPrefix:         config.OIDCGroupsPrefix,
   SupportedSigningAlgs: config.OIDCSigningAlgs,
   RequiredClaims:       config.OIDCRequiredClaims,
  })
  if err != nil {
   return nil, nil, err
  }
  tokenAuthenticators = append(tokenAuthenticators, authenticator.WrapAudienceAgnosticToken(config.APIAudiences, oidcAuth))
 }
 // WebhookTokenAuth 认证器
 if len(config.WebhookTokenAuthnConfigFile) > 0 {
  webhookTokenAuth, err := newWebhookTokenAuthenticator(config)
  if err != nil {
   return nil, nil, err
  }

  tokenAuthenticators = append(tokenAuthenticators, webhookTokenAuth)
 }

 if len(tokenAuthenticators) > 0 {
  // 将 Token 相关的认证器组合成一个 BearerToken 认证器
  tokenAuth := tokenunion.New(tokenAuthenticators...)
  // Optionally cache authentication results
  if config.TokenSuccessCacheTTL > 0 || config.TokenFailureCacheTTL > 0 {
   tokenAuth = tokencache.New(tokenAuth, true, config.TokenSuccessCacheTTL, config.TokenFailureCacheTTL)
  }
  authenticators = append(authenticators, bearertoken.New(tokenAuth), websocket.NewProtocolAuthenticator(tokenAuth))
  securityDefinitions["BearerToken"] = &spec.SecurityScheme{
   SecuritySchemeProps: spec.SecuritySchemeProps{
    Type:        "apiKey",
    Name:        "authorization",
    In:          "header",
    Description: "Bearer Token authentication",
   },
  }
 }

 if len(authenticators) == 0 {
  if config.Anonymous {
   // Anonymous 认证器，即匿名认证
   return anonymous.NewAuthenticator(), &securityDefinitions, nil
  }
  return nil, &securityDefinitions, nil
 }

 // 将上诉所有认证器组合成一个认证器
 authenticator := union.New(authenticators...)

 authenticator = group.NewAuthenticatedGroupAdder(authenticator)

 if config.Anonymous {
  // If the authenticator chain returns an error, return an error (don't consider a bad bearer token
  // or invalid username/password combination anonymous).
  authenticator = union.NewFailOnError(authenticator, anonymous.NewAuthenticator())
 }

 return authenticator, &securityDefinitions, nil
}
```

kube-apiserver 会注册一系列的认证器，包括 `RequestHeader` 、`ClientCA`、`BearerToken`（`TokenAuth`、`ServiceAccountAuth`、`BootstrapToken`、`OIDC`、`WebhookTokenAuth`），认证器会实现 `authenticator.Request` 接口，其中的 Token 认证器则实现 `authenticator.Token` 接口：

```
// k8s.io/apiserver/pkg/authentication/authenticator/interfaces.go

// Token 认证器
type Token interface {
 AuthenticateToken(ctx context.Context, token string) (*Response, bool, error)
}

// 认证器
type Request interface {
 AuthenticateRequest(req *http.Request) (*Response, bool, error)
}
```

所有认证器初始化完成后，使用 `union.New` 方法组合成一个 `unionAuth` 认证器，只要请求 API 时满足组合中的任何一个认证器，则认证成功，实现方法很简单：

```
// k8s.io/apiserver/pkg/authentication/request/union/union.go

type unionAuthRequestHandler struct {
 Handlers []authenticator.Request
 FailOnError bool
}

func New(authRequestHandlers ...authenticator.Request) authenticator.Request {
 if len(authRequestHandlers) == 1 {
  return authRequestHandlers[0]
 }
 return &unionAuthRequestHandler{Handlers: authRequestHandlers, FailOnError: false}
}

func NewFailOnError(authRequestHandlers ...authenticator.Request) authenticator.Request {
 if len(authRequestHandlers) == 1 {
  return authRequestHandlers[0]
 }
 return &unionAuthRequestHandler{Handlers: authRequestHandlers, FailOnError: true}
}

// 认证方法
func (authHandler *unionAuthRequestHandler) AuthenticateRequest(req *http.Request) (*authenticator.Response, bool, error) {
 var errlist []error
 // 遍历所有的认证器
 for _, currAuthRequestHandler := range authHandler.Handlers {
  // 认证
  resp, ok, err := currAuthRequestHandler.AuthenticateRequest(req)
  if err != nil {
   if authHandler.FailOnError {
    // 一旦发生错误立即返回
    return resp, ok, err
   }
   // 继续尝试使用下一个认证器进行认证
   errlist = append(errlist, err)
   continue
  }

  // 认证成功
  if ok {
   return resp, ok, err
  }
 }

 // 所有认证器都认证不通过，返回认证失败
 return nil, false, utilerrors.NewAggregate(errlist)
}
```

在组合成一个认证器后，还会再调用 `group.NewAuthenticatedGroupAdder` 方法，在 `unionAuth` 认证器外面再包装一个 `AuthenticatedGroup` 认证器：

```
// AuthenticatedGroupAdder adds system:authenticated group when appropriate
type AuthenticatedGroupAdder struct {
 // Authenticator is delegated to make the authentication decision
 Authenticator authenticator.Request
}

// NewAuthenticatedGroupAdder wraps a request authenticator, and adds the system:authenticated group when appropriate.
// Authentication must succeed, the user must not be system:anonymous, the groups system:authenticated or system:unauthenticated must
// not be present
func NewAuthenticatedGroupAdder(auth authenticator.Request) authenticator.Request {
 return &AuthenticatedGroupAdder{auth}
}

func (g *AuthenticatedGroupAdder) AuthenticateRequest(req *http.Request) (*authenticator.Response, bool, error) {
 // 实际调用的就是 unionAuth 认证器的认证方法
 r, ok, err := g.Authenticator.AuthenticateRequest(req)
 if err != nil || !ok {
  // unionAuth 认证器认证失败
  return nil, ok, err
 }

 // unionAuth 认证器认证成功

 // 如果用户是 system:anonymous ，则直接返回
 if r.User.GetName() == user.Anonymous {
  return r, true, nil
 }
 for _, group := range r.User.GetGroups() {
  // 如果用户组属于 system:authenticated 或 system:unauthenticated ，则直接返回
  if group == user.AllAuthenticated || group == user.AllUnauthenticated {
   return r, true, nil
  }
 }

 // 否则更新用户信息，将用户加入 system:authenticated 组
 newGroups := make([]string, 0, len(r.User.GetGroups())+1)
 newGroups = append(newGroups, r.User.GetGroups()...)
 newGroups = append(newGroups, user.AllAuthenticated)

 ret := *r // shallow copy
 ret.User = &user.DefaultInfo{
  Name:   r.User.GetName(),
  UID:    r.User.GetUID(),
  Groups: newGroups,
  Extra:  r.User.GetExtra(),
 }
 return &ret, true, nil
}
```

到这，认证器的流程就捋清了，至于组合中的一系列认证器的内部处理细节不是很复杂，自己去看就行。

接着看鉴权器的初始化，来到 `BuildAuthorizer` 方法：

```
// cmd/kube-apiserver/app/server.go

func BuildAuthorizer(s *options.ServerRunOptions, EgressSelector *egressselector.EgressSelector, versionedInformers clientgoinformers.SharedInformerFactory) (authorizer.Authorizer, authorizer.RuleResolver, error) {
 // 鉴权器的配置，s 是 kube-apiserver 的启动参数选项
 authorizationConfig := s.Authorization.ToAuthorizationConfig(versionedInformers)

 if EgressSelector != nil {
  egressDialer, err := EgressSelector.Lookup(egressselector.ControlPlane.AsNetworkContext())
  if err != nil {
   return nil, nil, err
  }
  authorizationConfig.CustomDial = egressDialer
 }

 // 初始化鉴权器
 return authorizationConfig.New()
}

//pkg/kubeapiserver/options/authorization.go

func (o *BuiltInAuthorizationOptions) ToAuthorizationConfig(versionedInformerFactory versionedinformers.SharedInformerFactory) authorizer.Config {
 return authorizer.Config{
  // 鉴权模式
  AuthorizationModes:          o.Modes,
  PolicyFile:                  o.PolicyFile,
  WebhookConfigFile:           o.WebhookConfigFile,
  WebhookVersion:              o.WebhookVersion,
  WebhookCacheAuthorizedTTL:   o.WebhookCacheAuthorizedTTL,
  WebhookCacheUnauthorizedTTL: o.WebhookCacheUnauthorizedTTL,
  VersionedInformerFactory:    versionedInformerFactory,
  WebhookRetryBackoff:         o.WebhookRetryBackoff,
 }
}
```

默认情况下，鉴权器的配置是在 kube-apiserver 启动前，即第 1 回的初始化默认启动参数时的 `NewServerRunOptions` 方法设置的：

```
// cmd/kube-apiserver/app/options/options.go
func NewServerRunOptions() *ServerRunOptions {
 s := ServerRunOptions{
  // 准入控制器的配置，略过
  Admission:               kubeoptions.NewAdmissionOptions(),
  // 认证器的配置，略过
  Authentication:          kubeoptions.NewBuiltInAuthenticationOptions().WithAll(),
  // 鉴权器的配置
  Authorization:           kubeoptions.NewBuiltInAuthorizationOptions(),
  // ...
 }
 // ...
 return &s
}

// pkg/kubeapiserver/options/authorization.go
func NewBuiltInAuthorizationOptions() *BuiltInAuthorizationOptions {
 return &BuiltInAuthorizationOptions{
  // 默认的鉴权模式是 AlwaysAllow
  Modes:                       []string{authzmodes.ModeAlwaysAllow},
  WebhookVersion:              "v1beta1",
  WebhookCacheAuthorizedTTL:   5 * time.Minute,
  WebhookCacheUnauthorizedTTL: 30 * time.Second,
  WebhookRetryBackoff:         genericoptions.DefaultAuthWebhookRetryBackoff(),
 }
}

func (o *BuiltInAuthorizationOptions) AddFlags(fs *pflag.FlagSet) {
 // 可以通过 authorization-mode 启动参数来配置鉴权模式
 fs.StringSliceVar(&o.Modes, "authorization-mode", o.Modes, ""+
  "Ordered list of plug-ins to do authorization on secure port. Comma-delimited list of: "+
  strings.Join(authzmodes.AuthorizationModeChoices, ",")+".")
 // ...
}
```

看完了鉴权器配置，继续看其初始化 `authorizationConfig.New()` 方法：

```
// pkg/kubeapiserver/authorizer/config.go

func (config Config) New() (authorizer.Authorizer, authorizer.RuleResolver, error) {
 if len(config.AuthorizationModes) == 0 {
  return nil, nil, fmt.Errorf("at least one authorization mode must be passed")
 }

 var (
  // 用来存放一系列的鉴权器
  authorizers   []authorizer.Authorizer
  // 用来存放一系列的规则解析器
  ruleResolvers []authorizer.RuleResolver
 )

 // 很简单的一个 superuserAuthorizer 鉴权器，给 system:masters 用户组放行
 superuserAuthorizer := authorizerfactory.NewPrivilegedGroups(user.SystemPrivilegedGroup)
 authorizers = append(authorizers, superuserAuthorizer)

 // 遍历所有支持的鉴权模式
 for _, authorizationMode := range config.AuthorizationModes {
  // Keep cases in sync with constant list in k8s.io/kubernetes/pkg/kubeapiserver/authorizer/modes/modes.go.
  switch authorizationMode {
  // Node 模式
  case modes.ModeNode:
   node.RegisterMetrics()
   graph := node.NewGraph()
   node.AddGraphEventHandlers(
    graph,
    config.VersionedInformerFactory.Core().V1().Nodes(),
    config.VersionedInformerFactory.Core().V1().Pods(),
    config.VersionedInformerFactory.Core().V1().PersistentVolumes(),
    config.VersionedInformerFactory.Storage().V1().VolumeAttachments(),
   )
   nodeAuthorizer := node.NewAuthorizer(graph, nodeidentifier.NewDefaultNodeIdentifier(), bootstrappolicy.NodeRules())
   authorizers = append(authorizers, nodeAuthorizer)
   ruleResolvers = append(ruleResolvers, nodeAuthorizer)
  // AlwaysAllow 模式，也是默认的模式
  case modes.ModeAlwaysAllow:
   alwaysAllowAuthorizer := authorizerfactory.NewAlwaysAllowAuthorizer()
   authorizers = append(authorizers, alwaysAllowAuthorizer)
   ruleResolvers = append(ruleResolvers, alwaysAllowAuthorizer)
  // AlwaysDeny 模式
  case modes.ModeAlwaysDeny:
   alwaysDenyAuthorizer := authorizerfactory.NewAlwaysDenyAuthorizer()
   authorizers = append(authorizers, alwaysDenyAuthorizer)
   ruleResolvers = append(ruleResolvers, alwaysDenyAuthorizer)
  // ABAC 模式
  case modes.ModeABAC:
   abacAuthorizer, err := abac.NewFromFile(config.PolicyFile)
   if err != nil {
    return nil, nil, err
   }
   authorizers = append(authorizers, abacAuthorizer)
   ruleResolvers = append(ruleResolvers, abacAuthorizer)
  // Webhook 模式
  case modes.ModeWebhook:
   if config.WebhookRetryBackoff == nil {
    return nil, nil, errors.New("retry backoff parameters for authorization webhook has not been specified")
   }
   clientConfig, err := webhookutil.LoadKubeconfig(config.WebhookConfigFile, config.CustomDial)
   if err != nil {
    return nil, nil, err
   }
   webhookAuthorizer, err := webhook.New(clientConfig,
    config.WebhookVersion,
    config.WebhookCacheAuthorizedTTL,
    config.WebhookCacheUnauthorizedTTL,
    *config.WebhookRetryBackoff,
   )
   if err != nil {
    return nil, nil, err
   }
   authorizers = append(authorizers, webhookAuthorizer)
   ruleResolvers = append(ruleResolvers, webhookAuthorizer)
  // RBAC 模式
  case modes.ModeRBAC:
   rbacAuthorizer := rbac.New(
    &rbac.RoleGetter{Lister: config.VersionedInformerFactory.Rbac().V1().Roles().Lister()},
    &rbac.RoleBindingLister{Lister: config.VersionedInformerFactory.Rbac().V1().RoleBindings().Lister()},
    &rbac.ClusterRoleGetter{Lister: config.VersionedInformerFactory.Rbac().V1().ClusterRoles().Lister()},
    &rbac.ClusterRoleBindingLister{Lister: config.VersionedInformerFactory.Rbac().V1().ClusterRoleBindings().Lister()},
   )
   authorizers = append(authorizers, rbacAuthorizer)
   ruleResolvers = append(ruleResolvers, rbacAuthorizer)
  default:
   return nil, nil, fmt.Errorf("unknown authorization mode %s specified", authorizationMode)
  }
 }

 // 同样地，最后会将所有鉴权器、规则解析器组合成一个，组合方法类似认证器，不再展开
 return union.New(authorizers...), union.NewRuleResolvers(ruleResolvers...), nil
}
```

大致流程和认证器一样，鉴权器会实现一个 `authorizer.Authorizer` 接口，而规则解析器则实现 `authorizer.RuleResolver` 接口：

```
// k8s.io/apiserver/pkg/authorization/authorizer/interfaces.go

// 鉴权器
type Authorizer interface {
 Authorize(ctx context.Context, a Attributes) (authorized Decision, reason string, err error)
}

// 规则解析器
type RuleResolver interface {
 RulesFor(user user.Info, namespace string) ([]ResourceRuleInfo, []NonResourceRuleInfo, bool, error)
}
```

鉴权器默认仅支持 `AlwaysAllow` 模式，但可以使用 `--authorization-mode` 启动参数来扩展支持 `AlwaysDeny` 、`ABAC` 、`Webhook` 、`RBAC` 、`Node` 模式，初始化时会将所有支持的模式添加到鉴权器列表中，最后使用 `union.New` 方法合并成一个鉴权器。

以默认的 `AlwaysAllow` 模式为例，其实现方法：

```
// k8s.io/apiserver/pkg/authorization/authorizerfactory/builtin.go

type alwaysAllowAuthorizer struct{}

// 鉴权方法
func (alwaysAllowAuthorizer) Authorize(ctx context.Context, a authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
 // 直接返回通过
 return authorizer.DecisionAllow, "", nil
}

// 鉴权规则解析
func (alwaysAllowAuthorizer) RulesFor(user user.Info, namespace string) ([]authorizer.ResourceRuleInfo, []authorizer.NonResourceRuleInfo, bool, error) {
 // 支持所有 API
 return []authorizer.ResourceRuleInfo{
   &authorizer.DefaultResourceRuleInfo{
    Verbs:     []string{"*"},
    APIGroups: []string{"*"},
    Resources: []string{"*"},
   },
  }, []authorizer.NonResourceRuleInfo{
   &authorizer.DefaultNonResourceRuleInfo{
    Verbs:           []string{"*"},
    NonResourceURLs: []string{"*"},
   },
  }, false, nil
}

func NewAlwaysAllowAuthorizer() *alwaysAllowAuthorizer {
 return new(alwaysAllowAuthorizer)
}
```

其它的不再一一贴出，举一反三即可。

下一回继续看准入控制。