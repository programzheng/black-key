package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/programzheng/black-key/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Sync() error
}

type zapLogger struct {
	*zap.Logger
	output *os.File
}

func New(logSystem string) Logger {
	// 預設輸出至 stdout
	output := os.Stdout

	if logSystem == "" {
		logSystem = config.Cfg.GetString("LOG_SYSTEM")
	}

	// 如果 LOG_SYSTEM 為 "file"，則將輸出寫入 log.txt 檔案中
	if logSystem == "file" {
		f, err := os.OpenFile(
			fmt.Sprintf("%s/%s", config.Cfg.GetString("LOG_PATH"), getLogFileNameByMode()),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
		)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		output = f
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(output), zap.NewAtomicLevelAt(zap.DebugLevel))

	if config.Cfg.GetString("LOG_ENV") == "development" {
		//AddCallerSkip(2) is used to skip the caller information of the logging function itself.
		logger, err := zap.NewDevelopment(zap.AddStacktrace(zapcore.DebugLevel), zap.AddCaller(), zap.AddCallerSkip(2))

		if err != nil {
			panic(fmt.Sprintf("Failed to create logger: %v", err))
		}
		return &zapLogger{
			Logger: logger,
			output: output,
		}

	} else {
		//AddCallerSkip(2) is used to skip the caller information of the logging function itself.
		logger := zap.New(core, zap.AddStacktrace(zapcore.DebugLevel), zap.AddCaller(), zap.AddCallerSkip(2))
		return &zapLogger{
			Logger: logger,
			output: output,
		}
	}
}

func NewWithLogger(logger *zap.Logger) Logger {
	return &zapLogger{logger, os.Stdout}
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *zapLogger) Sync() error {
	return l.Logger.Sync()
}

func Debug(msg string, fields ...zap.Field) {
	New("").Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	New("").Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	New("").Error(msg, fields...)
}

func Sync() error {
	return New("").Sync()
}
