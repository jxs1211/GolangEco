![image-20230612101712798](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230612101712798.png)

Architecture:

- Master Node:
  - API Server
  - Controller Manager
  - Scheduler
  - ETCD
- Worker Node:
  - Kubelet
  - Kube-proxy



![image-20230612102310311](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230612102310311.png)

Container in Pod:

- Localhost interface: all containers share the same network namespace and use the same IP address and communicate with each other 
- Individual network interfaces: created by Kubernetes CNI to communicate with other Pod and service in cluster





![image-20230612102548298](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230612102548298.png)

service is a common communication way between Pods

- Permanent IP
- Load balancer

![image-20230612103815123](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230612103815123.png)





![image-20230612104024464](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230612104024464.png)



![image-20230612104654250](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230612104654250.png)



Stateful application should use StatefuleSet but not deployment

![image-20230612105101996](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230612105101996.png)



![image-20230612105748794](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230612105748794.png)

![image-20230614095347863](C:\Users\xjshen\AppData\Roaming\Typora\typora-user-images\image-20230614095347863.png)

