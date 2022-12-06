https://mp.weixin.qq.com/s/lJTYWKrTtPp7rt658a251A

# Kubernetes 二次开发浅析

蔡奇 [k8s技术圈](javascript:void(0);) *2022-12-03 19:04* *Posted on 海南*



Kubernetes是一个开源的容器集群管理系统，可以实现容器集群的自动化部署、自动扩缩容、维护等功能。Kubernetes能满足绝大部分的开发运维需求，但是也存在需要基于Kubernetes进行二次开发实现特定业务逻辑的情况。本文将介绍kubernetes二次开发的相关知识，包含GVR、客户端工具、Informer机制、Code-generator代码生成等。

作者：蔡奇， 中国移动云能力中心软件开发工程师，专注于云原生、Istio、微服务、Spring Cloud 等领域。

##  

01

GVR与GVK



在Kubernetes体系中，资源是最重要的概念。Kubernetes使用Group、Version、Resource、Kind来描述。

 

![Image](https://mmbiz.qpic.cn/mmbiz_png/DmBLZYMe831r9VyOBLCWzsYq5mJvLERjWjuhxCShDZ5XbKkMJhhaV7gppWBjR1v5ZeofclIzpROc8zFvZ3B20w/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)



Group即资源组，在kubernetes中有两种Group：有组名资源组和无组名资源组(也叫核心资源组Core Groups)。如deployment有组名，pod没有组名。核心资源组API为api/v1，非核心资源组API为/apis/group/version。
Version即版本，kubernetes的版本分为三种：
1）Alpha：内部测试版本，如v1alpha1
2）Beta：经历了官方和社区测试的相对稳定版，如v1beta1
3）Stable：正式发布版，如v1、v2
在Kubernetes中，一个资源可能对应多个Version，也可能对应多个Group，因此通常使用GVK或GVR来区别特定的Kubernetes资源。二者有如下区别与联系：
1）GVR与HTTP请求里的PATH对应，如查询Pod的请求GET /api/v1/namespaces/{namespace}/pods就是一个GVR。GVK与存储在ETCD中的Object类型对应。
2）GVR与GVK通过REST映射可进行转化。
使用客户端工具如kubectl、clientSet、curl时，首先会根据GVR生成请求，然后Kubernetes API Server会查询HTTP PATH对应的Resource是否支持，并与ETCD进行交互。当API Server不支持该Resource时，Kubernetes会报错the server doesn't have a resource type "..."，使用kubectl api-resources命令可查看支持的Resource。



02

Client-go客户端工具





Client-go提供了四种不同的客户端工具对Kubernetes进行操作，分别是RestClient，ClientSet，DynamicClient、DiscoveryClient。

### 2.1. **RestClient**

RESTClient是client-go最基础的客户端，主要是对HTTP Reqeust进行了封装，对外提供RESTful风格的API，并且提供丰富的API用于各种设置，相比其他几种客户端虽然更复杂，但是也更为灵活。

使用RESTClient对kubernetes的资源进行增删改查的基本步骤如下：

1)确定要操作的资源类型(例如查找deployment列表)，去官方API文档中找到对于的path、数据结构等信息。

2)加载配置kubernetes配置文件。

3)根据配置文件生成配置对象，并且通过API对配置对象就行设置（例如请求的path、Group、Version、序列化反序列化工具等）。

4)创建RESTClient实例，入参是配置对象。

5)调用RESTClient实例的方法向kubernetes的API Server发起请求。

RestClient使用示例：

- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 

```
// 加载kubeconfigconfig, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfig)// 参考path： /api/v1/namespaces/{namespace}/podsconfig.APIPath = "api"config.GroupVersion = &corev1.SchemeGroupVersion// 指定序列化工具config.NegotiatedSerializer = scheme.CodecsrestClient, err := rest.RESTClientFor(config)if err != nil {    panic(err.Error())}// 保存pod结果result := &corev1.PodList{}err = restClient.Get().    Namespace("kube-system").    Resource("pods").    VersionedParams(&metav1.ListOptions{Limit: 50}, scheme.ParameterCodec).    Do(context.TODO()).    Into(result)
```



2.2. **ClientSet**

Clientset是所有Group、Version组成的客户端集合，每个GV客户端底层由RestClient实现。但ClientSet只能处理pod、deployment等事先已确定好GVR的Kubernetes内置资源。自定义的CR无法提前知道其GV信息和数据结构，无法使用ClientSet进行处理。

ClientSet使用示例：

- 
- 
- 
- 
- 
- 
- 
- 
- 
- 

```
// 加载kubeconfigconfig, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfig)if err != nil {    panic(err)}clientSet, err := kubernetes.NewForConfig(config)if err != nil {    panic(err)}deploymentList := clientSet.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
```



查看底层List方法，可以发现该方法由RestClient实现。

- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 

