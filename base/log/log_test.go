package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Error("error")
	Info("info")
	MfError("mferror")
}

func TestSetLevel(t *testing.T) {
	SetLevel(MfErrorLevel)
	Error("error")
	Info("info")
	MfError("mferror2")
}
