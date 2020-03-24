```
mongod  -f /usr/local/mongodb/etc/mongodb.conf --shutdown //结束进程

db.user.find({uid:1,name:"root"},{_id:0,name:1})  //查询uid= 1 && name = "root"
```

