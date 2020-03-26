package main

import (
	"github.com/dpCnx/go-study/demo/gin/logger"
	"github.com/dpCnx/go-study/demo/gin/router"
)

// @title 测试
// @version 1.0
// @description  API
// @BasePath /api/v1
func main() {

	logger.InitLogger()

	router := router.InitRouter()
	_ = router.Run()
}

/*
	https://blog.csdn.net/raogeeg/article/details/86743953
	https://www.tizi365.com/archives/288.html  gin使用eg
*/
