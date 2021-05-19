package models

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
}

func ResponseError(ctx *gin.Context, c ResponseCode) {

	rd := &ResponseData{
		Code:    c,
		Message: c.Msg(),
		Data:    nil,
	}

	response, _ := json.Marshal(rd)
	ctx.Set("response", string(response))
	ctx.JSON(http.StatusOK, rd)
}

func ResponseErrorWithMsg(ctx *gin.Context, code ResponseCode, errMsg string) {
	rd := &ResponseData{
		Code:    code,
		Message: errMsg,
		Data:    nil,
	}

	response, _ := json.Marshal(rd)
	ctx.Set("response", string(response))
	ctx.JSON(http.StatusOK, rd)
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	}

	response, _ := json.Marshal(rd)
	ctx.Set("response", string(response))
	ctx.JSON(http.StatusOK, rd)
}
