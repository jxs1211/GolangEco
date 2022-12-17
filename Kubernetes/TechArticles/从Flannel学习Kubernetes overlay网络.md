# [从 Flannel 学习 Kubernetes overlay 网络](https://mp.weixin.qq.com/s?__biz=MjM5OTg2MTM0MQ==&mid=2247486088&idx=1&sn=ee4197ebe2061e3ed15a66e1bd81d79c&scene=21&key=bceb956dadf05f825ed687d314a926e559806cf68c2e5faad3aa8537df0656fb85f451a4324d2c051ce6ce73e7e63d4673917e940b5633a90ce5d0462c48dd0b358b629437cf3ca0303321570d8f2ece0eaf0aec5be27650c7e7b6fbb09478ebe8e5e66de390e3be0a5248bfaa4a5e874525a22f4fcf91d9d9de2db2a4e31d80&ascene=7&uin=MjMxOTI2NTEwMA%3D%3D&devicetype=Windows+10+x64&version=63080029&lang=en&exportkey=n_ChQIAhIQPYirhmh9v9b9YoqGdDnZGRLYAQIE97dBBAEAAAAAAIlONmA9zocAAAAOpnltbLcz9gKNyK89dVj0qod1FcLyczoTbOBLogSlNlM3h3bmQnwPOk00cTrWARR2RAdsejOsgj2ZWAI8c%2BAiHpZGcM75hj%2F80HGjNIXf2stPLw1adqeN1UwMmMszuOc%2FjxegLCeBz5wtKYHlWzR9RPXPrqkMhIU7jnO1%2F7IbJppUDvD8nr9kyrma2QmoQcUDfLiZxtq6dgGArtbbFbucIb6%2F4uRqCsU7XKdBl12a5F66brxGxIxlpfruyhaVV%2F6TDA%3D%3D&acctmode=0&pass_ticket=t4pRIqcemiG7QLsUo2KEZYbj72f4dTXl2isssNHsBoCSkXwRh1cTWXYFWGOOfm8KxpjqnTmecjsy%2F6yRUhQ%2FOQ%3D%3D&wx_header=1&fontgear=2)

Original 张晓辉 [云原生指北](javascript:void(0);) *2022-12-15 08:19* *Posted on 广东*

收录于合集

\#kubernetes34个

\#网络9个

\#容器10个

\#CNI4个

这是 Kubernetes 网络学习的第四篇笔记。

