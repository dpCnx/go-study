### es 

```
配置/etc/hosts 192.168.1.1 es1

启动的时候报错:the default discovery settings are unsuitable for production use; at least one of [discovery.seed_hosts, discovery.seed_providers, cluster.initial_master_nodes] must be configured

修改配置:/etc/elasticsearch/elasticsearch.yml

cluster.name: my-application
node.name: es1
network.host: 0.0.0.0
discovery.seed_hosts: ["es1"]
cluster.initial_master_nodes: ["es1"]
```

