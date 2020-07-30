package main

import (
	"context"
	_ "github.com/dpCnx/go-study/demo/gin/logger"
	"github.com/dpCnx/go-study/demo/gin/router"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title 测试
// @version 1.0
// @description  API
// @BasePath /api/v1
func main() {

	r := router.InitRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 开启一个goroutine启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen:", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	// 相当于告诉程序我给你5秒钟的时间你把没完成的请求处理一下，之后我们就要关机啦
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting...")

}

/*
	https://blog.csdn.net/raogeeg/article/details/86743953
	https://www.tizi365.com/archives/288.html  gin使用eg
*/

/*
	在项目执行 swag init
	执行 go run main.go
	进入 http://127.0.0.1:8080/swagger/index.html 查看文档
*/
