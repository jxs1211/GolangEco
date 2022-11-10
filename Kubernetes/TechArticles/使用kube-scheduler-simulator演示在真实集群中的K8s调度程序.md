# 用kube-scheduler-simulator编写自己的调度程序

CNCF [CNCF](javascript:void(0);) *2022-08-26 10:45* *Posted on 香港*

*客座文章最初在**Miraxia 博客**[1]上发表*

由于默认的 Kubernetes 调度程序是高度可配置的，在许多情况下，我们不必编写任何代码来定制调度行为。然而，想要了解调度程序如何工作，以及如何与其他组件交互的人，可以尝试开发自己的调度程序。

在本文中，我将描述如何借助**kube-scheduler-simulator**[2]构建一个调度程序开发环境。

## 思路

1. 使用 kube-scheduler-simulator，它提供了一种简单的方法来开发调度程序，而无需准备真正的集群
2. 给 kube-scheduler-simulator 添加一个最小的调度器实现，因为默认的实现太灵活了，对初学者来说太复杂了
3. 修改和评估调度算法

## 设置

首先，让我们设置并尝试 kube-scheduler-simulator。这个过程很简单。

执行以下命令：

```sh
$ git clone https://github.com/kubernetes-sigs/kube-scheduler-simulator.git
$ cd kube-scheduler-simulator
$ git checkout 9de8c472f348b31437cce5ca2a34506f874cdddb
$ make docker_build_and_up
```

仅供参考，

- 我用 commit 9de8c472f348b31437cce5ca2a34506f874cdddb 进行了测试
- 如果你的电脑在代理后面，你可能需要取消设置 http_proxy 或者类似的东西
- 如果你在 ssh 上工作，请为 tcp 1212、3000 和 3131 创建隧道