```
func (c *deployments) List(ctx context.Context, opts metav1.ListOptions) (result *v1.DeploymentList, err error) {var timeout time.Durationif opts.TimeoutSeconds != nil {    timeout = time.Duration(*opts.TimeoutSeconds) * time.Second}result = &v1.DeploymentList{}err = c.client.Get().    Namespace(c.ns).    Resource("deployments").    VersionedParams(&opts, scheme.ParameterCodec).    Timeout(timeout).    Do(ctx).    Into(result)    return}
```



2.3. **DynamicClient**

对于无法使用ClientSet的CR资源，client-go提供了DynamicClient进行处理。client-go使用Unstructured来表示数据结构不确定的资源数据，Unstructured由Interface实现。

- 
- 
- 

```
type Unstructured struct {    Object map[string]interface{}}
```

Unstructed和具体资源类型如Pod直接的转化由runtime.unstructuredConverter的FromUnstructured和ToUnstructured方法分别实现。将Unstructured转化为pod：

- 
- 

```
pod := &apiv1.Pod{}err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), pod)
```

将pod转化为Unstructured：

- 

```
unstructured,err := runtime.DefaultUnstructuredConverter.ToUnstructured(podList)
```

DynamicClient使用示例：

- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 

```
// 加载kubeconfigconfig, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfig)if err != nil {    panic(err)}dynamicClient, err := dynamic.NewForConfig(config)if err != nil {    panic(err.Error())}// dynamicClient关联方法所需入参gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}// 调用dynamicClient查询方法unstructObj, err := dynamicClient.Resource(gvr).Namespace("kube-system").List(context.TODO(), metav1.ListOptions{})// 实例化一个podlist，接收从unstructObj转换后的结果podList := &apiv1.PodList{}err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList)
```

### 2.4. **DiscoveryClient**

RestClientSet、Clientset和DynamicClient都是面向资源对象的(如创建Pod实例、查看server实例等)，DiscoveryClient聚焦资源自身，例如查看当前Kubernetes有哪些Group、Version、Resource。

DiscoveryClient使用示例：

- 
- 
- 
- 
- 
- 
- 
- 
- 
- 
- 

```
// 加载kubeconfigconfig, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfig)if err != nil {    panic(err)}discoverClient, err := discovery.NewDiscoveryClientForConfig(config)if err != nil {    panic(err)}// 获取所有分组和资源数据APIGroup, APIResourceList, err := discoverClient.ServerGroupsAndResources()
```



03

Informer



Kubernetes基于声明式API的设计理念，所谓声明式API，即告诉Kubernetes Controller资源对象的期望状态，这样为Kubernetes在事件通知后，动作执行前这段过程里提供了更多的容错空间与扩展空间。这就需要Kubernetes Controller能够知道资源对象的当前状态，通常需要访问API Server才能获得资源对象，当Controller越来越多时，会导致API Server负载过大。

Kubernetes使用Informer代替Controller去访问API Server，Controller的所有操作都和Informer进行交互，而Informer并不会每次都去访问API Server。Informer使用ListAndWatch的机制，在Informer首次启动时，会调用LIST API获取所有最新版本的资源对象，然后再通过WATCH API来监听这些对象的变化，并将事件信息维护在一个只读的缓存队列中提升查询的效率，同时降低API Server的负载。

除了ListAndWatch，Informer还可以注册相应的事件，之后如果监听到的事件变化就会调用对应的EventHandler，实现回调。Informer主要包含以下组件：

1)Controller：Informer的实施载体，可以创建reflector及控制processLoop。processLoop将DeltaFIFO队列中的数据pop出，首先调用Indexer进行缓存并建立索引，然后分发给processor进行处理。

2)Reflector：Informer并没有直接访问k8s-api-server，而是通过一个叫Reflector的对象进行api-server的访问。Reflector通过ListAndWatch监控指定的 kubernetes 资源，当资源发生变化的时候，例如发生了 Added 资源添加等事件，会将其资源对象存放在本地缓存 DeltaFIFO 中。

3) DeltaFIFO：是一个先进先出的缓存队列，用来存储 Watch API 返回的各种事件，如Added、Updated、Deleted 。

4)Indexer：Indexer使用一个线程安全的数据存储来存储对象和它们的键值。需要注意的是，Indexer中的数据与etcd中的数据是完全一致的，这样client-go需要数据时，无须每次都从api-server获取，从而减少了请求过多造成对api-server的压力。一句话总结：Indexer是用于存储+快速查找资源。

5)Processor：记录了所有的回调函数（即 ResourceEventHandler）的实例，并负责触发回调函数。

 

![Image](https://mmbiz.qpic.cn/mmbiz_png/DmBLZYMe831r9VyOBLCWzsYq5mJvLERjxuojj9KPEP0NK4Yr6hkZ0HwJfhP9OxugcmLMwpyN8TtNuV3ia77teLA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)



Informer工作流程如下：
1) 第一次启动Informer的时候，Reflector 会使用List从API Server主动获取资源对象信息，并更新DeltaFIFO中的items。
2) 持续使用Reflector建立长连接，去Watch API Server发来的资源对象变更事件。
3) Reflector监控到k8s资源对象有增加删除修改之后，就把资源对象变更事件信息存放在DeltaFIFO中。
4) DeltaFIFO是一个先进先出队列， Controller调用processLoop从队列中不断pop出事件信息, 首先将其存储至Indexer中，然后通过processor触发事件回调函数。
5) 回调函数将资源对象的key放进workqueue。
6) 通过用户在custom controller中自定义的worker（包含Process Item程序）处理workqueue中的item。

