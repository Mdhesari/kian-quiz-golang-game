package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

var once sync.Once

func L() *zap.Logger {
	return logger
}

func init() {
	once.Do(func() {
		fileWriter := zapcore.AddSync(
			&lumberjack.Logger{
				Filename:  "./logs/logs.json",
				MaxSize:   1,
				MaxAge:    30,
				LocalTime: false,
				Compress:  false,
			},
		)
		stdoutWriter := zapcore.AddSync(os.Stdout)
		level := zapcore.InfoLevel

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdoutWriter, level),
			zapcore.NewCore(jsonEncoder, fileWriter, level),
		)

		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(level))
	})
}
