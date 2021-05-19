package log

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorlog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.Lshortfile|log.LstdFlags)
	infolog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.Lshortfile|log.LstdFlags)

	f, _  = os.OpenFile("f.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	mlog  = io.MultiWriter(f, os.Stdout)
	mflog = log.New(mlog, "error to file", log.Ldate|log.Ltime|log.Llongfile)

	loggers = []*log.Logger{errorlog, infolog, mflog}
	mu      sync.Mutex
)

var (
	Error   = errorlog.Println
	Errorf  = errorlog.Printf
	Info    = infolog.Println
	Infof   = infolog.Printf
	MfError = mflog.Println
)

const (
	InfoLevel = iota
	ErrorLevel
	MfErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if InfoLevel < level {
		infolog.SetOutput(ioutil.Discard)
	}

	if ErrorLevel < level {
		errorlog.SetOutput(ioutil.Discard)
	}

	if MfErrorLevel < level {
		mflog.SetOutput(ioutil.Discard)
	}
}
