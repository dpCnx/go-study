package models

type ResponseCode int64

const (
	CodeSuccess ResponseCode = 1000

	CodeInvalidParams ResponseCode = 1001
	CodeServerBusy    ResponseCode = 1002
)

var msgFlags = map[ResponseCode]string{
	CodeSuccess:       "success",
	CodeInvalidParams: "参数错误",
	CodeServerBusy:    "服务器繁忙",
}

func (c ResponseCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
