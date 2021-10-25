package logger

import (
	"os"

	zapLog "github.com/liuxiong332/kratos-starter/logger/zap"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() *zapLog.Logger {
	fileOut := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/elk.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     3, // days
	})
	stdout := zapcore.AddSync(os.Stdout)
	w := zap.CombineWriteSyncers(stdout, fileOut)

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), w, zap.InfoLevel)

	zlog := zap.New(core, zap.ErrorOutput(zapcore.AddSync(os.Stderr)), zap.AddCaller(), zap.AddCallerSkip(3), zap.AddStacktrace(zap.ErrorLevel))

	return zapLog.NewLogger(zlog)
}
