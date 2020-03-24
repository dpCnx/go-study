### mysql

* 安装

  ```
  http://dev.mysql.com/downloads/mysql/
  
  mysql-5.7.28-1.el7.x86_64.rpm-bundle.tar
  mysql-community-client 	//客户端应用程序
  mysql-community-common 	//数据库和客户端共享数据
  mysql-community-devel 	//客户端应用程序得库和头文件
  mysql-community-embedded	//嵌入式函数库
  mysql-community-embedded-compat 	//嵌入式兼容函数库
  mysql-community-embedded-devel		//头文件和库文件作为mysql的嵌入式库文件
  mysql-community-libs	//mysql数据库客户端应用程序的共享库
  mysql-community-libs-compat //客户端应用程序的共享兼容库
  
  systemctl stop mariadb
  rm -rf /etc/my.conf
  rm -rf /var/lib/mysql/*
  rpm -e --nodeps mariadb-server mariadb
  
  yum -y install perl-Data-Dumper perl-JSON perl-Time-HiRes //安装依赖包
  
  rpm -Uvh mysql-community-*.prm //安装包
  
  systemctl start mysqld //启动mysql服务
  systemctl enable mysqld //设置开机自启
  systemctl status mysqld //查看mysql服务状态
  ```

* 配置
  ```
  /var/lib/mysql //存放mysql的地方
  /etc/my.cnf    //配置文件
  /var/log/mysqld.log //初始密码的文件 grep 'temporary password' /var/log/mysqld.log
  /mysql/auto.cnf //数据库的uuid位置 每一台的uuid不能相同
  ```
  
* 修改密码策略

  ```
  show variables like "%password%";     //查看变量
  
  set global validate_password_policy=0;      //只验证长度 
  set global validate_password_length=6;     //修改密码长
  alter user root@"localhost" identified by "123456";  //修改登陆密
  
  永久修改
  systemctl stop mysqld
  
  vim /etc/my.cnf
  [mysqld]
  global validate_password_policy=0
  global validate_password_length=6  
  
  systemctl start mysqld
  ```

* 统计一些常识

  ```
  select database()					//显示当前的库
  show slave status  //查看是不是从库
  show master status //查看正在使用的日志名
 
  手动删除授权库：
    delete from mysql.user where user not in ("root","mysql.sys")
    flush privileges //需要手动刷新才能生效
  
  设置外键的表,只要外键的值符合主键,可以重复添加,可以设置为null
  eg:
    id name            user_id   money
    1  bob				1        125000
    2  jack				1        3000
    					    null     3000
  ```

* 表结构

   ```
   myisam 存储引擎 
   	.myi 文件存储的是索引
   	.myd 文件存储的是数据
   	.frm 文件存储的是表结构
   	
   innodb 
   	.frm 文件存储的是表结构
   	.ibd 文件存储的是数据和索引
   	
   	ib_logfile0 记录操作日志信息
   	ib_logfile1 记录操作日志信息
   	
   查询访问多的表 适合使用myisam存储引擎 节省系统资源
   写多的表 适合使用innodb存储引擎 并发访问量大
   ```

* 数据备份

   ```
   物理备份
       备份数据的时候要记得修改mysql的用户和用户组为mysql
       如果只是备份一个库，引擎为innodb时，会备份失败 
   
       cp -r /var/lib/mysql /root/mysqlall.bak
       scp -r /root/mysqlall.bak 192.168.4.51:/root
       systemctl stop mysqld
       rm -rf /var/lib/mysql
       cp -r /root/mysqlall.bak /var/lib/mysql
       chown -R mysql:mysql /var/lib/mysql
       systemctl start mysql
       
   逻辑备份
   	完全备份
   	缺点:恢复的时候会加写锁
   	
   	
   	增量备份 
   		启动binlog 不会记录select的操作
   		
   		vim /etc/my.conf
   		
   		[mysqld]
   		
   		log_bin //默认日志名和位置
   		#log_bin=mylog	//默认位置，自己指定的日志名
   		#log_bin=/logdir/mylog  //指定位置和名字
   		server_id=50   //1-255
   		binlog_format="mixed"
   		
   		如果是自己指定的位置 mkdir /logdir
   		chown mysql /logdir
   		
   		systemctl restart mysqld
   		 
   		 //在备份的时候新建binlog文件
   		mysqldump -u root -p123456 --flush-logs db4> /root/db4.sql
   		//查看正在使用的日志文件
   		show master status
   	
            mysqldump备份的时候，恢复的宁一种方法:
                use db1;
                source /root/db1.sql;  
    
   	XtraBackup备份:
   		innobackupex 支持myisam 但是每次备份都是完全备份，不是增量备份
   					 innodb的时候是增量备份
   ```