- [深入探索 Kubernetes 网络模型和网络通信](https://mp.weixin.qq.com/s?__biz=MjM5OTg2MTM0MQ==&mid=2247485925&idx=1&sn=d30ac0a1fdd39db3b8edffba78d20cae&scene=21#wechat_redirect)
- [认识一下容器网络接口 CNI](https://mp.weixin.qq.com/s?__biz=MjM5OTg2MTM0MQ==&mid=2247485946&idx=1&sn=9757648e456be052150d11dff0d3bc42&scene=21#wechat_redirect)
- [源码分析：从 kubelet、容器运行时看 CNI 的使用](https://mp.weixin.qq.com/s?__biz=MjM5OTg2MTM0MQ==&mid=2247485960&idx=1&sn=b2f881647680d4ef568dfd6173386021&scene=21#wechat_redirect)
- 从 Flannel 学习 Kubernetes VXLAN 网络（本篇）
- Cilium CNI 与 eBPF
- ...

------

## Flannel 介绍

Flannel 是一个非常简单的 overlay 网络（VXLAN），是 Kubernetes 网络 CNI 的解决方案之一。Flannel 在每台主机上运行一个简单的轻量级 agent `flanneld` 来监听集群中节点的变更，并对地址空间进行预配置。Flannel 还会在每台主机上安装 vtep `flannel.1`（VXLAN tunnel endpoints），与其他主机通过 VXLAN 隧道相连。

flanneld 监听在 `8472` 端口，通过 UDP 与其他节点的 vtep 进行数据传输。到达 vtep 的二层包会被原封不动地通过 UDP 的方式发送到对端的 vtep，然后拆出二层包进行处理。简单说就是用四层的 UDP 传输二层的数据帧。

vxlan-tunnel

在 Kubernetes 发行版 **K3S**[1] 中将 Flannel 作为默认的 CNI 实现。K3S 集成了 flannel，在启动后 flannel 以 go routine 的方式运行。

## 环境搭建

Kubernetes 集群使用 k3s 发行版，但在安装集群的时候，禁用 k3s 集成的 flannel，使用独立安装的 flannel 进行验证。

安装 CNI 的 plugin，需要在所有的 node 节点上执行下面的命令，下载 CNI 的官方 bin。

```
sudo mkdir -p /opt/cni/bin
curl -sSL https://github.com/containernetworking/plugins/releases/download/v1.1.1/cni-plugins-linux-amd64-v1.1.1.tgz | sudo tar -zxf - -C /opt/cni/bin
```

安装 k3s 的控制平面。

```
export INSTALL_K3S_VERSION=v1.23.8+k3s2
curl -sfL https://get.k3s.io | sh -s - --disable traefik --flannel-backend=none --write-kubeconfig-mode 644 --write-kubeconfig ~/.kube/config
```

安装 Flannel。**这里注意，Flannel 默认的 Pod CIRD 是 `10.244.0.0/16`，我们将其修改为 k3s 默认的 `10.42.0.0/16`。**

```
curl -s https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml | sed 's|10.244.0.0/16|10.42.0.0/16|g' | kubectl apply -f -
```

添加另一个节点到集群。

```
export INSTALL_K3S_VERSION=v1.23.8+k3s2
export MASTER_IP=<MASTER_IP>
export NODE_TOKEN=<TOKEN>
curl -sfL https://get.k3s.io | K3S_URL=https://${MASTER_IP}:6443 K3S_TOKEN=${NODE_TOKEN} sh -
```

查看节点状态。

```
kubectl get node
NAME          STATUS   ROLES                  AGE   VERSION
ubuntu-dev3   Ready    <none>                 13m   v1.23.8+k3s2
ubuntu-dev2   Ready    control-plane,master   17m   v1.23.8+k3s2
```

运行两个 pod：`curl` 和 `httpbin`，为了探寻

```
NODE1=ubuntu-dev2
NODE2=ubuntu-dev3
kubectl apply -n default -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: curl
  name: curl
spec:
  containers:
  - image: curlimages/curl
    name: curl
    command: ["sleep", "365d"]
  nodeName: $NODE1
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: httpbin
  name: httpbin
spec:
  containers:
  - image: kennethreitz/httpbin
    name: httpbin
  nodeName: $NODE2
EOF
```

## 网络配置

接下来，一起看下 CNI 插件如何配置 pod 网络。

### 初始化

Flannel 是通过 `Daemonset` 的方式部署的，每台节点上都会运行一个 flannel 的 pod。通过挂载本地磁盘的方式，在 Pod 启动时会通过初始化容器将二进制文件和 CNI 的配置复制到本地磁盘中，分别位于 `/opt/cni/bin/flannel` 和 `/etc/cni/net.d/10-flannel.conflist`。

通过查看 **kube-flannel.yml**[2] 中的 `ConfigMap`，可以找到 CNI 配置，flannel 默认委托（见 **flannel-cni 源码 `flannel_linux.go#L78`**[3]）给 **bridge 插件**[4] 进行网络配置，网络名称为 `cbr0`；IP 地址的管理，默认委托（见 **flannel-cni 源码 `flannel_linux.go#L40`**[5]） **host-local 插件**[6] 完成。

```
#cni-conf.json 复制到 /etc/cni/net.d/10-flannel.conflist
{
  "name": "cbr0",
  "cniVersion": "0.3.1",
  "plugins": [
    {
      "type": "flannel",
      "delegate": {
        "hairpinMode": true,
        "isDefaultGateway": true
      }
    },
    {
      "type": "portmap",
      "capabilities": {
        "portMappings": true
      }
    }
  ]
}
```

还有 Flannel 的网络配置，配置中有我们设置的 Pod CIDR `10.42.0.0/16` 以及后端（backend）的类型 `vxlan`。这也是 flannel 默认的类型，此外还有 **多种后端类型**[7] 可选，如 `host-gw`、`wireguard`、`udp`、`Alloc`、`IPIP`、`IPSec`。

```
#net-conf.json 挂载到 pod 的 /etc/kube-flannel/net-conf.json
{
  "Network": "10.42.0.0/16",
  "Backend": {
    "Type": "vxlan"
  }
}
```

Flannel Pod 运行启动 `flanneld` 进程，指定了参数 `--ip-masq` 和 `--kube-subnet-mgr`，后者开启了 `kube subnet manager` 模式。

### 运行

![Image](https://mmbiz.qpic.cn/mmbiz_png/tMghG0NOfxeZK0TOZcYFcVsgpxgpS9mBgVwmiaFeiaCHEV9uS1gTbBibrKNuz2MhS4R2yNCqpoOeQBrW0VpLLA71w/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

集群初始化时使用了默认的 Pod CIDR `10.42.0.0/16`，当有节点加入集群，集群会从该网段上为节点分配 **属于节点的 Pod CIDR** `10.42.X.1/24`。

flannel 在 `kube subnet manager` 模式下，连接到 apiserver 监听节点更新的事件，从节点信息中获取节点的 Pod CIDR。

```
kubectl get no ubuntu-dev2 -o jsonpath={.spec} | jq
{
  "podCIDR": "10.42.0.0/24",
  "podCIDRs": [
    "10.42.0.0/24"
  ],
  "providerID": "k3s://ubuntu-dev2"
}
```

然后在主机上写子网配置文件，下面展示的是其中一个节点的子网配置文件的内容。另一个节点的内容差异在 `FLANNEL_SUBNET=10.42.1.1/24`，使用的是对应节点的 Pod CIDR。

```
#node 192.168.1.12
cat /run/flannel/subnet.env
FLANNEL_NETWORK=10.42.0.0/16
FLANNEL_SUBNET=10.42.0.1/24
FLANNEL_MTU=1450
FLANNEL_IPMASQ=true
```

### CNI 插件执行

CNI 插件的执行是由容器运行时触发的，具体细节可以看上一篇 [《源码解析：从 kubelet、容器运行时看 CNI 的使用》](https://mp.weixin.qq.com/s?__biz=MjM5OTg2MTM0MQ==&mid=2247485960&idx=1&sn=b2f881647680d4ef568dfd6173386021&scene=21#wechat_redirect)。

![Image](https://mmbiz.qpic.cn/mmbiz_png/tMghG0NOfxeZK0TOZcYFcVsgpxgpS9mBM5JbDSAhpQBW6r9JxImfNxRoW6iaGIHnW6TOtNqSeguPLTjfzQEEWBA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)Flannel Plugin Flow

#### flannel 插件

`flannel` CNI 插件（`/opt/cni/bin/flannel`）执行的时候，接收传入的 `cni-conf.json`，读取上面初始化好的 `subnet.env` 的配置，输出结果，委托给 `bridge` 进行下一步。

```
cat /var/lib/cni/flannel/e4239ab2706ed9191543a5c7f1ef06fc1f0a56346b0c3f2c742d52607ea271f0 | jq
{
  "cniVersion": "0.3.1",
  "hairpinMode": true,
  "ipMasq": false,
  "ipam": {
    "ranges": [
      [
        {
          "subnet": "10.42.0.0/24"
        }
      ]
    ],
    "routes": [
      {
        "dst": "10.42.0.0/16"
      }
    ],
    "type": "host-local"
  },
  "isDefaultGateway": true,
  "isGateway": true,
  "mtu": 1450,
  "name": "cbr0",
  "type": "bridge"
}
```

#### bridge 插件

`bridge` 使用上面的输出连同参数一起作为输入，根据配置完成如下操作：

1. 创建网桥 `cni0`（节点的根网络命名空间）
2. 创建容器网络接口 `eth0`（ pod 网络命名空间）
3. 创建主机上的虚拟网络接口 `vethX`（节点的根网络命名空间）
4. 将 `vethX` 连接到网桥 `cni0`
5. 委托 ipam 插件分配 IP 地址、DNS、路由
6. 将 IP 地址绑定到 pod 网络命名空间的接口 `eth0` 上
7. 检查网桥状态
8. 设置路由
9. 设置 DNS

最后输出如下的结果：

```
cat /var/li/cni/results/cbr0-a34bb3dc268e99e6e1ef83c732f5619ca89924b646766d1ef352de90dbd1c750-eth0 | jq .result
{
  "cniVersion": "0.3.1",
  "dns": {},
  "interfaces": [
    {
      "mac": "6a:0f:94:28:9b:e7",
      "name": "cni0"
    },
    {
      "mac": "ca:b4:a9:83:0f:d4",
      "name": "veth38b50fb4"
    },
    {
      "mac": "0a:01:c5:6f:57:67",
      "name": "eth0",
      "sandbox": "/var/run/netns/cni-44bb41bd-7c41-4860-3c55-4323bc279628"
    }
  ],
  "ips": [
    {
      "address": "10.42.0.5/24",
      "gateway": "10.42.0.1",
      "interface": 2,
      "version": "4"
    }
  ],
  "routes": [
    {
      "dst": "10.42.0.0/16"
    },
    {
      "dst": "0.0.0.0/0",
      "gw": "10.42.0.1"
    }
  ]
}
```

#### port-mapping 插件

该插件会将来自主机上一个或多个端口的流量转发到容器。

## Debug

让我们在第一个节点上，使用 `tcpdump` 对接口 `cni0` 进行抓包。

```
tcpdump -i cni0 port 80 -vvv
```

从 pod `curl` 中使用 pod `httpbin` 的 IP 地址 `10.42.1.2` 发送请求：

```
kubectl exec curl -n default -- curl -s 10.42.1.2/get
```

### cni0

从在 `cni0` 上的抓包结果来看，第三层的 IP 地址均为 Pod 的 IP 地址，看起来就像是两个 pod 都在同一个网段。

![Image](https://mmbiz.qpic.cn/mmbiz_png/tMghG0NOfxeZK0TOZcYFcVsgpxgpS9mBJjtTErR0ziaRMw5dAtD8zp9zUKejsflU56nmHCQFRbRNrDGqEQUuk0g/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)tcpdump-on-cni0

### host eth0

文章开头提到 flanneld 监听 udp 8472 端口。

```
netstat -tupln | grep 8472
udp        0      0 0.0.0.0:8472            0.0.0.0:*                           -
```

我们直接在以太网接口上抓取 UDP 的包：

```
tcpdump -i eth0 port 8472 -vvv
```

再次发送请求，可以看到抓取到 UDP 数据包，传输的负载是二层的封包。

![Image](https://mmbiz.qpic.cn/mmbiz_jpg/tMghG0NOfxeZK0TOZcYFcVsgpxgpS9mBibULibvHneLfda9ia8BP25UCiafuciaRibbc7GFdGFwKWlyLyH3kwAIBRtQw/640?wx_fmt=jpeg&wxfrom=5&wx_lazy=1&wx_co=1)tcpdump-on-host-eth0

## Overlay 网络下的跨节点通信

在系列的第一篇中，我们研究 pod 间的通信时提到不同 CNI 插件的处理方式不同，这次我们探索了 flannel 插件的工作原理。希望通过下面的图可以对 overlay 网络处理跨节点的网络通信有个比较直观的认识。

![Image](https://mmbiz.qpic.cn/mmbiz_gif/tMghG0NOfxeZK0TOZcYFcVsgpxgpS9mB96wkia7Oib6naKdP2wyW4gClibKWzdZcLDo96vz8TWXs6xlo9fibrKWJTQ/640?wx_fmt=gif&wxfrom=5&wx_lazy=1)

当发送到 `10.42.1.2` 流量到达节点 A 的网桥 `cni0`，由于目标 IP 并不属于当前阶段的网段。根据系统的路由规则，进入到接口 `flannel.1`，也就是 VXLAN 的 vtep。这里的路由规则也由 `flanneld` 来维护，当节点上线或者下线时，都会更新路由规则。

```
#192.168.1.12
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
default         _gateway        0.0.0.0         UG    0      0        0 eth0
10.42.0.0       0.0.0.0         255.255.255.0   U     0      0        0 cni0
10.42.1.0       10.42.1.0       255.255.255.0   UG    0      0        0 flannel.1
192.168.1.0     0.0.0.0         255.255.255.0   U     0      0        0 eth0
#192.168.1.13
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
default         _gateway        0.0.0.0         UG    0      0        0 eth0
10.42.0.0       10.42.0.0       255.255.255.0   UG    0      0        0 flannel.1
10.42.1.0       0.0.0.0         255.255.255.0   U     0      0        0 cni0
192.168.1.0     0.0.0.0         255.255.255.0   U     0      0        0 eth0
```

`flannel.1` 将原始的以太封包使用 UDP 协议重新封装，将其发送到目标地址 `10.42.1.0` （目标的 MAC 地址通过 ARP 获取）。对端的 vtep 也就是 `flannel.1` 的 UDP 端口 8472 收到消息，解帧出以太封包，然后对以太封包进行路由处理，发送到接口 `cni0`，最终到达目标 pod 中。

响应的数据传输与请求的处理也是类似，只是源地址和目的地址调换。

### 参考资料

[1] K3S: *https://k3s.io/*

[2] kube-flannel.yml: *https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml*

[3] flannel-cni 源码 `flannel_linux.go#L78`: *https://github.com/flannel-io/cni-plugin/blob/v1.1.0/flannel_linux.go#L78*

[4] bridge 插件: *https://www.cni.dev/plugins/current/main/bridge/*

[5] flannel-cni 源码 `flannel_linux.go#L40`: *https://github.com/flannel-io/cni-plugin/blob/v1.1.0/flannel_linux.go#L40*

[6] host-local 插件: *https://www.cni.dev/plugins/current/ipam/host-local/*

[7] 多种后端类型: *https://github.com/flannel-io/flannel/blob/master/Documentation/backends.md*