Linux

* 目录

  ```
   / 根目录：所有的数据都是放在此目录下（Linux系统的起点）
   /dev:存放与设备相关的数据
   
   Linux文件格式:FAT32(默认) EXT4(RHEL6) xfs(RHEL7)
   			  swap:交换空间（虚拟内存 缓解真实物理内存不足）
   			  
  ```


* 磁盘

  ```
  /dev/sad    /dev/sdb     /dev/sdc
  
  hd,表示IDE设备
  sd,表示SCSI设备
  vd,表示虚拟化设备
  
  eg:/dev/sda7:表示SCSI设备，第一块磁盘，第7个分区
  
  /proc/ 不占用磁盘的空间
  ```

* 命令行

  ```
    ctrl + alt + fn(1 - 6 ) 窗口
    ctrl + alt + f1 返回桌面
    ctrl + shift + t 新打开一个终端
    # 超级管理员
    $ 普通的用户
    
    cat /etc/redhat-release 查看版本 
    cat -n /etc/redhat-release 显示行号
    cat /proc/meminfo 查看内存
    uname -r 查看内核 
    	==>3.10.0-957.el7.x86_64 //主版本 次版本 修订号 企业版Linux7 64位操作系统
    lscpu 查看cpu
    hostname 查看主机名
    hostname + 名字  修改主机名 eg:hostname dp
    vim /etc/hostname 设置主机名（设置永久主机名）
    ifconfig eth0 192.168.1.1 设置本机的ip地址为192.168.1.1  （临时的）
    poweroff 关机
    reboot 重启
    less 查看文件 输入/a 全文查找a （n,N切换） 
    			   输入q 退出
    head -3 /etc/password 查看前3行
    tail -1 /etc/password 查看后1行
  
    Esc + . 或 Alt + . :粘贴上一个命令的参数
    ctrl + u：清空至行首
    ctrl + w：往回删除一个单词（以空格界定）
    
    cd ~root 去往普通用户的家 
    /home：存放所用普通用户的家目录
    useradd dp 添加一个dp的用户
    
     ls -lh  显示详细信息 并且文件的大小加上了单位
     ls -ld  显示目录本身的信息
     ls -A 显示隐藏文件夹
     ls -R 递归显示 
     
     mkdir .go 创建一个隐藏的go文件夹
     man ls 查看手册 查看文件 输入/a 全文查找a （n,N切换） 
    						  输入q 退出
    
     rm -r -f  -f的优先级比-i高
     cp -i-f   -i的优先级比-f高
     cp -r /home/ /etc/ /mnt/ 把/home/和/etc/拷贝到/mnt/里
     
     date 查看时间
     date -s "2006-5-4 15:00:00" 修改时间
     
      > ：覆盖重定向
      >> : 追加重定向
      -a //and
      -o //or
      
      wc -l /etc/passwd //统计行数
      which find //查找find的目录
      
      systemctl restart chronyd //重起程序 最后的d字母代表守护进程
      systemctl enable chronyd //设置开机自启
      
      管道 | ：将前面命令的输出结果，交由后面命令处理
      eg:head -12 /etc/passwd | tail -5
    
    完整命令的格式
    	==>命令字  【选项】 【参数1】 【参数2】...
    	
  ```

* mount挂载操作

  ```
   光驱设备：ls /dev/cdrom
   挂载光驱到目录：mount /dev/cdrom /dvd
   卸载光驱目录：umount /dvd/
   
   注意事项：
   	1）卸载时，当前不要在访问点内
   	2）挂载时，最好是自己创建的目录，不要用系统的文件 
   	
  ```

* 通配符

  ```
  * 任意多个任意字符 eg： tty*
  ？ 单个字符       eg: tty?
  []只能匹配0-9     eg: tty[3-6]
  {}              eg: tty{1,2,3,5,6}
  ```

* 设置别名

  ```
  alias a='poweroff' 设置别名
  unalias a 删除别名
  \cp -r 临时取消别名
  1)vim /root/.bashrc
  2)alias b='ls' ===>设置永久别名
  ```

