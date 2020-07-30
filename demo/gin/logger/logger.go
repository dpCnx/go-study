package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var MyLogger zap.Logger

func init() {

	errWriteSyncer := getErrLogWriter()
	encoder := getEncoder()

	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)
	var core zapcore.Core

	// NewTee 可以指定多个日志配置
	core = zapcore.NewTee(
		// 创建一个将debug级别以上的日志输出到终端的配置信息
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		// 将error级别以上的日志输出到err文件
		zapcore.NewCore(encoder, errWriteSyncer, zapcore.ErrorLevel),
	)

	logger := zap.New(core, zap.AddCaller()) // 根据上面的配置创建logger
	zap.ReplaceGlobals(logger)               // 替换zap库里全局的logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig) // json格式日志
}

func getErrLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/err.log",
		MaxSize:    100,  // 日志文件大小 单位：MB
		MaxBackups: 1,    // 备份数量
		MaxAge:     7,    // 备份时间 单位：天
		Compress:   true, // 是否压缩
	}

	return zapcore.AddSync(lumberJackLogger)
}
