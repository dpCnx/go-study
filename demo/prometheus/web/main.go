package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	// listen and serve on 0.0.0.0:8080
	if err := r.Run(); err != nil {
		panic(err)
	}

}
