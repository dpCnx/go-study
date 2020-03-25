package router

import (
	"github.com/dpCnx/go-study/demo/gin/controller"
	"github.com/dpCnx/go-study/demo/gin/logger"
	"github.com/dpCnx/go-study/demo/gin/middleware"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	//r.Use(gin.Logger())
	//r.Use(gin.Recovery())

	r.Use(middleware.GinLogger(logger.GetLogger()))
	r.Use(middleware.GinRecovery(logger.GetLogger(), true))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/test", controller.Test)
	}

	return r
}
