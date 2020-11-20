### k8s

```shell
#安装
https://www.kubernetes.org.cn/7189.html 

#kubeadm安装k8s 组件controller-manager 和scheduler状态 Unhealthy
https://blog.csdn.net/hehj369986957/article/details/107855829

#0/1 nodes are available: 1 node(s) had taints that the pod didn't tolerate.

允许master节点部署pod
kubectl taint nodes --all node-role.kubernetes.io/master-
如果不允许调度
kubectl taint nodes master1 node-role.kubernetes.io/master=:NoSchedule
污点可选参数
      NoSchedule: 一定不能被调度
      PreferNoSchedule: 尽量不要调度
      NoExecute: 不仅不会调度, 还会驱逐Node上已有的Pod
      
Horizontal Pod Autoscaling 仅适用于 Deployment 和 ReplicaSet ，在 V1 版本中仅支持根据 Pod 的 CPU 利用率扩所容，在 v1alpha 版本中，支持根据内存和用户自定义的 metric 扩缩容
```

```shell
kubectl log myapp-pod -c test #查看容器的日志
kubectl get pod -n default #查看容器的日志
```

