package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
	"web/middleware"
)

func main() {

	r := gin.Default()

	r.Use(middleware.PromMiddleware(nil))
	r.GET("/metrics", middleware.PromHandler(promhttp.Handler()))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err.Error())
	}
}

/*
	https://www.cnblogs.com/xiao987334176/p/12340743.html 查看docker仪表盘
*/
