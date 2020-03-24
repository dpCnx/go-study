* 重定向输出

  ```
  > :只收集前面命令的正确输出信息,写入文本文件中
  2> :只收集前面命令的错误输出信息,写入文本文件中
  &> :收集前面命令的正确与错误输出信息,写入文本文件中
  
  黑洞设备 /dev/null
  ```

* SELinux

  ```
  enforcing 强制
  permissive 宽松
  disabled 彻底禁用
  
  任何一种运行模式,变成disabled都要经历重起系统
  
  切换运行模式
  	getenforce 查看当前模式
  	setenforce 1|0
  	永久修改:/etc/selinux/config 文件
  	
  SELinux布尔值
  	服务功能的开关 on 或 off
  	- 需要加-P 选项才能实现永久设置
  	getsebool -a | grep samba
  	setsebool samba_export_all_ro on //设置读服务功能
  	
  	
  网页文件默认存放路径:/var/www/html 
  默认网页文件的名字:index.html
  ```

* 防火墙

  ```
  防火墙:隔离作用
  
  硬件防火墙
  软件防火墙
  firewalld
  
  根据所在的网络场所区分，预设保护规则集
  
  - public:仅允许访问本机的ssh，ping,dhcp等少数几个服务
  - trusted:允许任何访问
  - block:阻塞任何来访请求
  - drop:丢弃任何来访的数据包 
  
  firewall-cmd --get-default-zone        //查看区域
  firewall-cmd --set-default-zone=block  //设置区域
  
  firewall-cmd --zone=public --list-all //列出public区域中规则
  firewall-cmd --zone=public --add-service=http //为public区域添加允许的协议http
  
  firewall-cmd --permanent --zone=public --add-service=http //永久配置(permanent)
  firewall-cmd --reload //重新加载
  
  firewall-cmd --zone=block --add-source=127.0.0.1
  firewall-cmd --zone=block --remove-source=127.0.0.1
  
  通过防火墙转发端口：
  firewall-cmd --permanent --zone=public --add-forward-port=port=5423:proto=tcp:toport=80
  firewall-cmd --reload //重新加载
  firewall-cmd --zone=public --list-all //列出public区域中规则
  
  互联网常见协议:
  	http:超文本传输协议 80
  	https:安全的超文本传输协议 443
  	FTP:文件传输协议 21
  	DNS:域名解析协议 53
  	SMTP:用户发邮件协议 25
  	pop3:用户收邮件协议 110
  	telnet:运程管理协议 23
  	TFTP:简单的文本传输协议 69
  	SNMP:网络管理协议 161
  ```
  
* 配置聚合链接

  ```
  1) 创建虚拟网卡team0
  # nmcli connection add type team con-name team0 ifname team0 autoconnect yes config '{"runner":{"name":"activebackup"}}'
  
  nmcli connection 添加 类型为 team(组队) 配置文件名 team0 网卡名 team0 每次开机自动启动 工作模式配置为 热备份(activebackup)
  
  //轮询式(roundrobin) -->流量负载均衡的作用
  
  ifconfig
  nmcli connection delete team0 //删除
  
  2) 添加成员
  nmcli connection add type team-slave con-name team0-1 ifname eth1 master team0
  nmcli connection add type team-slave con-name team0-2 ifname eth2 master team0
  
  nmcli connection 添加 类型 team-成员 配置文件名为 team0-2 网卡名 eth2
  
  nmcli connection delete team0-1 //删除
  
  3)为虚拟网卡team0配置IP地址
  nmcli connection modify team0 ipv4.method manual ipv4.addressess 192.168.1.1/24 connection.autoconnect yes
  
  4)激活配置
  nmcli connection up team0
  nmcli connection up team0-1
  nmcli connection up team0-2
  
  5)查看链路聚合的命令
  ifconfig eth1 down(up) //关闭(打开)eth1
  teamdctl team0 state //查看状态
  ```

* Samba

  ```
  为客户机提供共享使用的文件夹
  
  添加用户:pdbedit -a 用户名
  查询用户:pdbedit -L
  删除用户:pdbedit -x 用户名
  
  创建samba共享账号
  	useradd -s /sbin/nologin harry
  	pdbedit -a harry
  创建共享目录 
  	mkdir /common
  配置文件
  	
  	[自定共享名]
  	path = 文件夹绝对路径
  	;public = no|yes //默认no
  	;browseable = yes|no //默认yes
  	;read only = yes|no	//默认yes
  	;write list = 用户1.... //默认无
  	;valid users = 用户1.... //默认任何用户
  	;hosts allow = 客户机地址....
  	;hosts deny = 客户机地址....
  	
  	vim /etc/samba/smb.conf
  
  	[common]
  	path = /common
  	
  重启smb服务
  	systemctl restart smb
  	systemctl enable smb
  	
	
  samba-client(客户端软件)
  	1)yum -y install samba-client
  
  	2)smbclient -L 127.0.0.1 #查看对方有哪些共享
  
  	3)smbclient -U harry //127.0.0.1/common
  
  利用挂载的方式访问:
  	1)安装软件cifs-utils(让本机支持cifs文件系统)
  	yum -y install cifs-utils
  	
  	2)挂载访问
  	mount -o user=harry,pass=123 //127.0.0.1/common /mnt/nsd
  	df -h
  	
  	3)开机自动挂载
  	_netdev 声明网络设备,配置完所有的网络参数后,再进行挂载该设备
  	
  	vim /etc/fstab
  	
  	//127.0.0.1/common /mnt/nsd cifs defaults,user=harry,pass=123,_netdev 0 0
  	
  	umount /mnt/nsd
  	df -h
  	mount -a
  	df -h
  	
  	
  配置可读可写的samba
  
  	1) 建立新的共享目录
  		mkdir /devops
  		echo abc > /devops/abc.txt
  		
  	2) 修改主配置文件
  		vim /etc/samba/smb.conf
  		[devops]
  		path = /devops
  		write list = d
  	
  	3) 重启smb服务
  		systemctl restart smb
  		
  	4) SELinux布尔值
  		getsebool -a | grep samba
  		setsebool samba_export_all_rw on //设置读写服务功能
  		
  	5)修改本地目录的权限
  		 setfacl -m u:d:rwx /devops/
  		 getfacl /devops/
  ```
  
* NFS共享

  ```
  所需nfs-utils
  
  1)检测软件包是否安装
  	rpm -q nfs-utils
  2)修改vim /etc/exports  -文件夹路径   共享的ip地址(权限)
  	
  	/abc  *(ro)
  	/test 127.0.0.1/24(ro)
  	
  3)重启服务nfs-server
  	systemctl restart nfs-server
  	systemctl enable nfs-server
  	
  	
  服务端
  
  	vim /etc/fstab
  	127.0.0.1/abc /mnt/nsd nfs defaults,_netdev 0 0
  	mount -a
  	df -h
  	
  ```

  

  

