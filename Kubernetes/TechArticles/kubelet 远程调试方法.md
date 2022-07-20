# kubelet 远程调试方法

https://mp.weixin.qq.com/s/j5qt-1X_kCtZZbh17yEu0g

[Go编程时光](javascript:void(0);) *2022-07-20 09:06* *Posted on 福建*



## ![Image](https://mmbiz.qpic.cn/mmbiz_png/Bn6fDVNBvtQUwlEq1nPemS6Rk1iaERcXqL1w0mzm3ACII51wUKkSiadsSnLibvGdzGpqLiaicd7tkZNup81EnOVQvng/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## 1. kubelet启动命令分析 

kubelet是一个systemd服务，以使用Kubeadm工具安装的v1.23.4 k8s集群为例，该服务的配置文件路径为`/etc/systemd/system/kubelet.service.d/10-kubeadm.conf`, 内容如下:

```
# Note: This dropin only works with kubeadm and kubelet v1.11+
[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
# This is a file that "kubeadm init" and "kubeadm join" generates at runtime, populating the KUBELET_KUBEADM_ARGS variable dynamically
EnvironmentFile=-/var/lib/kubelet/kubeadm-flags.env
# This is a file that the user can use for overrides of the kubelet args as a last resort. Preferably, the user should use
# the .NodeRegistration.KubeletExtraArgs object in the configuration files instead. KUBELET_EXTRA_ARGS should be sourced from this file.
EnvironmentFile=-/etc/default/kubelet
ExecStart=
ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS
```

以我的测试环境为例，执行`ps -ef |grep /usr/bin/kubelet`, 可见kubelet启动的完整命令如下：

```
/usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --container-runtime=remote --container-runtime-endpoint=/run/containerd/containerd.sock --pod-infra-container-image=k8s.gcr.io/pause:3.6
```

如果需要修改kubelet命令，可以关闭服务后使用相同参数启动。或修改systemd配置文件后重启kubelet服务。

## 2. 编译kubelet

根据**k8s makefile 源码分析**[1]，kubelet编译命令如下：

https://github.com/kubernetes/kubernetes/blob/v1.22.4/hack/lib/golang.sh#L679

```
kube::golang::build_some_binaries() {
    ...
    go install "${build_args[@]}" "$@"
    ...
}
```

其中GOLDFLAGS, GOGCFLAGS配置如下：

https://github.com/kubernetes/kubernetes/blob/v1.22.4/hack/lib/golang.sh#L797-L799

```
kube::golang::build_binaries() {
    ...
    goldflags="${GOLDFLAGS=-s -w} $(kube::version::ldflags)"
    goasmflags="-trimpath=${KUBE_ROOT}"
    gogcflags="${GOGCFLAGS:-} -trimpath=${KUBE_ROOT}"
    ...
}
```

为了保留尽可能多的调试信息，我们需要重新设置这两个编译参数，所以编译kubelet的命令应为

```
git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes
git checkout v1.22.4
./build/shell.sh
make generated_files
make -o generated_files kubelet KUBE_BUILD_PLATFORMS=linux/amd64 GOLDFLAGS="" GOGCFLAGS="all=-N -l"
```

编译完成后，kubelet二进制文件位于`_output/bin/kubelet`

## 3. delve介绍

**delve**[2]是一个用于Go编程语言的调试器。**尽管我们也可以使用gdb调试go语言程序**[3], 但在调试用标准工具链构建的Go程序时，Delve是GDB更好的替代品。它比GDB更能理解Go的运行时、数据结构和表达式。

可以使用如下命令安装dlv:

```
go install github.com/go-delve/delve/cmd/dlv@latest
```

使用如下命令使用dlv进行调试:

```
dlv exec ./hello -- server --config conf/config.toml
```

以kubelet为例，使用dlv命令行调试的过程如下：

```
root@st0n3-host:~# dlv exec /usr/bin/kubelet -- --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --container-runtime=remote --container-runtime-endpoint=/run/containerd/containerd.sock --pod-infra-container-image=k8s.gcr.io/pause:3.6
Type 'help' for list of commands.
(dlv) b main.main
Breakpoint 1 set at 0x502e086 for main.main() _output/dockerized./cmd/kubelet/kubelet.go:39
(dlv) c
> main.main() _output/dockerized./cmd/kubelet/kubelet.go:39 (hits goroutine(1):1 total:1) (PC: 0x502e086)
    34:  _ "k8s.io/component-base/metrics/prometheus/restclient"
    35:  _ "k8s.io/component-base/metrics/prometheus/version" // for version metric registration
    36:  "k8s.io/kubernetes/cmd/kubelet/app"
    37: )
    38: 
=>  39: func main() {
    40:  command := app.NewKubeletCommand()
    41: 
    42:  // kubelet uses a config file and does its own special
    43:  // parsing of flags and that config file. It initializes
    44:  // logging after it is done with that. Therefore it does
(dlv) 
```

## 4. GoLand远程调试kubelet

我们当然可以使用上文描述的命令行形式进行调试，但kubernetes代码量巨大，使用IDE会更方便。

点击调试按钮左侧的Edit Configurations按钮，配置dlv的地址和端口：

