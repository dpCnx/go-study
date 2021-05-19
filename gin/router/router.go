package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-study/gin/conf"
	"go-study/gin/controller"
	"go-study/gin/middleware"
)

func LoadRouter() *gin.Engine {

	gin.SetMode(conf.C.Server.Mode)

	r := gin.Default()

	r.Use(middleware.RequestLog())
	r.Use(middleware.Jaeger())

	r.Use(middleware.PromMiddleware(nil))
	r.GET("/metrics", middleware.PromHandler(promhttp.Handler()))

	apiV1 := r.Group("/test/")
	{
		apiV1.GET("/test-get", controller.TestGet)
		apiV1.POST("/test-post", controller.TestPost)
	}

	return r
}
