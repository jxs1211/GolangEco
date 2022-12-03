13.2. Mock Exam -2 (Solutiohortal + master $ master $ 

```sh
kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec. taints}{"\n"}{end}'

 master [map[effect:Noschedule key:node-role.kubernetes.io/master]] node01 [map[effect:NoSchedule key: app_type value:alpha]] node02 node03 master $ So now when you run, those who can see it clearly that the muzzle has changed and Nordstrom has a tent, K 07:23 / 32:22 A 1 0939 NEST LX > 发送 720P ** ** 1) 
```


![image-20221112094709147](images/image-20221112094709147.png)

```sh
kubectl describe deployments.apps frontend-deployment | grep -i image
```

![image-20221112100711270](images/image-20221112100711270.png)

```sh
kubectl create deployment name --image=nginx
kubectl scale deployment name --replicas=3
```

![image-20221112101054940](images/image-20221112101054940.png)

![image-20221112101618979](images/image-20221112101618979.png)

![image-20221112101557227](images/image-20221112101557227.png)

![image-20221112101841284](images/image-20221112101841284.png)

![ ](images/image-20221112102252538.png)

![image-20221112105226902](images/image-20221112105226902.png)



![image-20221112111349829](images/image-20221112111349829.png)

```sh
kubectl run redis --image=redis --dry-run=client -o yaml > pod.yaml
```

![image-20221112111616976](images/image-20221112111616976.png)

```sh
kubectl expose pod redis --name redis-service --port 6379 --target-port 6379
```

![image-20221112112749153](images/image-20221112112749153.png)

```sh
kubectl describe svc redis-service
```



![image-20221112112811709](images/image-20221112112811709.png)



![image-20221112102555170](images/image-20221112102555170.png)



```sh
kubectl run redis --image=redis --dry-run=client -o yaml > pod.yaml
kubectl apply -f pod.yaml
kubectl edit pod redis
```



![image-20221112163430250](images/image-20221112163430250.png)

![image-20221112164053610](images/image-20221112164053610.png)

```sh
kubectl create -f replicaset-definition.yaml
kubectl get replicaset
kubectl get pods
```



![image-20221112165105959](images/image-20221112165105959.png)

![image-20221112165301053](images/image-20221112165301053.png)



#### how scale:

![image-20221112165831474](images/image-20221112165831474.png)

![image-20221113114554134](images/image-20221113114554134.png)

![image-20221113115835384](images/image-20221113115835384.png)

![image-20221113120104846](images/image-20221113120104846.png)

```sh
kubectl create serviceaccount dashboard-sa
kubectl get serviceaccount
kubectl describe serviceaccount dashboard-sa
kubectl describe secret dashboard-sa
```



![image-20221113120247661](images/image-20221113120247661.png)

![image-20221113120506133](images/image-20221113120506133.png)

![image-20221113120836327](images/image-20221113120836327.png)

![image-20221113121405163](images/image-20221113121405163.png)

![image-20221113121434688](images/image-20221113121434688.png)

![image-20221113122735731](images/image-20221113122735731.png)

![image-20221113123659216](images/image-20221113123659216.png)

![image-20221114084631767](images/image-20221114084631767.png)

![image-20221114084910670](images/image-20221114084910670.png)

![image-20221114085105309](images/image-20221114085105309.png)

```sh
kubectl taint nodes node1 app=blue:NoSchedule
kubectl taint nodes node1 key1=value1:NoSchedule-
```



![image-20221114085335986](images/image-20221114085335986.png)

![image-20221114090920986](images/image-20221114090920986.png)

![image-20221114090950385](images/image-20221114090950385.png)

```sh
kubectl explain pod --recursive | less
kubectl explain pod --recursive | grep tolerations -A5 
```

![image-20221114092345950](images/image-20221114092345950.png)

![image-20221114092315718](images/image-20221114092315718.png)



![image-20221114084754287](images/image-20221114084754287.png)

```sh
kubectl describe node master
kubectl taint node master xxx:NoSchedule-
```

![image-20221114092711371](images/image-20221114092711371.png)

configmap

![image-20221114092931791](images/image-20221114092931791.png)



![image-20221114095023913](images/image-20221114095023913.png)

node

```sh
kubectl label nodes foo unhealthy=true
```

![image-20221114095445750](images/image-20221114095445750.png)

![image-20221114095502627](images/image-20221114095502627.png)

![image-20221114100552880](images/image-20221114100552880.png)

![image-20221114100612412](images/image-20221114100612412.png)

![image-20221114101250091](images/image-20221114101250091.png)

![image-20221114101423884](images/image-20221114101423884.png)

![image-20221114101638484](images/image-20221114101638484.png)

```

```

![image-20221114132628123](images/image-20221114132628123.png)

![image-20221114133849596](images/image-20221114133849596.png)

Secret

![image-20221114134847530](images/image-20221114134847530.png)

![image-20221114135051108](images/image-20221114135051108.png)

![image-20221114135125154](images/image-20221114135125154.png)

![image-20221114135613667](images/image-20221114135613667.png)

![image-20221114135849605](images/image-20221114135849605.png)

![image-20221115082320677](images/image-20221115082320677.png)

taints and tolerations vs. node affinity

![image-20221115090248323](images/image-20221115090248323.png)

![image-20221115090540875](images/image-20221115090540875.png)

   ![image-20221115090705887](images/image-20221115090705887.png)

security context
Linux Capability

![image-20221115091937458](images/image-20221115091937458.png)

![image-20221115092007874](images/image-20221115092007874.png)

![image-20221115092025693](images/image-20221115092025693.png)

multi-container

![image-20221115092351055](images/image-20221115092351055.png)

![image-20221115092733703](images/image-20221115092733703.png)

![image-20221115092853441](images/image-20221115092853441.png)

![image-20221115093127068](images/image-20221115093127068.png)

![image-20221115093014757](images/image-20221115093014757.png)

![image-20221116091909002](images/image-20221116091909002.png)

![image-20221116092136852](images/image-20221116092136852.png)![image-20221116092825223](images/image-20221116092825223.png)

observability![image-20221116193958988](images/image-20221116193958988.png)

![image-20221116194641311](images/image-20221116194641311.png)

![image-20221116195018645](images/image-20221116195018645.png)

![image-20221116195527886](images/image-20221116195527886.png)

![image-20221116195647296](images/image-20221116195647296.png)

metric server

![image-20221117084414590](images/image-20221117084414590.png)

Pod design: labels, selector and annotation

![image-20221117085247553](images/image-20221117085247553.png)

![image-20221117085837399](images/image-20221117085837399.png)

![image-20221117085952860](images/image-20221117085952860.png)

![image-20221117090021891](images/image-20221117090021891.png)

![image-20221117090149635](images/image-20221117090149635.png)

Rollout

![image-20221117092510877](images/image-20221117092510877.png)

![image-20221117092635354](images/image-20221117092635354.png)

 ![image-20221117092820845](images/image-20221117092820845.png)

![image-20221117092931309](images/image-20221117092931309.png)

![image-20221117093848649](images/image-20221117093848649.png)

![image-20221117094001660](images/image-20221117094001660.png)

![image-20221117094127243](images/image-20221117094127243.png)

deployment

revision 1

![image-20221117170059541](images/image-20221117170059541.png)

revision2: modify the nginx version and apply again

![image-20221117174901752](images/image-20221117174901752.png)

revision3: modify the version of nginx by the cmd

![image-20221117174951336](images/image-20221117174951336.png)

revision4: undo to take us back to the deployment  of revision 2 that we applied at that time and apply it again.

![image-20221117175525830](images/image-20221117175525830.png)

revision5: modify the version of nginx to a error version

revison6: undo the deployement to rollback to former deployment of revision5 which actually use the deployment of revision4

![image-20221117180237662](images/image-20221117180237662.png)

Job and cron job

![image-20221117200640753](images/image-20221117200640753.png)

![image-20221117200742972](images/image-20221117200742972.png)

![image-20221117200850231](images/image-20221117200850231.png)

![image-20221117201149019](images/image-20221117201149019.png)

![image-20221117201428176](images/image-20221117201428176.png)

![image-20221117201545678](images/image-20221117201545678.png)

![image-20221117201649380](images/image-20221117201649380.png)

![image-20221117202145145](images/image-20221117202145145.png)

![image-20221117201932088](images/image-20221117201932088.png)

![image-20221117202829777](images/image-20221117202829777.png)

![image-20221117203041039](images/image-20221117203041039.png)

![image-20221117203127862](images/image-20221117203127862.png)

![image-20221117203233930](images/image-20221117203233930.png)

![image-20221117203913086](images/image-20221117203913086.png)

![image-20221117203649925](images/image-20221117203649925.png)

![image-20221117204303882](images/image-20221117204303882.png)

![image-20221117204445582](images/image-20221117204445582.png)

![image-20221117204417998](images/image-20221117204417998.png)

Services

![image-20221118082856967](images/image-20221118082856967.png)

![image-20221118083108761](images/image-20221118083108761.png)

![image-20221118083428835](images/image-20221118083428835.png)

![image-20221118083604939](images/image-20221118083604939.png)

![image-20221118083713805](images/image-20221118083713805.png)

![image-20221118083841008](images/image-20221118083841008.png)

![image-20221118083956289](images/image-20221118083956289.png)

![image-20221118084331808](images/image-20221118084331808.png)

![image-20221118084413521](images/image-20221118084413521.png)

![image-20221118084538906](images/image-20221118084538906.png)

![image-20221118084654939](images/image-20221118084654939.png)

![image-20221118085051950](images/image-20221118085051950.png)

![image-20221118085300154](images/image-20221118085300154.png)

![image-20221118085633608](images/image-20221118085633608.png)

![image-20221118085758990](images/image-20221118085758990.png)

![image-20221118090102066](images/image-20221118090102066.png)

![image-20221118090206882](images/image-20221118090206882.png)



![image-20221118090455399](images/image-20221118090455399.png)

![image-20221118090534340](images/image-20221118090534340.png)

![image-20221118090616853](images/image-20221118090616853.png)

![image-20221118090811921](images/image-20221118090811921.png)

![image-20221118091608197](images/image-20221118091608197.png)

![image-20221118091724952](images/image-20221118091724952.png)

![image-20221118092621568](images/image-20221118092621568.png)

![image-20221118092522931](images/image-20221118092522931.png)

Ingress

![image-20221118144658844](images/image-20221118144658844.png)

Why bring in a additional layer called proxy

![image-20221118144414921](images/image-20221118144414921.png)

![image-20221118143145113](images/image-20221118143145113.png)

![image-20221118145049440](images/image-20221118145049440.png)

set the DNS to point to the ip of the LB

![image-20221118145243494](images/image-20221118145243494.png)

![image-20221118145543904](images/image-20221118145543904.png)

![image-20221118150231823](images/image-20221118150231823.png)

![image-20221118150328653](images/image-20221118150328653.png)

![image-20221118150553886](images/image-20221118150553886.png)

![image-20221118150923861](images/image-20221118150923861.png)

![image-20221118165609999](images/image-20221118165609999.png)

![image-20221118165850717](images/image-20221118165850717.png)

![image-20221118171233724](images/image-20221118171233724.png)

![image-20221118172320437](images/image-20221118172320437.png)

![image-20221118173507011](images/image-20221118173507011.png)

![image-20221118173648533](images/image-20221118173648533.png)

![image-20221118173857883](images/image-20221118173857883.png)

![image-20221118195208262](images/image-20221118195208262.png)

![image-20221119102016796](images/image-20221119102016796.png)

![image-20221119102139982](images/image-20221119102139982.png)

![image-20221119102334644](images/image-20221119102334644.png)

![image-20221119102415911](images/image-20221119102415911.png)

![image-20221119102544269](images/image-20221119102544269.png)

![image-20221119102918542](images/image-20221119102918542.png)

![image-20221119103039440](images/image-20221119103039440.png)

![image-20221119103452950](images/image-20221119103452950.png)

![image-20221119103605782](images/image-20221119103605782.png)

![image-20221119103748142](images/image-20221119103748142.png)

Network Policy

![image-20221119104016829](images/image-20221119104016829.png)

![image-20221119104641540](images/image-20221119104641540.png)

![image-20221119104601050](images/image-20221119104601050.png)

![image-20221119104823305](images/image-20221119104823305.png)

![image-20221119110146033](images/image-20221119110146033.png)

![image-20221119105613258](images/image-20221119105613258.png)

![image-20221119105836950](images/image-20221119105836950.png)

![image-20221119110551798](images/image-20221119110551798.png)

![image-20221119110804699](images/image-20221119110804699.png)

![image-20221119111535215](images/image-20221119111535215.png)

![image-20221119111057027](images/image-20221119111057027.png)

![image-20221119111309130](images/image-20221119111309130.png)

![image-20221119111731727](images/image-20221119111731727.png)

![image-20221119111901195](images/image-20221119111901195.png)

![image-20221119112256464](images/image-20221119112256464.png)

![image-20221119113145192](images/image-20221119113145192.png)

Persist Volume

![image-20221119113725472](images/image-20221119113725472.png)

![image-20221119114211752](images/image-20221119114211752.png)

![image-20221119114353146](images/image-20221119114353146.png)

![image-20221119115036857](images/image-20221119115036857.png)

![image-20221119115218444](images/image-20221119115218444.png)



![image-20221119115417544](images/image-20221119115417544.png)

![image-20221119115635782](images/image-20221119115635782.png)

![image-20221119115859511](images/image-20221119115859511.png)

![image-20221119115933485](images/image-20221119115933485.png)

![image-20221119115954224](images/image-20221119115954224.png)

![image-20221119120213004](images/image-20221119120213004.png)

![image-20221119120328365](images/image-20221119120328365.png)

In Retain mode, the PV will both not be deleted and be reused by other pvc while the former bound pvc was deleted.

![image-20221119120414160](images/image-20221119120414160.png)

In Delete mode, the PV will both be deleted while the former bound pvc was deleted and thus the storage would be freeing up.

![image-20221119120629554](images/image-20221119120629554.png)

![image-20221119121452678](images/image-20221119121452678.png)

![image-20221119121713790](images/image-20221119121713790.png)

![image-20221119122458001](images/image-20221119122458001.png)

![image-20221119122714976](images/image-20221119122714976.png)

![image-20221119123235689](images/image-20221119123235689.png)

![image-20221119123544247](images/image-20221119123544247.png)

![image-20221119124025749](images/image-20221119124025749.png)

![image-20221119124133369](images/image-20221119124133369.png)



![image-20221119124515383](images/image-20221119124515383.png)

![image-20221119124707541](images/image-20221119124707541.png)

![image-20221119124804228](images/image-20221119124804228.png)

![image-20221119124834676](images/image-20221119124834676.png)

![image-20221119125325907](images/image-20221119125325907.png)

Pod --> volumes--> pvc--->pv--->storage

storage class

![image-20221120103406816](images/image-20221120103406816.png)

![image-20221120103612487](images/image-20221120103612487.png)

![image-20221120103804447](images/image-20221120103804447.png)

![image-20221120103933671](images/image-20221120103933671.png)

![image-20221120103947939](images/image-20221120103947939.png)

![image-20221120104307576](images/image-20221120104307576.png)

![image-20221120104607432](images/image-20221120104607432.png)

![image-20221120104941280](images/image-20221120104941280.png)

https://kubernetes.io/docs/concepts/storage/storage-classes/

https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html

Stateful sets

![image-20221120110506925](images/image-20221120110506925.png)

![image-20221120111045427](images/image-20221120111045427.png)

![image-20221120111237133](images/image-20221120111237133.png)

![image-20221120111447007](images/image-20221120111447007.png)

![image-20221120111555468](images/image-20221120111555468.png)

![image-20221120111608972](images/image-20221120111608972.png)

![image-20221120111808403](images/image-20221120111808403.png)

![image-20221120111940008](images/image-20221120111940008.png)

![image-20221120112444651](images/image-20221120112444651.png)

![image-20221120112511025](images/image-20221120112511025.png)

![image-20221120112536407](images/image-20221120112536407.png)

Headless service

![image-20221120113058482](images/image-20221120113058482.png)

![image-20221120113143917](images/image-20221120113143917.png)

![image-20221120113223611](images/image-20221120113223611.png)

![image-20221120113326392](images/image-20221120113326392.png)

![image-20221120113453791](images/image-20221120113453791.png)

![image-20221120113547377](images/image-20221120113547377.png)

![image-20221120113646815](images/image-20221120113646815.png)

![image-20221120113710939](images/image-20221120113710939.png)

![image-20221120113904858](images/image-20221120113904858.png)

![image-20221120114206608](images/image-20221120114206608.png)

![image-20221120114312126](images/image-20221120114312126.png)

![image-20221120114446638](images/image-20221120114446638.png)

![image-20221120114535862](images/image-20221120114535862.png)

![image-20221120114552646](images/image-20221120114552646.png)

![image-20221120114630051](images/image-20221120114630051.png)

![image-20221120114705189](images/image-20221120114705189.png)

![image-20221120114854287](images/image-20221120114854287.png)

![image-20221120115548113](images/image-20221120115548113.png)

![image-20221120115657522](images/image-20221120115657522.png)

Storage In stateful sets

![image-20221121084247508](images/image-20221121084247508.png)

![image-20221121084458037](images/image-20221121084458037.png)

![image-20221121084922102](images/image-20221121084922102.png)

![image-20221121085011432](images/image-20221121085011432.png)

![image-20221121085202970](images/image-20221121085202970.png)

![image-20221121085347641](images/image-20221121085347641.png)

Cluster Role

![image-20221121091316044](images/image-20221121091316044.png)

![image-20221121091350865](images/image-20221121091350865.png)

![image-20221121092107974](images/image-20221121092107974.png)

![image-20221121092129373](images/image-20221121092129373.png)



![image-20221121091620120](images/image-20221121091620120.png)

![image-20221121091706787](images/image-20221121091706787.png)

![image-20221121092341117](images/image-20221121092341117.png)

![image-20221121092426579](images/image-20221121092426579.png)

![image-20221121092758896](images/image-20221121092758896.png)

![image-20221121092839139](images/image-20221121092839139.png)

![image-20221121093105520](images/image-20221121093105520.png)

![image-20221121094006059](images/image-20221121094006059.png)

![image-20221121094204034](images/image-20221121094204034.png)

![image-20221121094708244](images/image-20221121094708244.png)

![image-20221121094729684](images/image-20221121094729684.png)

![image-20221121094847482](images/image-20221121094847482.png)

![image-20221121094834731](images/image-20221121094834731.png)

Define, build, modify docker image

![image-20221121145311785](images/image-20221121145311785.png)

![image-20221121150005152](images/image-20221121150005152.png)

Authentication, Authorization and Admission Control

![image-20221122082948590](images/image-20221122082948590.png)

Validating and mutating Admission Controller

![image-20221122083750437](images/image-20221122083750437.png)

![image-20221122083811664](images/image-20221122083811664.png)

![image-20221122084011698](images/image-20221122084011698.png)

![image-20221122084415341](images/image-20221122084415341.png)

![image-20221122084745858](images/image-20221122084745858.png)

![image-20221122084811609](images/image-20221122084811609.png)

![image-20221123083852045](images/image-20221123083852045.png)

![image-20221123083929136](images/image-20221123083929136.png)

![image-20221123084038109](images/image-20221123084038109.png)

![image-20221123084220981](images/image-20221123084220981.png)

![image-20221123084400333](images/image-20221123084400333.png)

![image-20221123084413057](images/image-20221123084413057.png)

![image-20221123084507276](images/image-20221123084507276.png)

 ![image-20221123084700010](images/image-20221123084700010.png)

Authentication

![image-20221123085505369](images/image-20221123085505369.png)

![image-20221123085556791](images/image-20221123085556791.png)

![image-20221123085637674](images/image-20221123085637674.png)

![image-20221123085719083](images/image-20221123085719083.png)

![image-20221123085902294](images/image-20221123085902294.png)

![image-20221123090102575](images/image-20221123090102575.png)

![image-20221123090306279](images/image-20221123090306279.png)

![image-20221123090454534](images/image-20221123090454534.png)

![image-20221123090620234](images/image-20221123090620234.png)

![image-20221123090931015](images/image-20221123090931015.png)

Authorization

![image-20221123091138079](images/image-20221123091138079.png)

![image-20221123091631363](images/image-20221123091631363.png)

![image-20221123091922250](images/image-20221123091922250.png)

![image-20221123092150097](images/image-20221123092150097.png)

![image-20221123092208605](images/image-20221123092208605.png)

![image-20221123092313671](images/image-20221123092313671.png)

![image-20221123092348634](images/image-20221123092348634.png)

![image-20221123092436273](images/image-20221123092436273.png)

![image-20221123092504323](images/image-20221123092504323.png)

![image-20221123092634982](images/image-20221123092634982.png)

![image-20221123092942995](images/image-20221123092942995.png)

![image-20221123092658540](images/image-20221123092658540.png)

![image-20221123092603055](images/image-20221123092603055.png)

![image-20221123092818708](images/image-20221123092818708.png)

![image-20221123092754310](images/image-20221123092754310.png)

CSR

![image-20221126101605459](images/image-20221126101605459.png)

![image-20221126101707998](images/image-20221126101707998.png)

![image-20221126101752698](images/image-20221126101752698.png)

![image-20221126102714751](images/image-20221126102714751.png)

![image-20221126103336520](images/image-20221126103336520.png)

Security Kubeconfig

![image-20221126103812155](images/image-20221126103812155.png)

![image-20221126103849711](images/image-20221126103849711.png)

![image-20221126104046630](images/image-20221126104046630.png)

![image-20221126104347645](images/image-20221126104347645.png)

![image-20221126104705440](images/image-20221126104705440.png)

![image-20221126104757071](images/image-20221126104757071.png)

```sh
[going@dev config]$ k config --help
Modify kubeconfig files using subcommands like "kubectl config set current-context my-context"

 The loading order follows these rules:

  1.  If the --kubeconfig flag is set, then only that file is loaded. The flag may only be set once and no merging takes
place.
  2.  If $KUBECONFIG environment variable is set, then it is used as a list of paths (normal path delimiting rules for
your system). These paths are merged. When a value is modified, it is modified in the file that defines the stanza. When
a value is created, it is created in the first file that exists. If no files in the chain exist, then it creates the
last file in the list.
  3.  Otherwise, ${HOME}/.kube/config is used and no merging takes place.

Available Commands:
  current-context Display the current-context
  delete-cluster  Delete the specified cluster from the kubeconfig
  delete-context  Delete the specified context from the kubeconfig
  delete-user     Delete the specified user from the kubeconfig
  get-clusters    Display clusters defined in the kubeconfig
  get-contexts    Describe one or many contexts
  get-users       Display users defined in the kubeconfig
  rename-context  Rename a context from the kubeconfig file
  set             Set an individual value in a kubeconfig file
  set-cluster     Set a cluster entry in kubeconfig
  set-context     Set a context entry in kubeconfig
  set-credentials Set a user entry in kubeconfig
  unset           Unset an individual value in a kubeconfig file
  use-context     Set the current-context in a kubeconfig file
  view            Display merged kubeconfig settings or a specified kubeconfig file

Usage:
  kubectl config SUBCOMMAND [options]

Use "kubectl <command> --help" for more information about a given command.
Use "kubectl options" for a list of global command-line options (applies to all commands).
```

![image-20221126105122752](images/image-20221126105122752.png)

![image-20221126105210591](images/image-20221126105210591.png)

![image-20221126105449414](images/image-20221126105449414.png)

![image-20221126105424645](images/image-20221126105424645.png)

![image-20221126105600674](images/image-20221126105600674.png)

![image-20221126105923014](images/image-20221126105923014.png)

![image-20221126105937571](images/image-20221126105937571.png)

![image-20221126110131871](images/image-20221126110131871.png)

![image-20221126110517047](images/image-20221126110517047.png)

https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/

![image-20221126110707055](images/image-20221126110707055.png)

![image-20221126111115427](images/image-20221126111115427.png)

![image-20221126111159903](images/image-20221126111159903.png)

![image-20221126111335542](images/image-20221126111335542.png)

![image-20221126111558818](images/image-20221126111558818.png)

Authorization

![image-20221127105221450](images/image-20221127105221450.png)

![image-20221127105430256](images/image-20221127105430256.png)

![image-20221127105653835](images/image-20221127105653835.png)

![image-20221127105956975](images/image-20221127105956975.png)

![image-20221127110013949](images/image-20221127110013949.png)

![image-20221127110215117](images/image-20221127110215117.png)

![image-20221127110324056](images/image-20221127110324056.png)

![image-20221127110622201](images/image-20221127110622201.png)

![image-20221127110921432](images/image-20221127110921432.png)

![image-20221127110945486](images/image-20221127110945486.png)

![image-20221127111107487](images/image-20221127111107487.png)

![image-20221127111200919](images/image-20221127111200919.png)

Time management

1. Attempt all questions
2. Don't get stuck in any question even it's easy, don't spend too much time on it, mark it down and skip to the next one
3. Do all the trouble shoot you want after you have attempted all the questions.
4. Get good with YAML

![image-20221128084307597](images/image-20221128084307597.png)

![image-20221128084500617](images/image-20221128084500617.png)

![image-20221128084520293](images/image-20221128084520293.png)

 Solution Lighting Lab 1

![image-20221128090538169](images/image-20221128090538169.png)

![image-20221128091051805](images/image-20221128091051805.png)

![image-20221128092211354](images/image-20221128092211354.png)

![image-20221128092636244](images/image-20221128092636244.png)

![image-20221128092702072](images/image-20221128092702072.png)

![image-20221128092748258](images/image-20221128092748258.png)

![image-20221128092608084](images/image-20221128092608084.png)



Solution Lighting Lab2

![image-20221129091221394](images/image-20221129091221394.png)

![image-20221129091343454](images/image-20221129091343454.png)

![image-20221129092226663](images/image-20221129092226663.png)

![image-20221129092251873](images/image-20221129092251873.png)





#### Mock exam

![image-20221111085052809](images/image-20221111085052809.png)



![image-20221111085115425](images/image-20221111085115425.png)

![image-20221111085314795](images/image-20221111085314795.png)

![image-20221111085627616](images/image-20221111085627616.png)

![image-20221111091909032](images/image-20221111091909032.png)

![image-20221111100633100](images/image-20221111100633100.png)

![image-20221111100957067](images/image-20221111100957067.png)



![image-20221109092537639](images/image-20221109092537639.png)



![image-20221111084858208](images/image-20221111084858208.png)




#### Reference:

[kubectl-commands](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands)