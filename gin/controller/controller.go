package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go-study/gin/models"
)

func TestGet(c *gin.Context) {

	ids := c.Query("id")
	id, err := strconv.Atoi(ids)

	if err != nil {
		models.ResponseErrorWithMsg(c, models.CodeInvalidParams, err.Error())
		return
	}

	ctx := c.Request.Context()

	var u models.User

	if err = u.Query(ctx, id); err != nil {
		models.ResponseErrorWithMsg(c, models.CodeServerBusy, err.Error())
		return
	}

	models.ResponseSuccess(c, u)
}

func TestPost(c *gin.Context) {

	var u models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		models.ResponseErrorWithMsg(c, models.CodeInvalidParams, err.Error())
		return
	}

	ctx := c.Request.Context()

	if err := u.Insert(ctx); err != nil {
		models.ResponseErrorWithMsg(c, models.CodeServerBusy, err.Error())
		return
	}

	models.ResponseSuccess(c, nil)
}
