### git

```shell
#常识：git编辑框编辑不到的时候 可以点Q字母

#对git自我介绍
git config --global user.name "zhangsan"
git config --global user.email "zhangsan@qq.com"

#查看用户信息
git config --list 
#初始化
git init 
#创建忽略文件
touch .gitignore 
#查看文件的状态
git status 

#回退上一个版本
git reset --hard HEAD^ 
#回退n个版本
git reset --hard HEAD~n  
#回退到某一个版本 id 是提交的是id ==》commit
git reset --hard id 
#查看所有的记录
git reflog  
#查看当前分支
git branch 
#创建一个为xx的分支
git branch + 名字
#切换分支
git checkout + 名字
#创建一个叫xx的分支，并且切换到这个分支上面
git checkout -b + 名字
#删除xx分支
git branch -d + 名字
#合并xx分支
git merge + 名字 

#--------------------------远程仓库----------------------------------------

#在用户主目录下，看看有没有.ssh目录，如果有，再看看这个目录下有没有id_rsa和id_rsa.pub这两个文件。如果没有,执行下面命令，创建SSH Key： （id_rsa是私钥，id_rsa.pub是公钥）
#建议使用真实的邮箱
ssh-keygen -t rsa -C "youremail@example.com"
#添加一个远程仓库
git remote add + "远程仓库的名" 
#查看远程仓库
git remote 
#往远程仓库推送代码
git push -u + "远程仓库的名字" + "需要推送的分支"    
#第一次连接远程仓库必须加-u 链接成功以后 下次在推送可以不加 
git push + "远程仓库的名字" + "需要推送的分支"
#克隆
git clone
#查看远程仓库分支
git branch -a
#拉取远程仓库的代码
git pull 
#add之前的撤销
git checkout -- code.txt
#add之后的撤销
git reset HEAD code.txt  
#当前code的代码与代码库的不同
git diff HEAD -- code.txt  
#对比这两个版本的不同
git diff HEAD HEAD^ -- code.txt 
#保存当前工作状态
git stash  
#恢复当前工作状态
git stash pop 
#如果后面什么都不跟的话 就是上一次add 里面的全部撤销了
git reset HEAD 
#就是对某个文件进行撤销了
git reset HEAD XXX/XXX/XXX.java 
#查看远程分支
git branch -r 
#拉取分支代码
git checkout -t origin/daily/1.4.1 
 
#---------------------------问题-----------------------------------------

#git commit --amend 可以取消上一次提交,修改完之后按esc键退出编辑状态，再按大写ZZ就可以保存退出vim编辑器

```

### go

```go
//设置代理
https://studygolang.com/articles/18576?fr=sidebar 
https://blog.csdn.net/cyberspecter/article/details/100602552
//--------------------------gin---------------------
//匹配任何方法路由
router.Any() 
//可以匹配404
router.NoRoute() 
//如果直接使用这个系列，发生错误的时候gin会直接写入400到响应头里面,不予许自己修改
mustbind:bind() bindjson() 
//如果直接使用这个系列，发生错误的时候自己处理
shouldbind:bind() shoouldbindjson() 
//--------------------------map----------------------
//map 遍历无序 可能是因为在不同的平台生成的key不一样，有32位和64位，不好维护，为了稳定和安全就做成无序的了。从那个桶开始，桶里面的链表位置都随机了

// 生成随机数 r
r := uintptr(fastrand())
if h.B > 31-bucketCntBits {
	r += uintptr(fastrand()) << 31
}

// 从哪个 bucket 开始遍历
it.startBucket = r & (uintptr(1)<<h.B - 1)
// 从 bucket 的哪个 cell 开始遍历
it.offset = uint8(r >> h.B & (bucketCnt - 1))
//------------------------beego--------------------
//beego 获取返回值
beego.InsertFilter("*", beego.AfterExec, func(c *context.Context) {}, false) //第二个参数必须位false
requeststr := string(c.Input.RequestBody)
outputBytes, _ := json.Marshal(c.Input.Data()["json"])
```

