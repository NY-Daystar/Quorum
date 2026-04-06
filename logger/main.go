package logger

import (
	"os"
	"quorum/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	jack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	maxSize    int = 5  // 3 megabytes per files
	maxBackups int = 10 // 3 files before rotate
	maxAge     int = 15 // 15 days
)

// Log custom log based on uber.zap
type Log struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LumberJack *jack.Logger
	ZapLogger  *zap.Logger
}

// Init setup custom logger based on uber zap logger
func Init() *Log {
	log := newLogger(utils.GetLogsFile(), maxAge, maxBackups, maxAge)

	defer log.ZapLogger.Sync() // flushes buffer, if any

	log.Debug("Zap logger set",
		zap.String("path", log.Filename),
		zap.Int("filesize", log.MaxSize), zap.Int("backupfile", log.MaxBackups),
		zap.Int("fileage", log.MaxAge),
	)

	return log
}

func newLogger(filename string, maxSize, maxBackup, maxFile int) *Log {
	logger := &jack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxFile,
	}

	core := zapcore.NewTee(
		createFileLogger(logger),
		createConsoleLogger(),
	)

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Log{
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		LumberJack: logger,
		ZapLogger:  log,
	}
}

// createFileLogger return core logger to write logs into a file
func createFileLogger(lumber *jack.Logger) zapcore.Core {
	configZap := zap.NewProductionEncoderConfig()
	configZap.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(configZap),
		zapcore.AddSync(lumber),
		zapcore.DebugLevel,
	)
}

// createConsoleLogger return core logger to display logs in the console
func createConsoleLogger() zapcore.Core {
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

// Sugar allows to log with sugar function
func (log *Log) Sugar() *zap.SugaredLogger {
	return log.ZapLogger.Sugar()
}

// Debug allows to log with Debug function
func (log *Log) Debug(msg string, fields ...zap.Field) {
	log.ZapLogger.Debug(msg, fields...)
}
