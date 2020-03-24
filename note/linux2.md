### 磁盘

* 分区

  ```
  MBR/msdos 分区模式
  	分区规则:
  		分区类型:主分区,扩展分区,逻辑分区(1-4个主分区,或者3个主分区 + 1个扩展分区(n个逻辑分区)) 
  		最大支持容量为2.2TB的磁盘
  		lsblk 查看分区
  		fdisk /dev/vdb 
  		help 查看 
  		n 创建新的分区 结束+10G
  		p 查看分区表  
  		d 删除分区
  		w 保存并退出
  		
  		eg:p--> n 结束+10G --> w
  		
  	格式化,赋予空间 文件系统(规则)
			
  		mkfs.  //查看可以格式化的类型
  
  		mkfs.ext4 /dev/vdb1
  		blkid /dev/vdb1
  		
  		mkfs.xfs /dev/vdb2
  		blkid /dev/vdb2
  		
  	挂载使用分区
  		
  		mount /dev/vdb1 /mypart1
  		mount /dev/vdb2 /mypart2
  		df -h //显示已经挂载的设备使用情况
  		
  	开机自动挂载
  		
  		vim /etc/fstab
  		设备路径 挂载点 类型 参数 备份标记 检测顺序
  		/dev/vdb1 /mypart1 ext4 defaults 0 0
  	
  		 备份标记  defaults 默认(rw GID UID ...),可以只设置ro wo ...
  				 defaults,ro 代表默认的权限都要,但是只是只读的
  		 备份标记 0 不备份 1 备份
  		 检测顺序 0 不检测 1 检测
  		 
  		 mount -a 检测挂载是否成功
  		 df -h //显示已经挂载的设备使用情况
  		 
  		 partprobe //刷新新的分区表
  		 lsblk 查看分区
  		 
  GPT分区模式,最大到18EB 
  	1EP=1000PB
  	1PB=1000TB 
  ```
  
  

* LVM 逻辑卷的管理

  ```
  作用：1.可以整合分散的空间
  	 2.容量大小可以扩大 
  
  零散空间存储-->整合的虚拟磁盘-->虚拟的分区
  物理卷：pv  pvs 查看详细信息
  卷组：VG  vgs 查看详细信息
  逻辑卷：LV lvs 查看详细信息
  
  1)直接创建卷组
  	vgcreate systemvg /dev/vdc1 /dev/vdc2
  2)通过卷组划分逻辑卷
  	lvcreate -n mylv -L 16G systemvg  
  	lvcreate -l 50 -n mylv2 systemvg //小写的l代表PE的个数(50个PE)
  	lvcreate -s 50 -n mylv2 systemvg //小写的s代表PE的大小(50PE)
  	lvs
  	vgs
  	mkfs.ext4 /dev/systemvg/mylv
  	blkid /dev/systemvg/mylv
  	
  	vim /etc/fstab
  	设备路径 挂载点 类型 参数 备份标记 检测顺序
  	/dev/systemvg/mylv /mypart1 ext4 defaults 0 0
  	 df -h //显示已经挂载的设备使用情况
  	
  	ls /dev/dm-0 逻辑卷的位置
  	ls -l /dev/systemvg/mylv 我们的逻辑卷的位置 (快捷方式)   
  	 ==> /dev/systemvg/mylv -> ../dm-0 也是指向系统直接生成的逻辑卷的位置
  	 
  逻辑卷的扩展
  
  1)扩展空间的大小
  	lvextand -L 18G /dev/systemvg/mylv //扩展到18G
  	lvextand -L +2G /dev/systemvg/mylv //添加2G
  	lvs 查看详细信息
  	
  2)扩展文件系统的大小
  	resize2fs:扩展ext4文件系统
  	xfs_growfs:扩展xfs文件系统
  	
  	resize2fs /dev/systemvg/mylv
  	df -h //显示已经挂载的设备使用情况
  	
  当卷组空间不足的时候:
  	vgextend systemvg /dev/vdc3 //先扩展卷组
  	vgs 
  	lvextand -L 25G /dev/systemvg/mylv //扩展到18G
  	df -h
  	resize2fs /dev/systemvg/mylv
  	df -h
  	
  vgdisplay systemvg //查看vg详细信息
  	
  ext4文件系统支持缩减
  xfs文件系统不支持缩减
  
  卷组划分空间的单位：PE 1PE = 4M
  vgchange -s 1M systemvg //修改systemvg的1PE=1M
  
  逻辑卷的删除：
  	1)先删除逻辑卷本身(先卸载)
  	2)再删除卷组
  	3)最后删除物理卷(可选)
  	
  	lvremove /dev/systemvg/mylv 
  	//如果报错,先删除挂载,在删除逻辑卷
  	umount/lvm 
  	lvremove /dev/systemvg/mylv
  	
  	//删除卷组 要删除所有基于卷组的逻辑卷
  	vgremove systemvg
  	
  	
  	pvremove /dev/vdc[1-3]
  	
  	
  ```
  
  