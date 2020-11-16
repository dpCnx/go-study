### 问题

```
ERROR: [1] bootstrap checks failed
[1]: the default discovery settings are unsuitable for production use; at least one of [discovery.seed_hosts, discovery.seed_providers, cluster.initial_master_nodes] must be configured
ERROR: Elasticsearch did not exit normally - check the logs at /usr/share/elasticsearch/logs/docker-cluster.log


添加     - "discovery.type=single-node"  解决
```

```
java.nio.file.AccessDeniedException
Error opening log file ‘logs/gc.log’: Permission denied

权限的问题：chmod 777 logs  chmod 777 data (elasticsearch/data elasticsearch/logs)
```

```
elasticsearch-head 连接不上

设置跨域：

http.cors.enabled: true
http.cors.allow-origin: "*"
```

```
head里面的文件通过git下载

git clone https://github.com/mobz/elasticsearch-head.git

再通过文档修改对应的js
```

