package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

var once sync.Once

func init() {
	once.Do(func() {
		fileWriter := zapcore.AddSync(
			&lumberjack.Logger{
				Filename:   "./logs/logs.json",
				MaxSize:    1,
				MaxAge:     30,
				LocalTime:  false,
				Compress:   false,
			},
		)
		stdoutWriter := zapcore.AddSync(os.Stdout)
		level := zapcore.InfoLevel

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(encoderCfg)
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, stdoutWriter, level),
			zapcore.NewCore(encoder, fileWriter, level),
		)
		
		Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(level))
	})
}
