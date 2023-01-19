# K8S - Creating a kube-scheduler plugin

Saying it in a few words, the K8S scheduler is responsible for assigning *Pods* to *Nodes.* Once a new pod is created it gets in the scheduling queue. The attempt to schedule a pod is split in two phases: the *Scheduling* and the *Binding cycle.*

In the *Scheduling cycle* the nodes are filtered, removing those that don’t meet the pod requirements. Next, the *feasible nodes* (the remaining ones)*,* are ranked based on a given score. Finally, the node with highest score is chosen. These steps are called *Filtering* and *Scoring [1].*

Once a node is chosen, the scheduler needs to make sure *kubelet* knows it needs to start the pod (containers) in the selected node. The step related to starting the pod into the selected node is called *Binding Cycle [2].*

The *Scheduling* and *Binding cycle* are composed by stages that are executed sequentially to calculate the pod placement. These stages are called extension points and can be used to shape the placement behavior. *Scheduling Cycles* for different pods are run sequentially, meaning that the *Scheduling Cycle* steps will be executed for one pod at a time, whereas *Binding Cycles* for different pods may be executed concurrently.

The components that implement the extension points of kubernetes scheduler are called *Plugins.* The native scheduling behavior is implemented using the *Plugin* pattern as well, in the same way that custom extensions, making the core of the kube-scheduler lightweight as the main scheduling logic is placed in the plugins.

The extension points where the plugins can be applied are shown in *Figure 1.* A plugin can implement one or more of the extension points and a detailed description of each can be found in [4] (I won’t be copying stuff here, check it there before continuing 😃).