* ssh

  ```
  ssh 用户名@ip地址 eg：ssh root@172.25.0.11 
  ssh -X 用户名@ip地址 eg：ssh -X root@172.25.0.11   //启动对方的图形程序
  ssh -i 私钥的路径 ip
  
  ```
  
* rpm

   ```
    wget http://a.txt 下载
    
    rpm  -q vsftpd//查询软件包是否安装
    rpm  -ivh  .rpm软件的路径 //安装软件包 eg:rpm -i vsftpd.rpm
    rpm -Uvh .rpm //安装并升级软件包
    rpm -e vsftpd //卸载
    rpm -ql vsftpd //列出安装清单
    rpm --import RPM-GPG-KEY-CentOS-7  //导入红帽签名文件
   ```
   
* Yum

   ```
   服务端：
   	1.众多的软件包
   	2.仓库清单文件(repodata/)
   	3.构建web服务传递数据
   	
   客户端：
   	客户端配置文件：/etc/yum.repos.d/*.repo ==>错误的配置文件会影响正确的配置文件
   	
    	eg:
    	[rhel7]  //仓库的名字
    	name=rhel7.0 //仓库描述信息
    	baseurl= http://classroom //制定服务器的位置
    	enabled=1	//是否启用该文件
    	gpgcheck=0	//是否检查红帽签名
    	
   yum repolist //列出仓库信息
   yum -y install sssd //安装
   yum remove sssa //卸载  
   yum clean all //清空缓存
   ```

* 配置网络

   ```
   1 配置永久的主机名
   
   	 vim /etc/hostname 设置主机名（设置永久主机名）
   	 
   2 配置永久的ip地址,子网掩码、网关地址
   
      	方式1
   
      	网卡配置文件/etc/sysconfig/network-scripts/ifcfg-ens33
      	命令修改网卡配置文件：nmcli connection
      	1)查看命令识别网卡的名称:nmcli connection show 
   
      	2)进行ip地址配置
      	nmcli connection modify virbr0 ipv4.method manual ipv4.addresses '192.168.122.6' connection.autoconnect yes
      	
      	进行ip地址配置,子网掩码，网关地址
      	eg:nmcli connection modify virbr0 ipv4.method manual ipv4.addresses '172.25.0.110/24 172.25.0.254' connection.autoconnect yes
      	
      	172.25.0.110/24 172.25.0.254 ==>ip地址/子网掩码 网关地址
      	
      	3)激活配置
      	nmcli connection up virbr0	
      	
      	方式2
      	
      	1) nmtui
      	2) edit a connection(回车)
      	3) system eth0(回车)
      	4）require ipv4 addressing for this connection 勾上
      	5) automatically connect 勾上
      	6) nmcli connection up virbr0	
      	
      	查看网关的命令 route
          
   3 配置永久的DNS服务器地址
   
      	vim /etc/resolv.conf
      	nameserver 172.25.254.254	
       
   ```
   
