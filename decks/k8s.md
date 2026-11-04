# Kubernetes

## Basics

### Workload Resourses
`Deployment`
`ReplicaSet`
`StatefulSet`
`DaemonSet`
`Job`
`CronJob`

### Control Plane Components
`kube-apiserver`
`etcd`
`kube-scheduler`
`kube-controller-manager`
`cloud-controller-manager`

### Node Components
`kubelet`
`kube-proxy`
`container runtime`

### Kubernetes provides you with
Sd `Service discovery`
Lb `Load balancing`
So `Storage orchestration`
Ar `Automated rollouts and rollbacks`
Ab `Automatic bin packing`
Sh `Self-healing`
Se `Secret and configuration management`

### Expose Services
?
`ClusterIP`
`NodePort`
`LoadBalancer`
`ExternalName`

### Container restart policy
`Always`
`OnFailure`
`Never`

### Pod status
PS `PodScheduled`
PN `PodHasNetwork`
CR `ContainersReady`
I `Initialized`
R `Ready`

### Pod lifecycle
`Pending`
`Running`
`Succeeded` / `Failed`
`Unknown`

### Probe
`startup`
`liveness`
`readiness`

### Probe checking mechanisms
`exec`
`grpc` 
`httpGet` 
`tcpSocket`

### Quality of Service Class
`Guaranteed`
`Burstable`
`BestEffort`

### RBAC
`Role-Based Access Control`

## kubectl

### kubectl basic commands
kubectl `get`
kubectl `describe`
kubectl `logs`
kubectl `exec`

### kubectl logs
kubectl `-n my-namespace` logs `my-pod`
kubectl logs `my-pod -c my-container -f`

### kubectl exec
kubectl exec `-it pod-name -- bash`

### kubectl set default namespace
kubectl `config set-context --current --namespace=my-namespace`

### kubectl get deployment info
kubectl `describe deployment name`
kubectl `get deployment name -o yaml`

## Docker

### Docker objects
`images`
`containers`
`networks`
`volumes`
`plugins`
and other objects

### Docker run
docker `run -i -t ubuntu:20.04 /bin/bash`
`-i, --interactive` Keep STDIN open even if not attached
`-t, --tty`  Allocate a pseudo-TTY

### Get Image Info
docker `image ls`
docker `image inspect image-tag`

### CNI
`Container Network Interface`

