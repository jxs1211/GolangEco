# [Investigate EKS API & K8S API for wave test](https://confluence.ubisoft.com/pages/viewpage.action?pageId=1752175429)

## **API overview**

API operation:

- [AssociateEncryptionConfig](https://docs.aws.amazon.com/eks/latest/APIReference/API_AssociateEncryptionConfig.html)
- [AssociateIdentityProviderConfig](https://docs.aws.amazon.com/eks/latest/APIReference/API_AssociateIdentityProviderConfig.html)
- [CreateAddon](https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateAddon.html)
- [CreateCluster](https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateCluster.html)
- [CreateFargateProfile](https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateFargateProfile.html)
- [CreateNodegroup](https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateNodegroup.html)
- [DeleteAddon](https://docs.aws.amazon.com/eks/latest/APIReference/API_DeleteAddon.html)
- [DeleteCluster](https://docs.aws.amazon.com/eks/latest/APIReference/API_DeleteCluster.html)
- [DeleteFargateProfile](https://docs.aws.amazon.com/eks/latest/APIReference/API_DeleteFargateProfile.html)
- [DeleteNodegroup](https://docs.aws.amazon.com/eks/latest/APIReference/API_DeleteNodegroup.html)
- [DeregisterCluster](https://docs.aws.amazon.com/eks/latest/APIReference/API_DeregisterCluster.html)
- [DescribeAddon](https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeAddon.html)
- [DescribeAddonVersions](https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeAddonVersions.html)
- [DescribeCluster](https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeCluster.html)
- [DescribeFargateProfile](https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeFargateProfile.html)
- [DescribeIdentityProviderConfig](https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeIdentityProviderConfig.html)
- [DescribeNodegroup](https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeNodegroup.html)
- [DescribeUpdate](https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeUpdate.html)
- [DisassociateIdentityProviderConfig](https://docs.aws.amazon.com/eks/latest/APIReference/API_DisassociateIdentityProviderConfig.html)
- [ListAddons](https://docs.aws.amazon.com/eks/latest/APIReference/API_ListAddons.html)
- [ListClusters](https://docs.aws.amazon.com/eks/latest/APIReference/API_ListClusters.html)
- [ListFargateProfiles](https://docs.aws.amazon.com/eks/latest/APIReference/API_ListFargateProfiles.html)
- [ListIdentityProviderConfigs](https://docs.aws.amazon.com/eks/latest/APIReference/API_ListIdentityProviderConfigs.html)
- [ListNodegroups](https://docs.aws.amazon.com/eks/latest/APIReference/API_ListNodegroups.html)
- [ListTagsForResource](https://docs.aws.amazon.com/eks/latest/APIReference/API_ListTagsForResource.html)
- [ListUpdates](https://docs.aws.amazon.com/eks/latest/APIReference/API_ListUpdates.html)
- [RegisterCluster](https://docs.aws.amazon.com/eks/latest/APIReference/API_RegisterCluster.html)
- [TagResource](https://docs.aws.amazon.com/eks/latest/APIReference/API_TagResource.html)
- [UntagResource](https://docs.aws.amazon.com/eks/latest/APIReference/API_UntagResource.html)
- [UpdateAddon](https://docs.aws.amazon.com/eks/latest/APIReference/API_UpdateAddon.html)
- [UpdateClusterConfig](https://docs.aws.amazon.com/eks/latest/APIReference/API_UpdateClusterConfig.html)
- [UpdateClusterVersion](https://docs.aws.amazon.com/eks/latest/APIReference/API_UpdateClusterVersion.html)
- [UpdateNodegroupConfig](https://docs.aws.amazon.com/eks/latest/APIReference/API_UpdateNodegroupConfig.html)
- [UpdateNodegroupVersion](https://docs.aws.amazon.com/eks/latest/APIReference/API_UpdateNodegroupVersion.html)

Related tool utility (boto3)

- [`associate_encryption_config()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.associate_encryption_config)
- [`associate_identity_provider_config()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.associate_identity_provider_config)
- [`can_paginate()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.can_paginate)
- [`close()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.close)
- [`create_addon()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.create_addon)
- [`create_cluster()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.create_cluster)
- [`create_fargate_profile()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.create_fargate_profile)
- [`create_nodegroup()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.create_nodegroup)
- [`delete_addon()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.delete_addon)
- [`delete_cluster()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.delete_cluster)
- [`delete_fargate_profile()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.delete_fargate_profile)
- [`delete_nodegroup()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.delete_nodegroup)
- [`deregister_cluster()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.deregister_cluster)
- [`describe_addon()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.describe_addon)
- [`describe_addon_versions()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.describe_addon_versions)
- [`describe_cluster()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.describe_cluster)
- [`describe_fargate_profile()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.describe_fargate_profile)
- [`describe_identity_provider_config()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.describe_identity_provider_config)
- [`describe_nodegroup()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.describe_nodegroup)
- [`describe_update()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.describe_update)
- [`disassociate_identity_provider_config()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.disassociate_identity_provider_config)
- [`get_paginator()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.get_paginator)
- [`get_waiter()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.get_waiter)
- [`list_addons()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.list_addons)
- [`list_clusters()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.list_clusters)
- [`list_fargate_profiles()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.list_fargate_profiles)
- [`list_identity_provider_configs()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.list_identity_provider_configs)
- [`list_nodegroups()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.list_nodegroups)
- [`list_tags_for_resource()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.list_tags_for_resource)
- [`list_updates()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.list_updates)
- [`register_cluster()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.register_cluster)
- [`tag_resource()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.tag_resource)
- [`untag_resource()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.untag_resource)
- [`update_addon()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.update_addon)
- [`update_cluster_config()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.update_cluster_config)
- [`update_cluster_version()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.update_cluster_version)
- [`update_nodegroup_config()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.update_nodegroup_config)
- [`update_nodegroup_version()`](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.update_nodegroup_version)

## **Prerequisite**

- aws credential
- eksctl
- Kubernetes python client
- awscliv2
- kubecfg.yaml



**comparison table of wave load test and k8s**

| test information       | k8s objects              | comment                 |
| :--------------------- | :----------------------- | :---------------------- |
| test instance(test id) | Deployment               |                         |
| test machine           | Pod                      |                         |
| agent                  | container                |                         |
| player                 | player                   |                         |
| Machines' Geography    | Cluster's Geography      | multi-geography cluster |
| Machine Type           | Pod template             | calculation of resource |
| OS                     | os of Node/Pod/Container |                         |



**calculation of resources in k8s**

Simulated Users: 480000 = 100(player) × 10(agent)container × 480(machine)Pod

Machine Type: r5.xlarge(4c32g) Pod template

container's requests:

\- CPU: 4c / 10 = 0.4c (400m)
\- mem: 32g / 10 = 3.2g (3200Mi)

Deployment YAML:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: eks-sample-linux-deployment // test.id-linux-deployment
  namespace: eks-sample-app // eks-waveloadtest-app
  labels:
    app: eks-sample-linux-app // test.id-linux-deployment
    testId: '2e9103ae-3268-11ed-b5cf-0242ac110033 // optional {{test-id}}'
spec:
  replicas: 480 // test.configuration.test_objective.docker.slaves
  selector:
    matchLabels: null
    app: eks-sample-linux-app // test.id-linux-deployment
  template:
    metadata: null
    labels:
      app: eks-sample-linux-app // test.id-linux-deployment
      testId: 2e9103ae-3268-11ed-b5cf-0242ac110033--p1 // optional
  spec:
    containers:
      - name: All-service-ltystore-10k-p1-c1 // test.id-p1-c1
        image: 500398181577.dkr.ecr.us-east-1.amazonaws.com/linux_base_image:v0.0.2//
          test.configuration.test_agent.agent_type
        resources:
          requests:
            memory: 3200Mi // 32g / 10 (test_machine.ec2_instance_type, test_objective.process)
            cpu: 400m // 4c / 10 (test_machine.ec2_instance_type, test_objective.process)
          limits:
            memory: 3200Mi
            cpu: 400m
          env:
            - name: TEST_ID
              value: 2e9103ae-3268-11ed-b5cf-0242ac110033
            - name: TOTAL_PLAYERS
              value: 480000// test_objective.docker.expected_ccu
            - name: PLAYERS_PER_AGENT
              value: 100
        imagePullPolicy: IfNotPresent
    nodeSelector:
      kubernetes.io/os: linux // test.configuration.test_machine.ec2.os

```



## **resource API design**



**Create deployment**

**POST** /spaces/{space_id}/global/harbourapac/waveloadtest/resource/k8s/deployment/

##### Request Body:

```json
{
    "metadata": {
      "name": "eks-sample-linux-deployment",
      "namespace": "eks-sample-test",
      "labels": {
        "app": "eks-test-linux-app-deployment",
        "testId": "2e9103ae-3268-11ed-b5cf-0242ac110033"
      }
    },
    "spec": {
      "replicas": 480,
      "selector": {
        "matchLabels": {
          "app": "eks-sample-linux-app"
        }
      },
      "template": {
        "metadata": {
          "labels": {
            "app": "eks-sample-linux-app-pod",
            "testId": "2e9103ae-3268-11ed-b5cf-0242ac110033--p1"
          }
        },
        "spec": {
          "containers": [
            {
              "name": "All-service-ltystore-10k-p1-c1",
              "image": "500398181577.dkr.ecr.us-east-1.amazonaws.com/linux_base_image:v0.0.2",
              "env": [
                {
                  "name": "TEST_ID",
                  "value": "2e9103ae-3268-11ed-b5cf-0242ac110033"
                },
                {
                  "name": "TOTAL_PLAYERS",
                  "value": "480000"
                },
                {
                  "name": "PLAYERS_PER_AGENT",
                  "value": "100"
                }
              ],
              "resources": {
                "limits": {
                  "cpu": "400m",
                  "memory": "3200Mi"
                },
                "requests": {
                  "cpu": "400m",
                  "memory": "3200Mi"
                }
              }
            }
          ],
          "nodeSelector": {
            "kubernetes.io/os": "linux"
          },
        }
      }
    }
  }
```

##### Response Body:

**creat deployment** Expand source

```json
{
    "namespace": "wave-linux-test",
    "name": "eks-sample-linux-deployment",
    "generation": 1,
    "image": "500398181577.dkr.ecr.us-east-1.amazonaws.com/linux_base_image:v0.0.2",
}
```





**Delete deployment**

**DELETE** /spaces/{space_id}/global/harbourapac/waveloadtest/resource/k8s/deployment/

##### Request Body:

 Expand source

```json
{
    "namespace": "wave-linux-test",
    "name": "eks-sample-linux-deployment",
}
```

##### Response Body:

 Expand source

```json
{
    "result": "ok"
}
```



**Create cluster**

**CREATE** /spaces/{space_id}/global/harbourapac/waveloadtest/resource/k8s/cluster/

##### Request Body:

 Expand source

```json
{
    name="string",
    roleArn="string",
    resourcesVpcConfig={
        "subnetIds": [
            "string",
        ],
        "securityGroupIds": [
            "string",
        ],
        "endpointPublicAccess": True|False,
        "endpointPrivateAccess": True|False,
        "publicAccessCidrs": [
            "string",
        ]
    }
}
```



##### Response Body:

 Expand source

```json
{
    "cluster": {
        "name": "string",
        "arn": "string",
        "createdAt": datetime(2015, 1, 1),
        "version": "string",
        "endpoint": "string",
        "roleArn": "string",
        "resourcesVpcConfig": {
            "subnetIds": [
                "string",
            ],
            "securityGroupIds": [
                "string",
            ],
            "clusterSecurityGroupId": "string",
            "vpcId": "string",
            "endpointPublicAccess": True|False,
            "endpointPrivateAccess": True|False,
            "publicAccessCidrs": [
                "string",
            ]
        },
        "kubernetesNetworkConfig": {
            "serviceIpv4Cidr": "string",
            "serviceIpv6Cidr": "string",
            "ipFamily": "ipv4"|"ipv6"
        },
        "logging": {
            "clusterLogging": [
                {
                    "types": [
                        "api"|"audit"|"authenticator"|"controllerManager"|"scheduler",
                    ],
                    "enabled": True|False
                },
            ]
        },
        "identity": {
            "oidc": {
                "issuer": "string"
            }
        },
        "status": "CREATING"|"ACTIVE"|"DELETING"|"FAILED"|"UPDATING"|"PENDING",
        "certificateAuthority": {
            "data": "string"
        },
        "clientRequestToken": "string",
        "platformVersion": "string",
        "tags": {
            "string": "string"
        },
        "encryptionConfig": [
            {
                "resources": [
                    "string",
                ],
                "provider": {
                    "keyArn": "string"
                }
            },
        ],
        "connectorConfig": {
            "activationId": "string",
            "activationCode": "string",
            "activationExpiry": datetime(2015, 1, 1),
            "provider": "string",
            "roleArn": "string"
        },
        "id": "string",
        "health": {
            "issues": [
                {
                    "code": "AccessDenied"|"ClusterUnreachable"|"ConfigurationConflict"|"InternalFailure"|"ResourceLimitExceeded"|"ResourceNotFound",
                    "message": "string",
                    "resourceIds": [
                        "string",
                    ]
                },
            ]
        },
        "outpostConfig": {
            "outpostArns": [
                "string",
            ],
            "controlPlaneInstanceType": "string"
        }
    }
}
```



**Delete cluster**

**DELETE** /spaces/{space_id}/global/harbourapac/waveloadtest/resource/k8s/cluster/

##### Request Body:

 Expand source

```json
{
    name="string"
}
```

##### Response Body:

 Expand source

```json
{
    "cluster": {
        "name": "string",
        "arn": "string",
        "createdAt": datetime(2015, 1, 1),
        "version": "string",
        "endpoint": "string",
        "roleArn": "string",
        "resourcesVpcConfig": {
            "subnetIds": [
                "string",
            ],
            "securityGroupIds": [
                "string",
            ],
            "clusterSecurityGroupId": "string",
            "vpcId": "string",
            "endpointPublicAccess": True|False,
            "endpointPrivateAccess": True|False,
            "publicAccessCidrs": [
                "string",
            ]
        },
        "kubernetesNetworkConfig": {
            "serviceIpv4Cidr": "string",
            "serviceIpv6Cidr": "string",
            "ipFamily": "ipv4"|"ipv6"
        },
        "logging": {
            "clusterLogging": [
                {
                    "types": [
                        "api"|"audit"|"authenticator"|"controllerManager"|"scheduler",
                    ],
                    "enabled": True|False
                },
            ]
        },
        "identity": {
            "oidc": {
                "issuer": "string"
            }
        },
        "status": "CREATING"|"ACTIVE"|"DELETING"|"FAILED"|"UPDATING"|"PENDING",
        "certificateAuthority": {
            "data": "string"
        },
        "clientRequestToken": "string",
        "platformVersion": "string",
        "tags": {
            "string": "string"
        },
        "encryptionConfig": [
            {
                "resources": [
                    "string",
                ],
                "provider": {
                    "keyArn": "string"
                }
            },
        ],
        "connectorConfig": {
            "activationId": "string",
            "activationCode": "string",
            "activationExpiry": datetime(2015, 1, 1),
            "provider": "string",
            "roleArn": "string"
        },
        "id": "string",
        "health": {
            "issues": [
                {
                    "code": "AccessDenied"|"ClusterUnreachable"|"ConfigurationConflict"|"InternalFailure"|"ResourceLimitExceeded"|"ResourceNotFound",
                    "message": "string",
                    "resourceIds": [
                        "string",
                    ]
                },
            ]
        },
        "outpostConfig": {
            "outpostArns": [
                "string",
            ],
            "controlPlaneInstanceType": "string"
        }
    }
}
```



**List cluster**

**GET** /spaces/{space_id}/global/harbourapac/waveloadtest/resource/k8s/cluster/

##### Request Body:

 Expand source

```json
{
    "maxResults"=123,
    "nextToken"="string",
    "include"=[
        "string",
    ]
}
```

##### Response Body:

 Expand source

```json
{
    "clusters": [
        "string",
    ],
    "nextToken": "string"
}
```



**Update cluster**

**PUT** /spaces/{space_id}/global/harbourapac/waveloadtest/resource/k8s/cluster/

##### Request Body:

```json
{
    "name": "string",
     "resourcesVpcConfig": {
         "subnetIds": [
             "string",
         ],
         "securityGroupIds": [
             "string",
         ],
         "endpointPublicAccess": True|False,
         "endpointPrivateAccess": True|False,
         "publicAccessCidrs": [
             "string",
         ]
     },
     "logging": {
         "clusterLogging": [
             {
                 "types": [
                     "api"|"audit"|"authenticator"|"controllerManager"|"scheduler",
                 ],
                 "enabled": True|False
             },
         ]
     },
     "clientRequestToken": "string"
 }
```

##### Response Body:

```json
{
    "update": {
        "id": "string",
        "status": "InProgress"|"Failed"|"Cancelled"|"Successful",
        "type": "VersionUpdate"|"EndpointAccessUpdate"|"LoggingUpdate"|"ConfigUpdate"|"AssociateIdentityProviderConfig"|"DisassociateIdentityProviderConfig"|"AssociateEncryptionConfig"|"AddonUpdate",
        "params": [
            {
                "type": "Version"|"PlatformVersion"|"EndpointPrivateAccess"|"EndpointPublicAccess"|"ClusterLogging"|"DesiredSize"|"LabelsToAdd"|"LabelsToRemove"|"TaintsToAdd"|"TaintsToRemove"|"MaxSize"|"MinSize"|"ReleaseVersion"|"PublicAccessCidrs"|"LaunchTemplateName"|"LaunchTemplateVersion"|"IdentityProviderConfig"|"EncryptionConfig"|"AddonVersion"|"ServiceAccountRoleArn"|"ResolveConflicts"|"MaxUnavailable"|"MaxUnavailablePercentage",
                "value": "string"
            },
        ],
        "createdAt": datetime(2015, 1, 1),
        "errors": [
            {
                "errorCode": "SubnetNotFound"|"SecurityGroupNotFound"|"EniLimitReached"|"IpNotAvailable"|"AccessDenied"|"OperationNotPermitted"|"VpcIdNotFound"|"Unknown"|"NodeCreationFailure"|"PodEvictionFailure"|"InsufficientFreeAddresses"|"ClusterUnreachable"|"InsufficientNumberOfReplicas"|"ConfigurationConflict"|"AdmissionRequestDenied"|"UnsupportedAddonModification"|"K8sResourceNotFound",
                "errorMessage": "string",
                "resourceIds": [
                    "string",
                ]
            },
        ]
    }
}
```



**Demonstration:**

operate k8s deployment:

![img](https://confluence.ubisoft.com/download/attachments/1752175429/image2022-9-14_7-53-7.png?version=1&modificationDate=1663113188000&api=v2)

operate eks cluster:

![img](https://confluence.ubisoft.com/download/attachments/1752175429/image2022-9-14_14-8-53.png?version=1&modificationDate=1663135734000&api=v2)

![img](https://confluence.ubisoft.com/download/attachments/1752175429/image2022-9-14_14-10-17.png?version=1&modificationDate=1663135818000&api=v2)



Something needs to consider:

Comparing to creating cluster by eksctl, there are many resources the cluster is dependent on that have to be created on our behalf in advance by using boto3 API.

These resources include stack for cluster control plane and node group, VPCs, subnets, security groups, route tables, routes, and internet and NAT gateways.



**Other need to be investigated and designed in the future:**

- how does other service call this deployment API
- container start policy, Pods communication, broadcast to Pods
- business data transfer implement at resource client
- weather single cluster control plane can manage both Linux and Windows Node

- how does container access S3

- - transfer script into Pod
  - gather output generated by the test
- storage
  - business-related data
  - DB design, service API design(model, serializer, view)

- observability

- - the running process of the user script
  - WMS-agent
  - RPC-server



**References:**

https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html

https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html

https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/eks.html#EKS.Client.list_clusters

https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html

https://docs.aws.amazon.com/eks/latest/userguide/sample-deployment.html

https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/

https://docs.aws.amazon.com/eks/latest/userguide/eks-managing.html

https://aws.amazon.com/cn/blogs/containers/operating-a-multi-regional-stateless-application-using-amazon-eks/

https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/

https://docs.aws.amazon.com/eks/latest/userguide/windows-support.html

https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/



Operate k8s by python sdk

https://github.com/kubernetes-client/python

https://leftasexercise.com/2019/04/01/python-up-an-eks-cluster-part-i/

https://leftasexercise.com/2019/04/15/python-up-an-eks-cluster-part-ii/

https://leftasexercise.com/2019/04/25/kubernetes-101-creating-pods-and-deployments/



- Resource API 文档（集群相关的API以及deployment相关的API，包括参数描述，返回值，错误码，异常情况等）
- 数据库设计