![Image](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtUmrtjTqlSyibpucgXRYyHEwnSQfrefMgJZoJ3vKiazrdU2HRQz0b5Kia4MWF8hRj5RpRKVI9iaakKTg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)img

使用IDE提示的命令启动kubelet，或将其配置到systemd服务中后重启服务：

```
root@st0n3-host:~# cat /etc/systemd/system/kubelet.service.d/10-kubeadm.conf 
...
ExecStart=/usr/bin/dlv --listen=:10086 --headless=true --api-version=2 --accept-multiclient exec /usr/bin/kubelet -- $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS
root@st0n3-host:~# systemctl daemon-reload
root@st0n3-host:~# systemctl restart kubelet.service
```

此时kubelet命令实际还未真正启动，在GoLand中运行刚刚添加的配置，连接上dlv后，kubelet才会运行。

下好断点，点击debug按钮，我们就可以在IDE中对kubelet进行调试了。

![Image](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtUmrtjTqlSyibpucgXRYyHEQmeaXfAuq1Na0viaZ5xoZVBoRN1r4tzeib2zQXfKRyiaiasjaZqtYwkwHw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## 5. 其他容器软件调试命令

### 5.1 runc

#### 编译

```
make shell
make EXTRA_FLAGS='-gcflags="all=-N -l"'
```

#### 调试

```
mv /usr/bin/runc /usr/bin/runc.bak
cat <<EOF > /usr/bin/runc
#!/bin/bash
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec /usr/bin/runc.debug -- $*
chmod +x /usr/bin/runc
```

### 5.2 docker-cli

#### 编译

修改`scripts/build/binary`中的编译命令如下：删除LDFLAGS，添加gcflags

```
root@st0n3:~/cli# git diff
diff --git a/scripts/build/binary b/scripts/build/binary
index e4c5e12a6b..155528e501 100755
--- a/scripts/build/binary
+++ b/scripts/build/binary
@@ -74,7 +74,7 @@ fi
 echo "Building $GO_LINKMODE $(basename "${TARGET}")"
 
 export GO111MODULE=auto
-
-go build -o "${TARGET}" -tags "${GO_BUILDTAGS}" --ldflags "${LDFLAGS}" ${GO_BUILDMODE} "${SOURCE}"
+go build -o "${TARGET}" -tags "${GO_BUILDTAGS}" -gcflags="all=-N -l" ${GO_BUILDMODE} "${SOURCE}"
 
 ln -sf "$(basename "${TARGET}")" "$(dirname "${TARGET}")/docker"
make -f docker.Makefile shell
make binary
```

#### 调试

```
cat <<EOF > docker.debug
#!/bin/bash
dlv --listen=:2344 --headless=true --api-version=2 --accept-multiclient exec ./docker-cli.debug -- $*
chmod +x docker.debug
```

### 5.3 dockerd

#### 编译

修改`hack/make/.binary`文件中的编译命令

```
root@st0n3:~/moby# git diff
diff --git a/hack/make/.binary b/hack/make/.binary
index d56e3f3126..3e23865c81 100644
--- a/hack/make/.binary
+++ b/hack/make/.binary
@@ -81,11 +81,11 @@ hash_files() {
 
        echo "Building: $DEST/$BINARY_FULLNAME"
        echo "GOOS=\"${GOOS}\" GOARCH=\"${GOARCH}\" GOARM=\"${GOARM}\""
-       go build \
+       set -x
+       go build -gcflags "all=-N -l"  \
                -o "$DEST/$BINARY_FULLNAME" \
                "${BUILDFLAGS[@]}" \
                -ldflags "
-               $LDFLAGS
                $LDFLAGS_STATIC_DOCKER
                $DOCKER_LDFLAGS
        " \
make BIND_DIR=. shell
hack/make.sh binary
```

#### 调试

```
/root/go/bin/dlv --listen=:2343 --headless=true --api-version=2 --accept-multiclient exec /usr/bin/dockerd.debug -- -D -H unix:///var/run/docker.sock --containerd=/run/containerd/containerd.sock
```

> 原文地址：https://ssst0n3.github.io/post/%E7%BD%91%E7%BB%9C%E5%AE%89%E5%85%A8/%E5%AE%89%E5%85%A8%E7%A0%94%E7%A9%B6/%E5%AE%B9%E5%99%A8%E5%AE%89%E5%85%A8/%E5%AE%B9%E5%99%A8%E9%9B%86%E7%BE%A4%E5%AE%89%E5%85%A8/k8s/%E6%BA%90%E7%A0%81%E5%AE%A1%E8%AE%A1/%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91%E5%B9%B6%E7%BC%96%E8%AF%91%E4%BB%A3%E7%A0%81/kubelet-%E8%BF%9C%E7%A8%8B%E8%B0%83%E8%AF%95.html

### 参考资料

[1]k8s makefile 源码分析: https://ssst0n3.github.io/post/网络安全/安全研究/容器安全/容器集群安全/k8s/源码审计/如何开发并编译代码/k8s-makefile-源码分析.html[2]delve: https://github.com/go-delve/delve[3]尽管我们也可以使用gdb调试go语言程序: https://go.dev/doc/gdb