![img](https://miro.medium.com/v2/resize:fit:700/1*EZwRBzbU4Bs3mE1tIRYPEg.png)

Figure 1

To configure the *Plugins* that should be executed in each extension point, and then change the scheduling behavior, kube-scheduler provides *Profiles* [3]*.* A scheduling *Profile* describes which plugins should be executed on each stage mentioned in [4]. It is possible to provide multiple profiles, which means that there’s no need to deploy multiple schedulers to have different scheduling behaviors [5].

# kube-scheduler

The kube-scheduler is implemented in Golang and *Plugins* are included to it in compilation time. Therefore, if you want to have your own plugin, you will need to have your own scheduler image.

A new plugin needs to be registered and get configured to the plugin API. Also, it needs to implement the extension points interfaces that are defined in the [kubernetes scheduler framework package](https://github.com/kubernetes/kubernetes/blob/ed3e0d302fb546653b78df583569b0311687a7a8/pkg/scheduler/framework/interface.go#L268). Check out how it looks:

<iframe src="https://medium.com/media/7c7a46ce88db72ba44e58b24d81b2dc9" allowfullscreen="" frameborder="0" height="369" width="680" title="Scheduler Plugin - Interfaces example" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 369px; position: absolute; left: 0px;"></iframe>

The scheduler’s code allows to add new plugins without having to fork it. For that, developers just need to write their own `main()` wrapper around the scheduler. As plugins must be compiled with the scheduler, writing a wrapper allows to re-use the scheduler’s code in a clean way [7].

To do that, the main function will import the `k8s.io/kubernetes/cmd/kube-scheduler/app` and use the `NewSchedulerCommand` to register the custom plugins, providing the respective name and the constructor function:

<iframe src="https://medium.com/media/273ed6b707fb4e7fa927e2b000505289" allowfullscreen="" frameborder="0" height="325" width="680" title="Scheduler Plugin - main.go" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 325px; position: absolute; left: 0px;"></iframe>

## Configuration

The kube-scheduler configuration is where the profiles can be configured. Each profile allows plugins to be enabled, disabled and configured according to the configurations parameters defined by the plugin. Each profile configuration is separated into two parts [9]:

1. A list of enabled plugins for each extension point and the order they should run. If one of the extension points list is omitted, the default list will be used.
2. An optional set of custom plugin arguments for each plugin. Omitting config args for a plugin is equivalent to using the default config for that plugin.

Plugins that are enabled in different extension points must be configured explicitly in each of them.

The configuration is provided through the `KubeSchedulerConfiguration` struct. To enable it, it needs to be written to a configuration file and its path provided as a command line argument to kube-scheduler. E.g.:

```
kube-scheduler --config=/etc/kubernetes/networktraffic-config.yaml
```

Below you can see an example configuration of the `NetworkTraffic` plugin. In the example, the `clientConnection.kubeconfig` points to the kubeconfig path used by the kube-scheduler, with its defined authorizations in the control plane nodes. The `profiles` section overwrites the `default-scheduler` score phase enabling the `NetworkTraffic` plugin and disabling the others defined by default. The `pluginConfig` sets the configuration of the plugin, that will be provided during its initialization [8].

<iframe src="https://medium.com/media/e0045789f9397f4527340fb7b1299da4" allowfullscreen="" frameborder="0" height="435" width="680" title="Network Traffic Plugin - Configuration" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 435px; position: absolute; left: 0px;"></iframe>

PS: If you have HA with multiple control plane nodes, the configuration needs to be applied for each of them.

# Creating a custom plugin

Now that we understand the basics of kube-scheduler, we can do what we came here for. As we’ve seen previously, adding a custom plugin requires to include our code during compilation time and we don't need to fork the scheduler code for that.

To proceed, we could create an empty repository and wrap the scheduler as described before, however, the project [scheduler-plugins](https://github.com/kubernetes-sigs/scheduler-plugins) already does that and provides some custom plugins that are good examples to follow. So, we will just start from there.

Fork the [scheduler-plugins](https://github.com/kubernetes-sigs/scheduler-plugins) repository and pull it into `$GOPATH/src/sigs.k8s.io`. With that done, we can start :)

To keep following the next steps, you need to:

1. Have a K8S cluster (I am using a cluster created with [kubespray](https://github.com/kubernetes-sigs/kubespray)).
2. Have prometheus configured with node-exporter. Check [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack).

## NetworkTraffic Plugin

For this example, we are going to build a Score Plugin named "NetworkTraffic" that favors nodes with lower network traffic. To gather that information we will query prometheus.

To start, create the folder `pkg/networktraffic` and the files `networktraffic.go` and `prometheus.go` inside your fork of scheduler-plugins. The structure should look like this:

```
|- pkg
|-- networktraffic
|--- networktraffic.go
|--- prometheus.go
```

In the `networktraffic.go` we are going to have the implementation of the ScorePlugin interface and in the `prometheus.go` we will keep the logic to interact with prometheus.

## Prometheus communication

In the `prometheus.go` we will start by declaring the struct used to interact with Prometheus. It will have the fields `networkInterface` and `timeRange`, which can be used to configure the query we will be executing. The field `address` points to the prometheus service on K8S and can also be configured. The field `api` will be used to store the prometheus client, which is created based on the `address` provided.

```
type PrometheusHandle struct {
  networkInterface string
  timeRange        time.Duration
  address          string
  api              v1.API
}
```

Now that we have the basic structure we can also implement the querying. We will be using the sum of the received bytes in a time range per node in a specific network interface. The `kubernetes_node` filter will query the metrics for the node provided, as described by the query below. The `device` filter will query the metrics on the provided network interface, and the last value between `[%s]` defines the time range taken into account. `sum_over_time` will sum all the values in the provided time range.

```
sum_over_time(node_network_receive_bytes_total{kubernetes_node=\"%s\",device=\"%s\"}[%s])
```

At the end, the `prometheus.go` file will look like this:

<iframe src="https://medium.com/media/5bed130b841a179d2fb0386b1c6f26b1" allowfullscreen="" frameborder="0" height="1596" width="680" title="networktraffic - prometheus.go" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 1596px; position: absolute; left: 0px;"></iframe>

## ScorePlugin interface

Having the interaction with Prometheus done, we can move to the implementation of the Score Plugin. As mentioned, we will need to implement the Score Plugin Interface from the scheduler framework:

<iframe src="https://medium.com/media/19c333384f20627841379639e0d929e2" allowfullscreen="" frameborder="0" height="562" width="680" title="Scheduler Plugin - Score Plugin Interface" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 562px; position: absolute; left: 0px;"></iframe>

The `Score` function is called for each node and returns whether it was successful and an integer indicating the rank of the node. At the end of the Score plugin execution, we should have a Score value in the range from 0 to 100. In some cases it could be difficult to have a value within that range without knowing the score of other nodes, for example. For those scenarios, we can use the `NormalizeScore` function implemented in the `ScoreExtensions` interface. The `NormalizeScore` function receives the result of all nodes and allows them to be changed.

Moreover, the ScorePlugin interface also have the `Plugin` interface as an [embedded](https://travix.io/type-embedding-in-go-ba40dd4264df) field. So, we must implement its `Name() string` function.

Now that we understand the ScorePlugin interface, let's go to the `networktraffic.go` file. We will start by defining the `NetworkTraffic` struct:

```
// NetworkTraffic is a score plugin that favors nodes based on their
// network traffic amount. Nodes with less traffic are favored.
// Implements framework.ScorePlugin
type NetworkTraffic struct {
  handle     framework.FrameworkHandle
  prometheus *PrometheusHandle
}
```

With the structure defined, we can proceed with the `Score` function implementation. It will be straightforward. We will only call the `GetNodeBandwidthMeasure` function from our Prometheus structure providing the node name. The call will return a `Sample` which holds the value in the `Value` field. We will basically return it for each node.

<iframe src="https://medium.com/media/e9a853d3dc912aa09873dccb7274826a" allowfullscreen="" frameborder="0" height="254" width="680" title="Network Traffic Plugin - Score function" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 253.984px; position: absolute; left: 0px;"></iframe>

Network Traffic plugin Score function

Next, we will have returned the total bytes received by each node in a determined period of time. However, the scheduler framework expects a value from 0 to 100, thus, we still need to normalize the values to fulfill this requirement.

To do the normalization, we will implement the `ScoreExtensions` interface mentioned before. We will implement the interface embedded in the `NetworkTraffic` struct. In the `ScoreExtensions` function we will simply return the struct which implements the interface. The logic is placed under the `NormalizeScore` function.

<iframe src="https://medium.com/media/8adb746ecae00f504acda9fd5dfcb69b" allowfullscreen="" frameborder="0" height="474" width="680" title="Network Traffic Plugin - Score Extensions" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 473.984px; position: absolute; left: 0px;"></iframe>

The `NormalizeScore` basically will take the highest value returned by prometheus and use it as the highest possible value, corresponding to the `framework.MaxNodeScore` (100). The other values will be calculated relatively to the highest score using the rule of three.

Finally, we will have a list where the nodes with more network traffic have a greater score in the range of [0,100]. If we use it like it is, we would favor nodes that have higher traffic, so, we need to reverse the values. For that, we will simply replace the node score with the result of the rule of three, subtracted by the max score.

An example of the calculation which take as an example three nodes (*a, b* and *c*), the values are in bytes, is given below:

```
a => 1000000   # 1MB
b => 1200000   # 1,2MB
c => 1400000   # 1,4MBhigherScore = 1400000Y = (node.Score * framework.MaxNodeScore) / higherScoreYa = 1000000 * 100 / 1400000
Yb = 1200000 * 100 / 1400000
Yc = 1400000 * 100 / 1400000Ya = 71,42
Yb = 85,71
Yc = 100Xa = 100 - Ya
Xb = 100 - Yb
Xc = 100 - YcXa = 28,58
Xb = 14,29
Xc = 0
```

With that explained, we have the main pieces of our plugin. However, that's not all. As mentioned before, the scheduler plugins can be configured, and there are three configurations we will allow in our Network Traffic plugin, which were already mentioned:

- Prometheus address
- Prometheus query time range
- Prometheus query node network interface

Those values will be provided during the instantiation of the `NetworkTraffic` plugin by the scheduler framework, and we will need to declare a new struct called `NetworkTrafficArgs` that will be used to parse the configuration provided in the `KubeSchedulerConfiguration`. For that, we need to add a new function with the logic to create an instance of the `NetworkTraffic` plugin, described below:

<iframe src="https://medium.com/media/cfbd6d4756def29acc02132f63c879cf" allowfullscreen="" frameborder="0" height="320" width="680" title="Network Traffic Plugin - Score constructor" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 320px; position: absolute; left: 0px;"></iframe>

The `New` function follows the scheduler framework`PluginFactory` interface.

We still haven't declared the `NetworkTrafficArgs`struct, and that will come next. However, we have (almost) all we need for `networktraffic.go`:

<iframe src="https://medium.com/media/52e4eb8ea2e3b995cbfb5e75415f1176" allowfullscreen="" frameborder="0" height="1728" width="680" title="Network Traffic Plugin - networktraffic.go" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 1728px; position: absolute; left: 0px;"></iframe>

## Configuration

The scheduler-plugins project holds the configurations under `pkg/apis` folder. So, we will have ours plugin config there as well.

We will add the configuration in two places: `pkg/apis/config/types.go` and `pkg/apis/config/v1beta1/types.go`. The `config/types.go` holds the struct we will use in the `New` function, while the `v1beta1/types.go` holds the struct used to parse the information from the `KubeSchedulerConfiguration`.

Also, the config struct must follow the name pattern `<Plugin Name>Args`, otherwise, it won't be properly decoded and you will face issues.

<iframe src="https://medium.com/media/98b31ac6cd96a38309c4001ed030ef1e" allowfullscreen="" frameborder="0" height="325" width="680" title="Network Traffic Plugin - config/types.go" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 325px; position: absolute; left: 0px;"></iframe>

config/types.go

<iframe src="https://medium.com/media/59c93ebebfaaa8d03d1d1f3c82e2efed" allowfullscreen="" frameborder="0" height="347" width="680" title="Network Traffic Plugin - config/v1beta1/types.go" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 347px; position: absolute; left: 0px;"></iframe>

config/v1beta1/types.go

With the structs added, we need to execute the `hack/update-codegen.sh` script. It will update the generated files with functions as `DeepCopy` for the added structures.

Furthermore, we will add a new function `SetDefaultNetworkTrafficArgs` in the `config/v1beta1/defaults.go`. The function will set the default values for the `NetworkInterface` and `TimeRangeInMinutes` values, but `Address` still needs to be provided.

<iframe src="https://medium.com/media/e97a10aeb9b680b636081b673ea89130" allowfullscreen="" frameborder="0" height="303" width="680" title="Network Traffic Plugin - default.go" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 303px; position: absolute; left: 0px;"></iframe>

Default values for plugin arguments

To finish the default values configuration, we need to make sure the function above is registered in the v1beta1 schema. Thus, make sure that it is registered in the file `pkg/apis/config/v1beta1/zz_generated.defaults.go`.

<iframe src="https://medium.com/media/9e8309b28b18a07b767da174ea510885" allowfullscreen="" frameborder="0" height="479" width="680" title="Network Traffic Plugin - Register Default Values function" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 479px; position: absolute; left: 0px;"></iframe>

## Registering Plugin and Configuration

Now that the arguments structure is defined, our Plugin is ready. However, we still need to register the plugin and the configuration in the scheduler framework.

The scheduler-plugins project already has a couple plugins registered which makes things a bit easier as we have examples. The registration for the plugin configuration is placed under `pkg/apis/config`. In the file `register.go` we need to add the `NetworkTrafficArgs` in the call to the `AddKnownTypes` function. The same needs to be done in the `pkg/apis/config/v1beta1/register.go` file. With both files changed, the configuration registration is done.

Next, we move to the plugin registration, which is done in the `cmd/scheduler/main.go` file. In the *main* function, the `NetworkTraffic` plugin name and constructor need to be provided as arguments to the `NewSchedulerCommand`. It should look like this:

```
command := app.NewSchedulerCommand(
  app.WithPlugin(networktraffic.Name, networktraffic.New),
)
```

Also, notice that in the in the `main.go` file we have the import of `sigs.k8s.io/scheduler-plugins/pkg/apis/config/scheme`, which initializes the scheme with all configurations we have introduced in the `pkg/apis/config` files.

With that we are done from a code perspective. The full implementation can be found [here](https://github.com/juliorenner/scheduler-plugins), it also includes a couple of unit tests, so check it out!

## Deploying and using the Plugin

Now that we have the plugin done, we can deploy it in our K8S cluster and start using it. In the scheduler-plugins repository, there is a documentation on how to do it, check it [here](https://github.com/kubernetes-sigs/scheduler-plugins/blob/master/doc/install.md). We would basically need to adapt those steps with the Plugin we just implemented.

Nonetheless, before applying the changes to the cluster, make sure that you have build the scheduler container image and pushed it to a container registry which is accessible from your kubernetes. I won't go into the details as it will differ based to the environment used. You can check the Makefile as well, as there are some commands to build and push the image and also [this development](https://github.com/kubernetes-sigs/scheduler-plugins/blob/master/doc/develop.md#how-to-build) doc may help you.

As our plugin doesn't introduce any CRD, a couple steps in the [scheduler-plugins install doc](https://github.com/kubernetes-sigs/scheduler-plugins/blob/master/doc/install.md) can be skipped. As I mentioned, I am using a cluster created with kubespray with HA. Therefore, I will need to repeat the following steps on each control plane node.

1. Log into the control plane node.
2. Backup `kube-scheduler.yaml`

```
cp /etc/kubernetes/manifests/kube-scheduler.yaml /etc/kubernetes/kube-scheduler.yaml
```

\3. Create `/etc/kubernetes/networktraffic-config.yaml` and change the values according to your environment.

<iframe src="https://medium.com/media/e0045789f9397f4527340fb7b1299da4" allowfullscreen="" frameborder="0" height="435" width="680" title="Network Traffic Plugin - Configuration" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 435px; position: absolute; left: 0px;"></iframe>

\4. Modify `/etc/kubernetes/manifests/kube-scheduler.yaml` to run scheduler-plugins with Network Traffic. The changes we have made are:

- Add the command arg `--config=/etc/kubernetes/networktraffic-config.yaml`.
- Change the `image` name.
- Add a `volume` pointing to the configuration absolute path.
- Add a `volumeMount` to make the configuration available to the scheduler pod.

Check the example below:

<iframe src="https://medium.com/media/7f634d07f1a6797273b37fb5e5aa51d9" allowfullscreen="" frameborder="0" height="1425" width="680" title="Network Traffic Plugin - kube-scheduler pod" class="es n gl dh bf" scrolling="no" style="box-sizing: inherit; top: 0px; width: 680px; height: 1425px; position: absolute; left: 0px;"></iframe>

Now, we can start taking advantage of our custom plugin. Once you check the logs of the running pod, you should see lines with the node bandwidth returned from prometheus and you can make sure the behavior is as expected. Below, we can see that `node4` correctly has the higher score, as it is the node with less network traffic:

![img](https://miro.medium.com/v2/resize:fit:700/1*Tv2hNn0vY3a_Y_2byU1v0w.png)

Hope this post is useful to you and feel free to give feedbacks on the comments, they are very appreciated!

# References

1: https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/

2: https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/

3: https://kubernetes.io/docs/reference/scheduling/config/#profiles

4: https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#extension-points

5: https://kubernetes.io/docs/reference/scheduling/config/#multiple-profiles

6: https://github.com/kubernetes/enhancements/blob/master/keps/sig-scheduling/624-scheduling-framework/README.md

7: https://github.com/kubernetes/enhancements/blob/master/keps/sig-scheduling/624-scheduling-framework/README.md#custom-scheduler-plugins-out-of-tree

8: https://github.com/kubernetes/enhancements/blob/master/keps/sig-scheduling/624-scheduling-framework/README.md#optional-args

9: https://github.com/kubernetes/enhancements/blob/master/keps/sig-scheduling/624-scheduling-framework/README.md#configuring-plugins