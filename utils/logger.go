package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// CreateFileLogger return core logger to write logs into a file
func CreateFileLogger(lumber *lumberjack.Logger) zapcore.Core {
	configZap := zap.NewProductionEncoderConfig()
	configZap.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(configZap),
		zapcore.AddSync(lumber),
		zapcore.DebugLevel,
	)
}

// CreateConsoleLogger return core logger to display logs in the console
func CreateConsoleLogger() zapcore.Core {
	configConsole := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "",
		TimeKey:       "",
		NameKey:       "",
		CallerKey:     "",
		FunctionKey:   "",
		StacktraceKey: "",

		LineEnding: zapcore.DefaultLineEnding,
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(configConsole),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
}
