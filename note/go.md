### go

```go
//go 夜聊
https://talkgo.fm/
//查询三方库
https://pkg.go.dev/
//设置代理
https://studygolang.com/articles/18576?fr=sidebar 
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