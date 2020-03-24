### Nginx

* 安装nginx

  ```
  01. yum -y install gcc pcre-devel openssl-devel        //安装依赖包 
  02. useradd -s /sbin/nologin nginx 
  03. tar -xf nginx-1.10.3.tar.gz 
  04. cd  nginx-1.10.3 
  ----------------------启动-----------------------------
  05. ./configure
  06. --prefix=/usr/local/nginx    //指定安装路径 
  07. --user=nginx                 //指定用户
  08. --group=nginx                //指定组
  09. --with-http_ssl_module       //开启SSL加密功能 (--with- 添加额外的功能)
  												 (--without- 去掉不需要的功能)
  11. make && make install   		 //编译并安装
  
  eg：./configure --prefix=/usr/local/nginx --user=nginx --group=nginx  --with-http_ssl_module
  ```

* nginx命令的用法

  ```
  01.  /usr/local/nginx/sbin/nginx                    //启动服务 
  02.  /usr/local/nginx/sbin/nginx -s stop            //关闭服务 
  03.  /usr/local/nginx/sbin/nginx -s reload        	//重新加载配置文件 
  04.  /usr/local/nginx/sbin/nginx -V                	//查看软件信息 
  05.  ln -s /usr/local/nginx/sbin/nginx /sbin/       //方便后期使用
  ```

* 升级Nginx服务器 

  ```
  1）编译新版本nginx软件
  01. tar -zxvf nginx-1 .12.2.tar.gz 
  02. cd nginx-1 .12.2 
  03. ./configure
  04. --prefix=/usr/local/nginx 
  05. --user=nginx 
  06. --group=nginx 
  07. --with-http_ssl_module 
  08.  make
  
  2) 备份老的nginx主程序，并使用编译好的新版本nginx替换老版本
  01. mv /usr/local/nginx/sbin/nginx /usr/local/nginx/sbin/nginxold 
  02.	cp objs/nginx  /usr/local/nginx/sbin  	//拷贝新版本 
  04. make upgrade                            //升级(可以使用killall nginx 然后重启nginx) 
  	/usr/local/nginx/sbin/nginx -t 06. nginx: the configuration file 					/usr/local/nginx/conf/nginx.conf syntax is ok 07. nginx: configuration file 		/usr/local/nginx/conf/nginx.conf test is successful 08. kill -USR2 `cat 			/usr/local/nginx/logs/nginx.pid` 09. sleep 1 10. test -f 							/usr/local/nginx/logs/nginx.pid.oldbin 11 . kill -QUIT `cat 						/usr/local/nginx/logs/nginx.pid.oldbin` 
  05.  /usr/local/nginx/sbin/nginx – v                //查看版
  ```

* 其他

  ```
  注意:第一次创建要加-c 继续添加不用加-c 
  htpasswd -c /usr/local/nginx/pass   tom        //创建密码文件  
  
  设置nginx开机自启:
    vim /etc/rc.local
  设置nginx可以压缩的文件类型：
    可以参考/usr/local/nginx/conf/mime.types
  ```