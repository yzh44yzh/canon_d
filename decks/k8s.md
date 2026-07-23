# Kubernetes

## Basics

### Workload Resourses
Deployments
ReplicaSet
StatefulSets
DaemonSet
Jobs
CronJob

### Control Plane Components
kube-apiserver
etcd
kube-scheduler
kube-controller-manager
cloud-controller-manager

### Node Components
kubelet
kube-proxy
container runtime

### Container restart policy
Always
OnFailure
Never

### Probe
startup
liveness
readyness

### Quality of Service Class
Guaranteed
Burstable
BestEffort

### RBAC
Role-Based Access Control

## kubectl

### kubectl basic commands
kubectl get
kubectl describe
kubectl logs
kubectl exec

### kubectl logs
kubectl logs my-pod -c my-container -f

### kubectl exec
kubectl exec -it pod-name -- bash

### kubectl set default namespace
kubectl config set-context --current --namespace=my-namespace

### kubectl get deployment info
kubectl describe deployment name
kubectl get deployment name -o yaml
