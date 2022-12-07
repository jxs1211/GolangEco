# Kubernetes 证书管理系列（一）

Original 张晋涛 [MoeLove](javascript:void(0);) *2022-12-07 08:21* *Posted on 北京*

收录于合集

\#Kubernetes59个

\#证书管理1个

大家好，我是张晋涛。

![img](http://mmbiz.qpic.cn/mmbiz_png/uO0QratgttoHxia3WgquTOUTibXA2nf2nDXJeVpR8A9G1S1tyTZeH5NpicsOZGHDVg1T6tQKNQng70ExtSg0pnibZw/0?wx_fmt=png)

**MoeLove**

不只限于 Container, Docker, Kubernetes 等技术，与你分享更多实用且具有前景的技术。欢迎关注

206篇原创内容



公众号

这是一个系列文章，将会通过七篇内容和大家一起聊聊 Kubernetes 中的证书管理。

以下是内容概览：

![Image](https://mmbiz.qpic.cn/mmbiz_png/uO0QratgttppSpAtesucLbFDuTIZRsOcEoID7AFNCpiazzjElgdCdLoSFUbtMwhk15EKiauaV0OwD4yjoaxz7zYw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

img

如上所示，在第一篇中，我们将从原理出发，来理解 Kubernetes 中的证书及其相关的作用，然后从需求的角度来理解 Kubernetes 证书管理器在实际生产中所起的作用。然后以 cert-manager 为例，结合架构、组件、生态兼容等来一点一点的拆解学习 cert-manager 。

自第二篇起，集中在了实战部分，毕竟有了理论知识也还是要付诸实践的。实践是检验真理的唯一标准。

内容较长，小伙伴们可自行斟酌跳转到感兴趣的部分，我也希望在评论区收到小伙伴们的反馈，谢谢。

# 证书

我们先来聊聊证书的一些基础知识。

## X.509 证书

### 1. 简介

在我们谈到 Kubernetes 证书时，我们指的是数字证书 Digital Certificate ，也就是基于 X.509 V3 标准的证书。

它通过将身份与一对可用于对数字信息进行加密、签名和解密的电子密钥绑定，以实现认证和数据安全（一致性、保密性）的保障。

每一个 X.509 证书都是根据公钥和私钥组成的密钥对来构建的，它们能够用于加解密、身份验证、信息安全性确认。

> “
>
> https://datatracker.ietf.org/doc/html/rfc5280
>
> ”

> “
>
> （备注：【 X.509 V3】--  X.509标准到目前为止有三个版本。第一个版本在1988年发布，名为X.500；到1996年为了顺应网络发展要求而扩展，定义了 V2 版本；2019年定义了现在的标准版本，即X.509 V3 标准。）
>
> ”

#### 1.1 加解密信息

X.509 标准使用了一种抽象语法表示法 One (ASN.1)的接口描述语言，来定义、和编解码客户端与证书颁发机构之间传输的证书请求和证书。

以下便是使用 ASN.1 的证书表示语法。

```
SignedContent ::= SEQUENCE 
{
  certificate         CertificateToBeSigned,
  algorithm           Object Identifier,
  signature           BITSTRING
}

CertificateToBeSigned ::= SEQUENCE 
{
  version                 [0] CertificateVersion DEFAULT v1,
  serialNumber            CertificateSerialNumber,
  signature               AlgorithmIdentifier,
  issuer                  Name
  validity                Validity,
  subject                 Name
  subjectPublicKeyInfo    SubjectPublicKeyInfo,
  issuerUniqueIdentifier  [1] IMPLICIT UniqueIdentifier OPTIONAL,
  subjectUniqueIdentifier [2] IMPLICIT UniqueIdentifier OPTIONAL,
  extensions              [3] Extensions OPTIONAL
}
```

公钥和私钥都是由一长串随机数组成的。公钥是公开的，由长度来决定保护强度，但是信息会通过公钥来加密。私钥只在接受者处秘密存储，接受者通过使用公钥关联的私钥才能解密读取信息。

用于生成公钥的最常见加密算法有以下三种：

- Rivest–Shamir–Adleman (RSA)  ：RSA来自 Ron Rivest、Adi Shamir 和 Leonard Adleman 这三个人的姓氏，他们于 1977 年公开描述了“ RSA”算法。RSA 根据两个大质数和一个辅助值创建并发布公钥。质数是保密的。消息可以由任何人通过公钥加密，但只能由知道素数的人解码。
- 椭圆曲线密码学 (ECC)：Elliptic Curve Cryptography，是一种基于椭圆曲线数学的公开密钥加密算法。它可以在保障跟 RSA 同等的安全级别下，使用的字符串长度更小。它还可以基于 Weil pairing 或者 Tate pairing 定义群之间的双线性映射。
- 数字签名算法 (DSA) ：1991年由美国国家标准技术研究所（NIST）提出，有专利，并且 NIST 在实行买断式授权。

#### 1.2 证书的编码

证书内容的编码（即，文件中存储的内容编码）在 X.509 标准中还没有被界定下来。

目前常见编码模式有两种：

- 可分辨编码规则(DER)：二进制格式，最常见，因为DER能处理大部分数据。DER编码的证书是二进制文件，文本编辑器无法直接读取。聊到这里，你可能会好奇，那么为什么叫它“可分辨编码”呢？这是由于它的编码方式都是公开的，就比如前面我给出的示例中，有个 `SEQUENCE`，它通过 DER 编码的时候标记编号是 `0x30` ，这样就可以很容易分辨出来了；
- 隐私增强邮件(PEM)：ASCII 文本格式，这是一种加密的电子邮件编码规则，可将DER编码的证书转换为文本文件。用 OpenSSL 工具的话，很简单 `openssl x509 -in <filename>.cer -inform DER -out <filename>.pem -outform PEM` 就可以了。

#### 1.3 Kubernetes 集群中的证书

在 `/etc/kubernetes/pki `目录中，我们可以看到一些文件。文件的结尾符，以及说明如下图。



img

使用 openssl 以文本模式打开一份 apiserver.crt 公钥证书：



img

说明：

- Serial Number 序列号，它是一个大端数字，也是唯一标识符。
- Subject  CN = kube-apiserver 是主体信息 ，（可以用来表示证书持有者的信息：国籍(C:Country 两个字母表示，中国CN), 省级行政区(ST: State or Province)，城市(L:Locality), 公司组织名(O:Organization), 组织部门名(OU: Organization Unit), 通用名(CN: Common Name), email地址等等）
- Validity 有效时间，结尾有时区标注 （GMT）
- Subject Public Key Info 是公钥信息和公钥算法
- 【Extension 扩展信息部分】- Key Usage 证书中公钥的用途，（大概有9种，digitalSignature、nonRepudiation、keyEncipherment、dataEncipherment、keyAgreement、keyCertSign、cRLSign、encipherOnly、decipherOnly。通常都是用来做 TLS 中的数字签名的）



img

### 2. 证书的颁发及信任链

Certification Authority 简称 CA，它是证书的认证权威机构。它的体系是一个树形结构，每一个 CA 可以有一到多个子 CA ，顶层 CA 被称为根 CA 。除了根 CA 以外，其他的 CA 证书的颁发者是它的上一级 CA 。这种层级关系组成了信任链（Trust Chain）。



img

以一个实际的例子来看，比如我的博客 `moelove.info` 我通过 cloudflare 提供了 SSL 证书，那么查看证书的时候，就可以看到它的根是 Baltimore CyberTrust Root，中间证书是 Cloudflare Inc ECC CA-3，最后则是 `sni.cloudflaressl.com```



img

证书分为两种类型（没有本质区别）：

- CA证书（CA Certificate）
- 终端实体证书（End Entity Certificate）  接受CA证书的最终实体。

### 3. CRL

X.509标准还定义了证书吊销列表（Certificate revocation list **，**CRL）的使用。CRL 是尚未到期就被证书颁发机构 CA 吊销的证书的名单。在 CRL 中的证书就不会受到信任了。

根据在 RFC 3280 中的定义，吊销有两种不同的状态：

- 吊销：不可逆。被吊销的最常见的原因是私钥泄露。
- 吊扣：是可逆的。

客户端应验证服务器提供的证书或用作 CA 的证书的序列号未出现在 CRL 中。理想情况下，每次验证证书时，都会根据当前版本的域 CRL 检查这些序列号。实际上，每个连接都拉取一遍更新的 CRL 会给流程增加很多开销。所以大多数客户端会使用缓存，缓存虽然不能保证是最新的，但可以避免在托管 CRL 的站点不可用时出现问题。

### 4. PKI

Public Key Infrastructure，简称 PKI，是公钥基础结构。PKI 的核心是在客户端、服务器和证书颁发机构 (CA) 之间建立的信任。这种信任是通过证书的生成、交换和验证来建立和传播的。

下图例举出了 Authentication 和 Certification 的差别，（双方和三方的区别）。



img

传输层安全协议 (Transport Layer Security ，TLS)  使用 PKI 证书来验证相互通信的各方以及加密通信会话。TLS 建立在 1990 年代后期的 SSL 标准之上，这里 TLS 连接会出现的常见错误大概有以下几种：

- 名称不匹配，证书的 CN 部分或信息存在不一致。
- 证书已过期，需要由 CA 重新颁发。
- 用于签署证书的根 CA 不在客户端的受信任密钥库中。

## K8S 基于CA 签名的双向数字证书



img

在 Kubernetes 中，各个组件提供的接口中包含了集群的内部信息。出于对集群安全的考虑，需要组件之间的通信需要采用双向 TLS 认证，以防出现接口被非法访问的情况。（即客户端和服务器端都需要验证对方的身份信息）

在 https://medium.com/littlemanco/the-magic-of-tls-x509-and-mutual-authentication-explained-b2162dec4401 中有介绍过双向 TLS 的必要性，感兴趣的小伙伴可以详细看下。简单的说，就是 A 向 B 发出信息前，需要确认 B 是个真的 B；而 B 对 A 作出响应，也需要确认 A 是真的 A 。



img

K8S 集群中的证书主要包含如下部分：

- Kubelet 客户端证书，用于 API 服务器身份验证
- Kubelet 服务端证书， 用于 API 服务器与 Kubelet 的会话
- API 服务器端证书
- 集群管理员的客户端证书，用于 API 服务器身份认证
- API 服务器的客户端证书，用于和 Kubelet 的会话，还有和 etcd 的会话
- kube-controller-manager 的客户端证书，用于和 API 服务器的会话
- kube-scheduler 的客户端证书或 kubeconfig，用于和 API 服务器的会话
- 前端代理的客户端及服务端证书
- etcd 相关，用于客户端和其他对等节点进行身份验证。

通过 kubeadm 安装 Kubernetes，大多数证书都存储在 `/etc/kubernetes/pki`。 本文档中的所有路径都是相对于该目录的，但用户账户证书除外，kubeadm 将其放在 `/etc/kubernetes` 中。

# K8S中的证书管理

聊完证书的主要内容，我们来看看 Kubernetes 中的证书管理。

## 需求

一个 K8S 集群中，可以存在多个以独立进程形式存在的组件。这些组件通过相互通信来实现，集群的运行、管理等工作。上文也已解释过，对于各个需要通信的组件来讲，关键的是需要验证通信双方的身份是否符合预期，以免受到安全威胁。



img

在两个组件需要进行双向认证时，就涉及到了一系列证书（参照上一部分罗列的证书文件）。

对于这些证书的管理涉及到了从生成、分发、续签等等实际生产中需要用到的方方面面。



img

当我们使用 Kubernetes 相关的组件和工具时，可以从下方直观的看到，涉及的命令参数繁多。



img

官方文档中罗列了一些具体的参数和作用，这里就不再展开了。



img

综上，K8S的证书管理过程中，我们需要面对一些问题：

- 证书种类繁多
- 证书数量繁多
- 证书管理命令行工具参数繁多，存在一定的使用复杂度
- 管理过程及检验过程不直观
- 人为操作风险较高

至此，一些 Kubernetes 中的证书管理器的项目应运而生。

## K8S 集群的证书管理项目

在这一部分，我们主要会介绍 cert-manager 的发展历程。

为什么要选择介绍 cert-manager 呢？这是因为这个项目几乎已经成为了 Kubernetes 中证书管理领域的事实标准了。



img

从2016年至今，证书管理器经历了支持+扩展的历程。如今，我们可以很便捷的使用 cert-manager 来实现 Kubernetes 集群中的证书管理工作啦。

当然，这里有必要说明一下，cert-manager 所管理的证书，主要是为部署在 Kubernetes 中的服务所使用的，而非给 Kubernetes 自身。通常我们可以使用 `kubeadm certs` 命令来完成大多数通过 kubeadm 部署的 Kubernetes 集群的证书相关操作。

# cert-manager

## 简介



img

cert-manager 将证书和证书颁发者作为自定义资源类型添加到 Kubernetes 集群中，并简化了这些证书的获取、更新和使用过程。



img

cert-manager 可以从各种受支持的来源颁发证书，包括 Let's Encrypt、HashiCorp Vault和Venafi以及私有 PKI。

- Let's Encrypt 是全球证书颁发机构 (CA)，可以获取、更新和管理 SSL/TLS 证书。网站可以使用 Let's Encrypt 的证书来启用安全的 HTTPS 连接。Let's Encrypt 是一个非营利组织，对于大多数浏览器和操作系统都信任来自 Let's Encrypt 的证书。在不停机的情况下自动/手动颁发和安装证书的话，需要借助 Certbot （一款免费的开源软件工具）来启用 HTTPS。
- HashiCorp Vault 是一个开源的密钥和隐私数据管理工具。它提供了非常丰富的功能，除了常规的密钥存储外，还支持证书签发，策略管理等能力。
- Venafi 机器身份管理平台。

## 部署使用流程

### 1. 使用 ACME 签发的部署流程



img

上图详细的描述了 使用 ACME 签发证书的 cert-manager 的部署安装流程。

### 2. 使用 Vault 签发的部署流程



img

上图详细的描述了 使用 Vault 签发证书的 cert-manager 的部署安装流程。

以上内容，我会在后续文章中进行详细展开，敬请期待。

## Others

### 1. 证书资源

使用证书资源是请求签名证书的最简单和最常用的方法。证书资源将 certificate request 展示成了一种可读的CRD。打开 cert-manager 的 yaml ，在`certificate.spec.issuerRef` 部分，即指明从哪个颁发者处获取证书（默认类型是 Issuer，也可以通过更改类型指定成 ClusterIssuer）。

> “
>
> 注意：如果你想创建的 Issuer 可以被所有 Certificate namespaces 中的资源引用，那么需要将`certificate.spec.issuerRef.kind`字段设置为 ClusterIssuer 。
>
> ”

`certificate.spec.secretTemplate` 部分是可选部分，这里需要说的是，使用 `renewBefore`和`duration` 字段来控制证书持续时间以及更新时，需要安装 webhook 组件，并且需要注意字符串格式（用 s，m，h）。

webhook 组件部署为一个独立的 pod 来进行：

- 验证：确保在创建或更新 cert-manager 资源时，合规（ API 规则）。
- 变更/恢复默认：在创建和更新操作期间更改资源的内容。
- 转换：主要为了应对多版本的 API 适配。

默认情况下，私钥不会自动轮换。但是，如果 `certificate.spec.privateKey.rotationPolicy`设置成 Always ，则与证书对象关联的私钥 Secret 可以配置为在操作触发证书对象的重新发布时立即轮换。

cert-manager 会等到 Certificate 对象被正确签名后再覆盖`tls.key`Secret 中的文件，整个过程热操作不会中断。

重新颁发证书对象的方式：

- `status.renewalTime`到期，即 X.509 证书即将到期时。
- 证书中以下字段更改时：`commonName`, `dnsNames`, `ipAddresses`, `uris`, `emailAddresses`, `subject`, `isCA`, `usages`,`duration`或`issuerRef`。
- 手动执行：cmctl renew cert 。

如果 `certificate.spec.privateKey.rotationPolicy`设置成 Never，仅在启动时加载一次私钥和签名证书，如果遇到更新轮换时，需要手动重启 pod ：`kubectl rollout restart` ，或通过运行一个 Secret 控制器 ： wave 来进行自动操作。

cert-manager 在删除 Certificate 资源时，默认不会删除对应的 Secret 资源。若希望同步删除的方式：

- 配置参数 --enable-certificate-owner-ref。
- 手动删除。

cert-manager 可以自动续订证书，它根据颁发证书的持续时间`spec.duration` 字段和到期前多久进行更新`spec.renewBefore`字段来进行计算确认更新时间。

注：

- `spec.duration`默认值为90 天。
- `spec.duration`最小值为1 小时；`spec.renewBefore`最小值为 5 分钟。
- `spec.duration`> `spec.renewBefore`。

### 2. Prometheus Metrics

这部分需要结合部署模式来进行：

- 使用 Helm 部署 cert-manager：ServiceMonitor 可进行配置启用。
- 使用 YAML manifests ：首先，需要在 cert-manager 的部署 yaml 中追加 containers 部分；然后，创建一个 PodMonitor。

后续文章中会有如何进行完整的监控实践，敬请期待。

### 3. 与 Ingress 集成

cert-manager 中的一个组件 ingress-shim 负责通过向 Ingress 资源添加注释来实现请求 TLS 签名证书来保护 Ingress 资源。ingress-shim 会监控集群中的 Ingress 资源，如果 Ingress 的 annotation 中出现了 cert-manager 所支持的部分，则会对应的创建证书资源。

后续文章中也会具体介绍。

### 4. 与GateWay API 集成

Gateway API 是由 Kubernetes 社区的 SIG-Network 管理的开源项目，是一组安装在 Kubernetes 集群上的 CRD。主要目标是替换 Ingress，号称下一代 Ingress 。

GateWay API 有三种主要类型的对象：

- GatewayClass 定义一组具有通用配置和行为的网关。它是集群范围的资源。通常由基础设施供应商进行操作和管理。

- Gateway 定义了如何将流量转发到集群内服务。

- Routes 将来自 Gateway 的流量映射到服务。它定义了将请求从网关映射到 Kubernetes 服务的规则。从 v1alpha2 开始，API 包含四种 Route 资源类型：

- - [7层] HTTPRoute ：用于多路复用 HTTP 或进行 HTTPS 卸载。即对使用 HTTP 请求的数据进行路由或修改。
  - [4-7层间] TLSRoute ：用于多路复用 TLS 连接，通过 SNI 进行区分。
  - [4层] TCPRoute（和 UDPRoute）：将一个或多个端口映射到单个后端。
  - [7层] GRPCRoute ：用于路由 gRPC 流量。



img

cert-manager 可以通过向 Gateway 添加注释来实现为 Gateway 资源生成 TLS 证书。Gateway 资源是 Gateway API 的一部分。

具体启用操作也与安装模式有关。这个在后续文章也会详细讲解以及操作演示。需要注意的是，Gateway API CRD 应该在 cert-manager 启动之前安装，或者在安装 Gateway API CRD 之后重启 cert-manager。

### 5. 与 Istio 集成

istio-csr 是一个 agent，允许使用 cert-manager 保护 Istio workload 和控制平面组件 。istio-csr 实现了 gRPC Istio 证书服务，该服务对来自 Istio workload 的传入证书签名请求进行身份验证、授权和签名，并通过安装在集群中的 cert-manager 路由所有证书的处理。具体操作还是请期待后续文章。

下图，摘自 https://istio.io/latest/about/service-mesh/ 感兴趣的小伙伴可以详细了解下。



img

### 6. 与 SPIRE 集成



img

SPIRE （SPIRE Runtime Environment）是一个 API 工具链，用于在各种托管平台的软件系统之间建立信任。SPIRE 只是 SPIFFE 规范的一种实现。SPIRE 公开了 SPIFFE Workload API ，它可以保障正在运行的软件系统并向它们颁发 SPIFFE ID 和 SVID 。

- SPIFFE ID 是用于标识资源或调用者的结构化字符串，SPIFFE 组件致力于 SPIFFE ID 的发布及验证。
- SVID（SPIFFE 可验证身份文件）即，计算端点可以通过密码验证来信任或者拒绝 某一 SPIFFE ID 。SVID 可以引用相关联的非对称密钥对，还可以用于形成安全通信通道。

SPIRE 还可以使 workload 能够安全地向存储、数据库或云厂商进行身份验证。

一个 SPIRE 由一个 SPIRE Server 以及一个或者多个 SPIRE Agent 组成。如下图所示，



img

- SPIRE Server 负责管理和发布其配置的 SPIFFE 信任域中的所有 ID 。它还存储注册条目（选定选择器来明确应发布的 SPIFFE ID 的条件）和签名密钥，使用 Node API 自动验证 Agent 的身份，并在经过身份验证的 Agent 请求时为 Workload 创建 SVID。

- - Node attestor plugins 与 agent node attestors 一起验证 Agent 身份。
  - Node resolver plugins 扩展了 Server 可以用来通过验证节点的属性。
  - Datastore plugins 服务器使用它来存储、查询和更新数据。
  - Key manager plugins 存储用于签署 X.509-SVID 和 JWT-SVID 的私钥。
  - Upstream authority plugins 默认情况下，SPIRE Server 就是自己的 CA，也可以有其他选择。

- SPIRE Agent 从 Server 请求 SVID 并缓存，等待 Workload 请求其 SVID；将 SPIFFE Workload API 暴露给节点上的 Workload，并证明其身份，然后为已识别的 Workload 提供 SVID。

- - Node attestor plugins 与 server node attestors 一起验证 Agent 的身份。
  - Workload attestor plugins 验证节点上 workload 进程的身份。
  - Key manager plugins 为发布给 workload 的 X.509-SVID 生成和使用私钥。

那么， SPIRE 的工作流程及部署使用实战就请期待后续的文章吧。

# 总结

本篇从证书的基本概念聊起，结合了 Kubernetes 环境中对证书的应用，以及分析了在 Kubernetes 环境中证书管理的复杂度。之后分别介绍了 cert-manager 及其主要的应用场景，具体的操作实践将在后续文章中为大家展现。

欢迎大家在评论区留言讨论，也请点赞再看，谢谢。