##  

04

Code-generator



Kubernetes支持创建CRD以根据业务定制所需要的资源类型，创建好的CRD对象会保存在ETCD中，但如果仅仅是在ETCD保存，那对象只是一条数据而已，没有什么实质性作用。因此，在创建CRD后，通常会实现client、informer、controller来操作、监听指定对象。

Code-generator 是Kubernetes提供的一个用于代码生成的项目，项目地址为https://github.com/kubernetes/code-generator，它提供了以下工具为 Kubernetes 中的资源生成代码。

1)deepcopy-gen: 生成深度拷贝方法，为每个T类型生成 func (t* T) DeepCopy() *T 方法，API 类型都需要实现深拷贝。

2) client-gen: 为资源生成标准的 clientset。

3) informer-gen: 生成 informer，提供事件机制来响应资源的事件。

4) lister-gen: 生成Lister，为get和list请求提供只读缓存层（通过 indexer 获取）。

在使用Code-generator之前，首先需要初始化doc.go,register.go,types.go三个文件。doc.go主要是用来声明要使用deepconpy-gen以及groupName。types.go主要是定义crd资源对应的go中的结构。register.go注册资源。

```
type.go// +genclient// +genclient:noStatus// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Objecttype Demo struct {    metav1.TypeMeta `json:",inline"`    metav1.ObjectMeta `json:"metadata,omitempty"`    Spec MydemoSpec `json:"spec"`}type MydemoSpec struct {    Uuid string `json:"uuid"`    Name string `json:"name"`}// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Objecttype DemoList struct {    metav1.TypeMeta `json:",inline"`    metav1.ListMeta `json:"metadata"`

    Items []Demo `json:"items"`}

doc.go// +k8s:deepcopy-gen=package// +groupName=cq.iopackage v1

register.govar SchemeGroupVersion = schema.GroupVersion{    Group:   "cq.io",    Version: "v1",}var (    SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)    AddToScheme   = SchemeBuilder.AddToScheme)func Resource(resource string) schema.GroupResource {    return SchemeGroupVersion.WithResource(resource).GroupResource()}func Kind(kind string) schema.GroupKind {    return SchemeGroupVersion.WithKind(kind).GroupKind()}func addKnownTypes(scheme *runtime.Scheme) error {    scheme.AddKnownTypes(        SchemeGroupVersion,        &Demo{},        &DemoList{},    )    metav1.AddToGroupVersion(scheme, SchemeGroupVersion)    return nil}
```

该段代码还包含一些tag，分为全局tag和局部tag。全局tag必须在doc.go文件中声明。+k8s:deepcopy-gen=package表示为整个包里的所有类型定义自动生成DeepCopy方法，+groupName=cq.io定义包对应的API 组的名字。

局部tag要么直接声明在类型之前，要么位于类型之前的第二个注释块中。通常在types.go文件中声明。+genclient表示为资源类型生成对应的 Client 代码。+genclient:noStatus 表示在生成的Client里，没有Status字段，实现spec-status分离。+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object 表示在生成DeepCopy 的时候，实现 Kubernetes 提供的 runtime.Object 接口。否则，在某些版本的 Kubernetes 里，类型定义会出现编译错误。

在定义好类似文件后，就可以使用code-generator来生成相应的代码。以开源项目kubespher为例，该项目使用以下命令生成pkg/client目录下clientset、Informer等代码。

```shell
./hack/generate_group.sh "client,lister,informer" kubesphere.io/kubesphere/pkg/client kubesphere.io/api "${GV}" --output-base=./  -h "$PWD/hack/boilerplate.go.txt"mv kubesphere.io/kubesphere/pkg/client ./pkg/

hack/generate_group.sh "deepcopy" kubesphere.io/api kubesphere.io/api ${GV} --output-base=staging/src/  -h "hack/boilerplate.go.txt"
```

Generator_group为code-generator源码中提供的代码生成脚本，通常会复制到自己项目的/hack目录下。第一个参数"client,lister,informer"，"deepcopy"代表以上4种标准代码生成器。当包含所有4种生成器时，可用all替代。第二个参数用于指定生成的client、informer、lister所在的包名。第三个参数指定API组的基础包名，deepcopy生成的zz_generated.deepcopy.go文件的包名与此参数相同。第四个参数指定API组的版本号，可以包含多个API信息，格式可参考"groupA:v1 groupB:v2"。第五个参数--output-base指定包输出的基础路径。第六个参数-h指定版权信息文件boilerplate.go.txt，通常也会复制到/hack目录下。

## **参考文献**

https://blog.csdn.net/boling_cavalry/article/details/113487087
https://zhuanlan.zhihu.com/p/391465614
https://andblog.cn/3196

