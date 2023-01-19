# 第1回：kube-apiserver 启动及前期调试准备

[CNCF](javascript:void(0);) *2023-06-06 10:10* *Posted on 香港*
[第1回：kube-apiserver 启动及前期调试准备 (qq.com)](https://mp.weixin.qq.com/s?__biz=Mzg2NTU3NjgxOA==&mid=2247488266&idx=1&sn=07226e3e82c90aeaf6f0d10782768a8e&chksm=ce56a117f921280155bb822887a4165d3a20bdb3cb789b8532c15af71dc8341c3e5025f03c73&cur_album_id=2958341226519298049&scene=190#rd)

The following article is from gopher云原生 Author 邹俊豪

[![img](http://wx.qlogo.cn/mmhead/Q3auHgzwzM62KB2Ce0NLjthrSTtH4JPKvYFRNkibGotXyOvNgIEIiavw/0)**gopher云原生**.技术log](https://mp.weixin.qq.com/s?__biz=MzI5ODk5ODI4Nw==&mid=2247533694&idx=2&sn=bbeabfacaef68b007e39c2dc1dda573c&chksm=ec9f4b1edbe8c208e64158a0e70f1a3fedef7073c19ce801066232b524efbcc91bc71fc23e6f&mpshare=1&scene=1&srcid=0606tjPjHggJAybi5xQTC8HH&sharer_sharetime=1686055932326&sharer_shareid=85cfd1be9a9cdea1f202dd1f395a2697&key=640c1bac93237266b64be71ed6cf7f3556d82b0bd5eca155c98d1a0c140b6b6fa38c9b1de35376c38b63bd169ad23294307ae5de25fb75bac6e76c64ca12ec8d3a27473980ba1c11cfb834169cff65a312e4f74e010eb7f3cbeb1d935c8b362e03ebf31932ac15316db8e6b36ef297218bc3fc75ae18c81eeb9bf5edf7eba9e4&ascene=1&uin=MjMxOTI2NTEwMA%3D%3D&devicetype=Windows+10+x64&version=63090549&lang=en&countrycode=CN&exportkey=n_ChQIAhIQoDh03wi7p0dE4MGWxNUZnBLYAQIE97dBBAEAAAAAAKHzL%2BWfE%2FIAAAAOpnltbLcz9gKNyK89dVj0ViW00xshCLGD6EspXRrgcsdj9L1ZZ1%2BE1Bmq99OiIWqkDVg4SVStdZVC2Gmfw1S7XkgMPObTIRSPVqt9ELbtPaPPe%2BLIrE3jqH8i4fgVCsvCQB9Pc7KdN0SIGjVQYGwzQInVdX%2B6RX6kb%2FrXF2ueUdzkltsYtNTQJVCoUkvukPca0HyuuY%2FAkGSlZKrjBzT0EMZhSpszJYWXzD%2FDydA%2F3yICIjYvq9askFLm1q37Wv7%2Fzw%3D%3D&acctmode=0&pass_ticket=XBCsG3Kp1Jx8IFM4ty9nPQf3GAkHB036w%2BIBEW9LCkLA86zig2ihCWZSiCnxiFIX&wx_header=1&fasttmpl_type=0&fasttmpl_fullversion=6713510-en_US-zip&fasttmpl_flag=3#)

## 前言

立个 flag，啃 k8s 源码。

版本：Kubernetes **1.27.2**[1]

代码结构：

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmeSOzNRNEjazLwjC02QfJPDB8d8yx92uKiayeL6FxLecS4hhLTicX8WsczZsa4gZ2nDNdrZgE5dufHQ/640?wx_fmt=png&wxfrom=13&tp=wxpic)

- `api`: 存放 OpenAPI 的 spec 文件
- `build`: 包含构建 Kubernetes 的工具和脚本
- `cluster`: 包含用于构建、测试和部署 Kubernetes 集群的工具和脚本
- `cmd`: 包含 Kubernetes 所有组件入口的源代码，例如 kube-apiserver、kube-scheduler、kube-controller-manager、kubelet、kube-proxy、kubectl 等
- `hack`: 包含用于构建和测试 Kubernetes 的脚本和工具
- `pkg`: 包含 Kubernetes 的核心公共库和工具代码
- `plugin`: 包含 Kubernetes 插件的源代码，例如认证插件、授权插件等
- `staging`: 存放部分核心库的暂存代码，这些暂存代码会定期发布到各自的顶级 **k8s.io**[2] 存储库
- `test`: 包含 Kubernetes 测试的源代码和测试工具
- `third_party`: 包含 Kubernetes 使用的第三方工具代码
- `vendor`: 包含 Kubernetes 使用的所有依赖库代码

Kubernetes 组件架构：

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmeSOzNRNEjazLwjC02QfJPDicyhIXbwicBmk70Ik81xt9PVGKkYa6LJqmrGvvAy5XibtxshuPXwicgKxQ/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

## 卷一：kube-apiserver

API 服务器是 Kubernetes 控制平面的组件， 该组件负责公开了 Kubernetes API，负责处理接受请求的工作。API 服务器可以看作是 Kubernetes 控制平面的前端。

### 第1回：apiserver 启动及前期调试准备

apiserver 的入口：

```
func main() {
 command := app.NewAPIServerCommand()
 code := cli.Run(command)
 os.Exit(code)
}
```

命令行工具的功能是使用的 **github.com/spf13/cobra**[3] 库，不要陷入细节，直接跳到 `app.NewAPIServerCommand()` 方法。

首先会初始化默认的启动参数，得到 `*ServerRunOptions` 类型的 `s` 变量：

```
s := options.NewServerRunOptions()
```

打个断点感受下 `s` 的默认值：

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

其中各个参数的介绍可以看：**https://kubernetes.io/zh-cn/docs/reference/command-line-tools-reference/kube-apiserver/**[4]

得到 `s` 参数后，还要对其进行一系列的 `Complete` 处理，才可以得到最终的完整选项配置：

```
// 设置默认参数选项
completedOptions, err := Complete(s)
if err != nil {
 return err
}
```

设置完后，还需要对参数进行合法性校验：

```
// 参数选项的校验
if errs := completedOptions.Validate(); len(errs) != 0 {
 return utilerrors.NewAggregate(errs)
}
```

在这里，如果参数不合法或参数配置失败就会校验失败，`errs` 会返回具体的错误（比如我这里是没有创建自签名证书文件的权限）：

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

为了让 apiserver 正常启动，必须为其配置参数，其中最基本的 etcd 的参数配置（这里为了专注于 k8s 自身组件，etcd 就搭建了一个本地单机版）：

```
--etcd-servers=http://127.0.0.1:2379
```

而 apiserver 的 tls 证书配置，可以用 **easyrsa**[5] 来生成：

```
curl -LO https://dl.k8s.io/easy-rsa/easy-rsa.tar.gz
tar xzf easy-rsa.tar.gz
cd easy-rsa-master/easyrsa3
./easyrsa init-pki
export MASTER_IP=127.0.0.1
export MASTER_CLUSTER_IP=127.0.0.1
./easyrsa --batch "--req-cn=${MASTER_IP}@`date +%s`" build-ca nopass
./easyrsa --subject-alt-name="IP:${MASTER_IP},"\
"IP:${MASTER_CLUSTER_IP},"\
"DNS:kubernetes,"\
"DNS:kubernetes.default,"\
"DNS:kubernetes.default.svc,"\
"DNS:kubernetes.default.svc.cluster,"\
"DNS:kubernetes.default.svc.cluster.local" \
--days=10000 \
build-server-full server nopass
```

将生成的证书（ `pki/ca.crt`、`pki/issued/server.crt` 和  `pki/private/server.key`）拷贝到自定义目录（例如：`cert`），并添加到 apiserver 的启动参数中：

```
--client-ca-file=cert/ca.crt
--tls-cert-file=cert/server.crt
--tls-private-key-file=cert/server.key
```

另外 pod 访问 kube-apiserver 还需要用到 service account 证书（`pki/private/ca.key`），需要指定服务帐号令牌颁发者 ：

```
--service-account-issuer=api
--service-account-key-file=cert/ca.crt
--service-account-signing-key-file=cert/ca.key
```

总结，启动所需要的参数：

```
--etcd-servers=http://127.0.0.1:2379 --client-ca-file=cert/ca.crt --tls-cert-file=cert/server.crt --tls-private-key-file=cert/server.key --service-account-issuer=api --service-account-key-file=cert/ca.crt --service-account-signing-key-file=cert/ca.key
```

此时，`cert` 目录有以下证书文件（包括 ca 证书、服务端证书）：

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

apiserver 启动成功后，可以直接使用 curl 充当一个客户端调用 apiserver 的接口，但是得先为客户端创建客户端证书：

```
# 生成客户端私钥
openssl genrsa -out cert/client.key 2048
# 生成证书签名请求（CSR）
openssl req -new -key cert/client.key -out cert/client.csr -subj "/CN=<client-name>"
# 使用 kube-apiserver 的 CA 签署 CSR，生成客户端证书
openssl x509 -req -in cert/client.csr -CA cert/ca.crt -CAkey cert/ca.key -CAcreateserial -out cert/client.crt -days 365
```

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

有了客户端证书，就可以调用 apiserver 接口了：

```
$ curl --cacert cert/ca.crt --cert cert/client.crt --key cert/client.key https://127.0.0.1:6443/version
{
  "major": "",
  "minor": "",
  "gitVersion": "v0.0.0-master+$Format:%H$",
  "gitCommit": "$Format:%H$",
  "gitTreeState": "",
  "buildDate": "1970-01-01T00:00:00Z",
  "goVersion": "go1.20.3",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

对于 kubectl 的使用，是同样的道理：

```
# 添加新集群（apiserver地址和ca证书）
kubectl config set-cluster devk8s --server=https://127.0.0.1:6443 --certificate-authority=cert/ca.crt
# 添加用户（客户端证书）
kubectl config set-credentials devk8s --client-certificate=cert/client.crt --client-key=cert/client.key
# 添加上下文（绑定集群和用户）
kubectl config set-context devk8s --user=devk8s --cluster=devk8s
# 切换当前上下文
kubectl config use-context devk8s
```

验证 kubectl ：

```
$ kubectl get ns
NAME              STATUS   AGE
default           Active   2d6h
kube-node-lease   Active   2d6h
kube-public       Active   2d6h
kube-system       Active   2d6h
```

创建一个 pod ：

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

如果遇到 `pods "nginx" is forbidden: error looking up service account defau lt/default: serviceaccount "default" not found` 错误，意味着还没有默认的 `serviceaccount` ，则先执行：

```
kubectl create sa default
```

创建 pod 后观察（ k8s 的组件未全部启动，处于 Pending 状态为正常现象）：

```
$ kubectl get pod
NAME    READY   STATUS    RESTARTS   AGE
nginx   0/1     Pending   0          13s
```

当然，有了 kubectl ，就可以直接使用 proxy 模式转发 apiserver 的接口了：

```
kubectl proxy --port=8080
```

apiserver 的 API 使用了 openapi 规范，暴露出了 `/openapi/v2` 接口，所以如果有需要，可以自己搭建一个 swagger-ui 服务导入该接口地址进行调试。

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

此时，etcd 存储的数据如下（使用 ETCD Keeper 可视化展示）：

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmeSOzNRNEjazLwjC02QfJPDfgYhv8ao5icLZhm0tLaRpsoPbLcaJyk31DsorAfxGJkWwOicb4pjZOQA/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

至此，先告一段落，现在 apiserver 已经成功启动，为了后续方便调试 apiserver ，`kubectl` 工具，`etcd` 数据展示也已准备好。

### 参考资料

[1]1.27.2: https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.27.md#v1272

[2]k8s.io: http://k8s.io/

[3]github.com/spf13/cobra: http://github.com/spf13/cobra

[4]https://kubernetes.io/zh-cn/docs/reference/command-line-tools-reference/kube-apiserver/: https://kubernetes.io/zh-cn/docs/reference/command-line-tools-reference/kube-apiserver/

[5]easyrsa: https://kubernetes.io/zh-cn/docs/tasks/administer-cluster/certificates/#distributing-self-signed-ca-certificate