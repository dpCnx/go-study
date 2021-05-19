package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go-study/gin/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

)

func init() {

	infoWriteSyncer := initInfoLogWriter()
	errWriteSyncer := initErrLogWriter()
	encoder := getEncoder()

	/*	level := zap.NewAtomicLevel()
		level.SetLevel(zap.DebugLevel)*/

	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.InfoLevel
	})

	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.InfoLevel && lev >= zap.DebugLevel
	})

	errPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriteSyncer, infoPriority),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdin), debugPriority),
		zapcore.NewCore(encoder, errWriteSyncer, errPriority),
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

func initInfoLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.C.Log.Filename,
		MaxSize:    conf.C.Log.MaxSize,    // 日志文件大小 单位：MB
		MaxBackups: conf.C.Log.MaxBackups, // 备份数量
		MaxAge:     conf.C.Log.MaxAge,     // 备份时间 单位：天
		Compress:   true,                  // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

func initErrLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.C.LogErr.Filename,
		MaxSize:    conf.C.LogErr.MaxSize,    // 日志文件大小 单位：MB
		MaxBackups: conf.C.LogErr.MaxBackups, // 备份数量
		MaxAge:     conf.C.LogErr.MaxAge,     // 备份时间 单位：天
		Compress:   true,                     // 是否压缩
	}

	return zapcore.AddSync(lumberJackLogger)
}
