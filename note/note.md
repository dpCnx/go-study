### Note

```
https://studygolang.com/articles/18576?fr=sidebar   //设置代理
```

### Gin

```
https://blog.csdn.net/cyberspecter/article/details/100602552

c.Abort() 在被调用的函数中阻止挂起函数。注意这将不会停止当前的函数
router.Any() 匹配任何方法路由
router.NoRoute() 可以匹配404
mustbind:bind() bindjson() 如果直接使用这个系列，发生错误的时候gin会直接写入400到响应头里面,不予许自己修改
shouldbind:bind() shoouldbindjson() 如果直接使用这个系列，发生错误的时候自己处理
gin.session 
```

### sql

```
sql.nullstring

beego 获取返回值
beego.InsertFilter("*", beego.AfterExec, func(c *context.Context) {}, false) //第二个参数必须位false
requeststr := string(c.Input.RequestBody)
outputBytes, _ := json.Marshal(c.Input.Data()["json"])
```

