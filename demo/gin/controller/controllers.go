package controller

import (
	_ "github.com/dpCnx/go-study/demo/gin/docs"
	"github.com/dpCnx/go-study/demo/gin/model"
	"github.com/gin-gonic/gin"
)

// @Summary Test
// @Description gin添加Api测试
// @Produce  json
// @Param param formData string true "参数"
// @Success 0 {object} model.ResponseData
// @router /test [post]
func Test(c *gin.Context) {
	p := c.PostForm("param")
	model.ResponseSuccess(c, p)
}
