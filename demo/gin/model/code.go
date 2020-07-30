package model

type ResposeCode int64

const (
	CodeSuccess         ResposeCode = 1000
	CodeInvalidParams   ResposeCode = 1001
	CodeServerBusy      ResposeCode = 1005
)

var msgFlags = map[ResposeCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeServerBusy:      "服务繁忙",
}

func (c ResposeCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