* 用户

   ```
   UID：用户账号标识
   系统程序用户 默认1~999
   所有普通用户UID默认从1000起始
   
   组账号分类：基本组  附加组（从属组）
   所用用户至少属于一个组
   基本组：由Linux系统创建，由Linux将用户加入，与用户同名
   附加组：由root用户创建
   
   useradd dp 添加用户
   useradd -u 1600 dp 添加用户并指定uid
   useradd -d /opt/nsd07 dp 添加用户并指定家目录的地址为/opt/nsd07 
   useradd -s /sbin/nologin dp 添加用户并指定解释器为/sbin/nologin
   
   groupadd dog //添加一个dog组
   useradd -G dog bob 添加用户并指定附加组为dog 
   
   usermod -u 1700 -s /sbin/nologin -d /opt/abc -G tarena dp //修改dp属性
   
   userdel [-r] 用户名 //-r：连同用户家目录一并删除
   
   groupadd [ -g 组ID] 组名 //创建组
   groupdel 组名 //删除组
   
   gpasswd -a dp stugrp //添加dp到stugrp中
   gpasswd -d dp stugrp //删除dp从stugrp组中 
    
   vim /etc/passwd 查看用户列表
   id dp //查看基本信息 
   
   passwd dp //修改dp用户的密码 ==> 局限于root用户 普通用户修改自己的密码直接输入passwd （交互式）
   echo 123 | passwd -- stdin dp //不通过交互式设置密码 局限于root用户
   
   su - dp //切换到dp的用户
   
   /etc/shadow //存放用户密码的文件 
   root:$6$yssMrQSj/lV1jtTJ$l7zQ/5zPBt4DGDX3JsNJtJiS6snBV3It9kNaKTlBVeG6MfWXIQks6l.GF63tzjKVOiTh.VrppWT6ZmtwAB7lq0::0:99999:7:::
   用户名：密码加密字符串(!!如果是两个感叹号 是锁定状态)：上一次修改密码的时间
   
   head -1 /etc/passwd //存放用户信息的文件 
   root:x:0:0:root:/root:/bin/bash
   用户名：密码占位符：UID：基本组GID：用户描述信息：家目录：解释器
   
   /etc/group //存放组基本信息
   
   stugrp:x:2002:dp
   组名：密码占位符：组的id（GID）:组的成员列表
   ```
   
*  tar备份与恢复

   ```
   常见的压缩格式及命令工具
   .gz --> gzip
   .bz2 --> bzip2
   .xz --> xz 
   
   tar工具的常用选项
   	-c 创建归档
   	-x 释放归档
   	-f 指定归档文件名称
   	-z -j -J :调用.gz .bz2 .xz 格式的工具进行处理
   	-t 显示归档中的文件清单
   	-P 保持归档内文件的绝对路径
   	-C 释放的位置
   打包
   	格式：tar 选项 /路径/tar包名字 /路径/源文件 /路径/源文件
   	tar -zcf 备份文件.tar.gz 被备份的文档...
   	tar -jcf 备份文件.tar.bz2 被备份的文档...
   	tar -Jcf 备份文件.tar.xz 被备份的文档...
   解包
   	格式：tar 选项 /路径/tar包名字 /路径/释放的位置
   	
   	tar -xf tar包名字 -C /路径/释放的位置
   	
   显示包的文件清单
      	
      	tar -tf 备份的文档
   ```

* NTP时间同步  

  ```
  chrony 同步时间的软件
  vim /etc/chrony.conf  
     1）server 0.rhel.pool.ntp.org iburst  //修改这个配置
     2）systemctl restart chronyd //重起程序 最后的d字母代表守护进程
     3）systemctl enable chronyd //设置开机自启 
  ```

* cron计划任务

  ```
  周期性任务 crond 
  
  日志文件：/var/log/crond
  
    *:匹配范围内任意时间
    ,:分割多个不连续的时间点 //1,3,
    -:指定连续时间范围 //1-3
    /n:指定时间频率 //*/2
  
     crontab -e -u 用户名 //编辑
     crontab -l -u 用户名 //查看
     crontab -r -u 用户名 //清除所有
  
  eg:crontab -e -u root 会通过vim打开文件
  
  	* * * * * date >> /opt/time.txt
  		 
  /var/spool/cron 保存crond任务的文件
  
  配置任务格式
   分 时 日 月 周 任务命令行(绝对路径)
  
   eg: 30 23 * * * poweroff
   eg: 30 23 1 * 1 poweroff //每月1号和周1满足一个就可以
  
  ```