* 主从配置

   ```
   change master to master_host="192.168.4.45"  //修改主数据库的链接地址
   
   stop slave //设置临时不同步
   
   删除从库：
   rm -rf master.info
   rm -rf mysql52-relay-bin.*
   rm -rf relay-log.info
   systemctl restart mysqld
    
   设置只同步哪些库
    
    vim /etc/my.conf
     		
    [mysqld]
     		
     log_bin //默认日志名和位置
     #log_bin=mylog	//默认位置，自己指定的日志名
     #log_bin=/logdir/mylog  //指定位置和名字
     server_id=50   //1-255
     binlog_format="mixed"
     binlog_do_db=db1,db2 //只同步哪些库到binlog日志里
     #binlog_ignore_db=db1,db2 //不同步哪些库到binlog日志里
    
    主从从配置的时候必须加log_slave_updates 
    如果不加：50 --> 51 --> 52
    51通过io线程从50上面读取到的binlog信息不能同步到51上的binlog文件里面
    那么52通过io线程去51上面读取数据就读不到
   ```
   
* 配置读写分离

   ```
   配置maxscale
   
   01 . [root@maxscale mysql]# vim /etc/maxscale.cnf.template 
   02. [maxscale] 
   03. threads=auto         //运行的线程的数量，根据计算机的内核数决定的 
   
   [MySQL Monitor]        //定义监控的数据库服务器 
   25. user=scalemon      //监视数据库服务器时连接的用户名scalemon 
   26. passwd=123456      //密码123456
   	==> 主要监视51，52是不是在运行
   		主从结构是否正常
   		谁是主库
   		谁是slave库
   		
   [Read-Write Service]         //定义读写分离服务 
   45. user=maxscaled           //用户名 验证连接代理服务时访问数据库服务器的用户是否存在 
   46. passwd=123456             //密码
   	==> 这个账号密码是用来登陆到server里面去检查客户端的账号密码是否存在的
   ```

* 配置MySQL多实例

   ```
   创建进程运行的所有者和组 mysql 初始化授权库 可以不做 直接启动多实例 会自动去授权多实例
   ```

* MHA集群

   ```
   配置54从库:
   	 vim /etc/my.cnf
   	 
   	 server_id=54
   	 relay_log_purge=0 //不删除中继日志
        plugin-load="rpl_semi_sync_slave=semisync_slave.so"
        rpl-semi-sync-slave-enabled = 1 
        
   安装：
   	yum -y install perl-*   //安装perl依赖包
   	yum -y install  perl-*.rpm //这步操作是有些依赖包在光盘上没有,需要联网下载，如果是直接连接的外网，应该就不需要了
   	
    master_ip_failover_script=/usr/local/bin/master_ip_failover  //故障重启脚本
   
   配置完成后,需要自己手动在51上绑定vip：
   	ifconfig eth0:1 192.168.4.100
   	ifconfig eth0  //查看ip是否绑上
   	
   集群启动成功以后： 访问 mysql -uroot -p123456 -h192.168.4.100
   ```
   
* 视图

   ```
   create view v1 as select name ,uid from user where shell = "/bin/bash"
   
   insert into v1 vaules("aa",11,"nobin/bash")
   
   由于视图在创建的时候添加了条件,所以试图里面查不到这条记录,但是user可以 
   
    //有相同的名字的时候需要自己定义别名
   create view v3(a,b,c,d) as select * from t1,t2 where t1.name = t2.name
   
   algorithm //试图的算法 
   		undefined //默认格式是merge ==>  
   	   	merge //替换方式 ==> 执行sql语句,直接执行当前的操作
          	temptable  //具体化方式 ==> 执行sql语句时，先执行as后面的，在执行当前的操作
   ```
   
* 存储过程

   ```
   select body from mysql.proc where name = "p1"/G   ==> 查看存储过程的用途
   ```

* 触发器

   ```
   https://blog.csdn.net/babycan5/article/details/82789099
   ```

   