### redis

```
****************************
哨兵服务器
bind 0.0.0.0 //本机的所有地址，不写就是本机的127.0.0.1的地址
port 26379 //不写就默认是26379，可以自己修改
sentinel   monitor   redis51    192.168.4.51   6379   
sentinel   auth-pass  redis51 123456  //连接主库密码，若主库有密码加上这一行

****************************
AOF与RDB一起开的时候,默认先加载AOF 。 最好一起开启

appendfsync always 有新写操作立即记录 -->存盘
appendfsync everysec 每秒记录一次 -->存盘
appendfsync no 从不记录 -->只记录到aof文件,不存盘

```

