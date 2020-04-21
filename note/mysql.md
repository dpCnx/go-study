### mysql

* 配置
  ```
  /var/lib/mysql //存放mysql的地方
  /etc/my.cnf    //配置文件
  /var/log/mysqld.log //初始密码的文件 grep 'temporary password' /var/log/mysqld.log
  /mysql/auto.cnf //数据库的uuid位置 每一台的uuid不能相同
  ```
  
* 表结构

   ```
   myisam 存储引擎 
   	.myi 文件存储的是索引
   	.myd 文件存储的是数据
   	.frm 文件存储的是表结构
            支持表级锁
   	
   innodb 
   	.frm 文件存储的是表结构
   	.ibd 文件存储的是数据和索引
   	支持行级锁
  
   	ib_logfile0 记录操作日志信息
   	ib_logfile1 记录操作日志信息
   	
   查询访问多的表 适合使用myisam存储引擎 节省系统资源
   写多的表 适合使用innodb存储引擎 并发访问量大
   ```

* 触发器

   ```
   https://blog.csdn.net/babycan5/article/details/82789099
  https://www.cnblogs.com/cyhbyw/p/8869855.html //乐观锁，悲观锁
   ```
 
   