* 文本文件

   ```
   以-开头:代表文本文件
   以d开头:代表目录
   以l开头:快捷方式
   
   r
   w
   x:执行权限,能够cd切换到此目录 
   
   附属权限
   	Set GID权限 ==> 适用于目录,Set GID 可以使目录下新增的文档自动设置与父目录相同的属组，附加在属组的x位上 (让新增的子文档，自动继承父目录的所属组)
   	chmod g+s 文档
   	--S 代表只有附加权限
   	--s 代表有x,s权限
   	
   	Set UID ==> 附加在属主的x位上(可以让使用者具有一个属主的身份)
   	chmod u+s 文档
   	--S 代表只有附加权限
   	--s 代表有x,s权限
   	
   	Sticky Bit ==> 附加在其他人的x位上(不予许其他人删除，修改操作)阻止用户滥用 w 写入权限
   	chmod o+t 文档
   	--T 代表只有附加权限
   	--t 代表有x,s权限
   
   权限冲突时 ===> 所有者>所属组>其他人 （匹配及停止）
   
   chmod [-R] 文档 //修改权限 -R 递归
   
   chmod u-w /dp
   chmod u+r /dp
   chmod u=rwx,g=rx,o=rx /dp
   chmod ugo=rwx /dp
   chmod o=--- /dp
   
   chown   [-R] 属主 文档
   chown   [-R] :属组 文档
   chown   [-R] 属主:属组 文档
   
   eg:
   chown dp:superman /demo
   chown :superman /demo
   chown dp /demo
   
   
   acl 为个别的用户设置权限 
   
   setfacl [-R] -m u:用户名:权限类别 文档
   setfacl [-R] -m g:组名:权限类别 文档
   setfacl [-R] -b 文档                  //删除全部acl权限
   setfacl [-R] -x u:用户名:权限类别 文档  //删除一个人的acl权限
   
   setfacl -m u:dp:rx /demo 
   getfacl /demo
   ```
   
* LDAP认证

   ```
   轻量级目录访问协议(客户端配置)
   	1 安装sssd软件包,与LDAP沟通的软件 ===> yum -y install sssd
   	2 安装authconfig-gtk图形的工具,配置sssd软件包
   		===> yum -y install authconfig-gtk
   	3 运行authconfig-gtk图形的工具 ===> authconfig-gtk
   	配置里面的dc为域控制器
   	
   	配置完以后
   	systemctl restart sssd //重起程序 最后的d字母代表守护进程
       systemctl enable sssd //设置开机自启
       
    NFS 共享服务
   	
   	1)查看NFS资源 ===> showmount -e classroom.example.com 
   	
   	  /home/guests 172.25.0.0/255.255.0.0 
   	  共享的路径 共享给那些客户端
   	
   	2)访问虚拟机classroom的共享
   	
   		mount classroom.example.com:/home/guests /home/guests
   			   服务端：服务端共享路径                 本地挂载点
   ```
   
* grep

   ```
     grep工具 查找条件 eg:grep root /etc/passwd （查找passwd里面的root） 
     -v 取反匹配 eg:grep -v root /etc/passwd
     -i 忽略大小写
     grep ^root /etc/passwd //以字符串root开头
     grep root$ /etc/passwd //以字符串root结尾
     grep ^$ /etc/passwd //显示空行
     grep ^# /etc/passwd //显示注释行
   ```

* find

   ```
   -type  类型(f 文本文件 d 目录 l 快捷方式)    eg: find /boot -type d
   -name  "文档名称" eg: find /boot -type "passwd"
   -size  +|-文件大小(k(小写) M(大写) G(大写))
   	-size +10w //大于10M
   	-size -10w //小于10M
   	eg:find /boot/ -size +10M
   -user  用户名(按照文档的所有者查看) eg:find / -user student
   -iname 忽略大小写
   -group 更具所属组 eg:find /home -group student
   -mtime 根据文件修改时间 
   	-mtime +10 //过去的10天之前修改和创建的文档
   	-mtime -10 //过去的10天之内修改和创建的文档
	eg:find /var/log -mtime +10
   -maxdepth 限制目录查找的深度 eg:find /home -maxdepth 3 name "*.conf"
   -find .... -exec 处理命令{}\;
   	eg:find /etc/ -name "*tab" -exec cp {} /mnt/ \;
   
   条件可以加多个:find /root/ -name "abb*" -a -type f
   			find /root/ -name "abb*" -type f  //-a 可以省略
   ```
   
* 环境变量

   ```
    影响指定用户的bash 解释环境
    ~/.bashrc
    
     影响所有用户的bash 解释环境
    /etc/bashrc
   ```

   