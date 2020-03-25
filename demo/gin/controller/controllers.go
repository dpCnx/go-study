package controller

import (
	_ "github.com/dpCnx/go-study/demo/gin/docs"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Test
// @Description gin添加Api测试
// @Produce  json
// @Param param formData string true "参数"
// @Success 0 {object} Respose
// @router /test [post]
func Test(c *gin.Context) {
	p := c.PostForm("param")
	c.JSON(http.StatusOK, gin.H{"p": p})
}