然后，打开http://localhost:3000，尝试添加一些节点和pod。默认行为非常直观，pod分布在节点上，从而使每个节点的负载相等。
![image.png](https://www.miraxia.com/wpimages/uploads/2022/06/image-3.png)

## 向 kube-scheduler-simulator 添加一个最小调度程序

我使用了由 Kensei Nakada-san 开发的“minisched”，他是 kube-scheduler-simulator 的作者，作为我们将要开发的新调度程序的基础实现。

minisched 是**mini-kube-scheduler**[3]的一部分，mini-kube-scheduler 是一个演示系统，为教育目的而设计。虽然 mini-kube-scheduler 是基于 kube-scheduler-simulator 代码的，但是你只能使用 mini-kube-scheduler 来开发你的调度程序。而由于 mini-kube-scheduler 似乎几个月没有更新了，我决定把这两个结合起来。

为了在 kube-scheduler-simulator 中使用 minisched，需要以下步骤。

1. 将**mini-kube-scheduler/minisched**[4]（从 initial-random-scheduler 分支）复制到 kube-scheduler-simulator
2. 修改 kube-scheduler-simulator/scheduler/scheduler 以使用 minisched(见下面附上的补丁)
3. 检查行为变化（minisched 随机绑定 pod 和节点）

```go
Patch license: Apache-2.0 (same as kube-scheduler-simulator)

diff --git a/scheduler/scheduler.go b/scheduler/scheduler.go
index a5d5ca2..8eb931d 100644
--- a/scheduler/scheduler.go
+++ b/scheduler/scheduler.go
@@ -3,6 +3,8 @@ package scheduler
 import (
     "context"

+    "github.com/kubernetes-sigs/kube-scheduler-simulator/minisched"
+
     "golang.org/x/xerrors"
     v1 "k8s.io/api/core/v1"
     clientset "k8s.io/client-go/kubernetes"
@@ -14,7 +16,6 @@ import (
     "k8s.io/kubernetes/pkg/scheduler/apis/config"
     "k8s.io/kubernetes/pkg/scheduler/apis/config/scheme"
     "k8s.io/kubernetes/pkg/scheduler/apis/config/v1beta2"
-    "k8s.io/kubernetes/pkg/scheduler/profile"

     simulatorschedconfig "github.com/kubernetes-sigs/kube-scheduler-simulator/scheduler/config"
     "github.com/kubernetes-sigs/kube-scheduler-simulator/scheduler/plugin"
@@ -59,7 +60,6 @@ func (s *Service) ResetScheduler() error {
 // StartScheduler starts scheduler.
 func (s *Service) StartScheduler(versionedcfg *v1beta2config.KubeSchedulerConfiguration) error {
     clientSet := s.clientset
-    restConfig := s.restclientCfg
     ctx, cancel := context.WithCancel(context.Background())

     informerFactory := scheduler.NewInformerFactory(clientSet, 0)
@@ -71,36 +71,10 @@ func (s *Service) StartScheduler(versionedcfg *v1beta2config.KubeSchedulerConfig

     s.currentSchedulerCfg = versionedcfg.DeepCopy()

-    cfg, err := convertConfigurationForSimulator(versionedcfg)
-    if err != nil {
-        cancel()
-        return xerrors.Errorf("convert scheduler config to apply: %w", err)
-    }
-
-    registry, err := plugin.NewRegistry(informerFactory, clientSet)
-    if err != nil {
-        cancel()
-        return xerrors.Errorf("plugin registry: %w", err)
-    }
-
-    sched, err := scheduler.New(
+    sched := minisched.New(
         clientSet,
         informerFactory,
-        profile.NewRecorderFactory(evtBroadcaster),
-        ctx.Done(),
-        scheduler.WithKubeConfig(restConfig),
-        scheduler.WithProfiles(cfg.Profiles...),
-        scheduler.WithPercentageOfNodesToScore(cfg.PercentageOfNodesToScore),
-        scheduler.WithPodMaxBackoffSeconds(cfg.PodMaxBackoffSeconds),
-        scheduler.WithPodInitialBackoffSeconds(cfg.PodInitialBackoffSeconds),
-        scheduler.WithExtenders(cfg.Extenders...),
-        scheduler.WithParallelism(cfg.Parallelism),
-        scheduler.WithFrameworkOutOfTreeRegistry(registry),
     )
-    if err != nil {
-        cancel()
-        return xerrors.Errorf("create scheduler: %w", err)
-    }

     informerFactory.Start(ctx.Done())
     informerFactory.WaitForCacheSync(ctx.Done())
```

minisched 随机绑定 pod 和节点

![image-1.png](https://www.miraxia.com/wpimages/uploads/2022/06/image-1-1.png)

## 修改算法

你可以通过编辑**minisched/minisched.go#L37**[5]，像这样轻松地修改 minisched：

```go
Patch license: MIT (same as mini-kube-scheduler)

diff --git a/minisched/minisched.go b/minisched/minisched.go
index 82c0043..ae02597 100644
--- a/minisched/minisched.go
+++ b/minisched/minisched.go
@@ -34,7 +34,7 @@ func (sched *Scheduler) scheduleOne(ctx context.Context) {
        klog.Info("minischeduler: Got Nodes successfully")

        // select node randomly
-       selectedNode := nodes.Items[rand.Intn(len(nodes.Items))]
+       selectedNode := nodes.Items[0]

        if err := sched.Bind(ctx, pod, selectedNode.Name); err != nil {
                klog.Error(err)
```

我们修改的调度程序将所有的 pod 绑定到同一个节点。

![image-2.png](https://www.miraxia.com/wpimages/uploads/2022/06/image-2-1.png)

## 总结

- 你可以使用 kube-scheduler-simulator 开发自己的调度程序，它从不需要真正的集群
- mini-kube-scheduler/minisched 实现帮助你从最少的代码开始

### 参考资料

[1]Miraxia 博客: *https://www.miraxia.com/engineers-blog/writing-your-own-scheduler-with-kube-scheduler-simulator/*

[2]kube-scheduler-simulator: *https://github.com/kubernetes-sigs/kube-scheduler-simulator*

[3]mini-kube-scheduler: *https://github.com/sanposhiho/mini-kube-scheduler*

[4]mini-kube-scheduler/minisched: *https://github.com/sanposhiho/mini-kube-scheduler/tree/initial-random-scheduler/minisched*

[5]minisched/minisched.go#L37: *https://github.com/sanposhiho/mini-kube-scheduler/blob/initial-random-scheduler/minisched/minisched.go#L37*