package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-study/gin/conf"
	"go-study/gin/router"
	"go.uber.org/zap"

	_ "go-study/gin/log"
)

func main() {

	r := router.LoadRouter()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.C.Server.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)

	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Debug("go-project stop")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error(err.Error())
		return
	}

	zap.L().Debug("go-project quit")
}

/*
	https://blog.csdn.net/raogeeg/article/details/86743953 gin使用swag
	https://www.tizi365.com/archives/288.html  gin使用session
*/

/*
	在项目执行 swag init
	执行 go run main.go
	进入 http://127.0.0.1:8080/swagger/index.html 查看文档
*/

/*
	https://www.cnblogs.com/xiao987334176/p/12340743.html 查看docker仪表盘
*/
