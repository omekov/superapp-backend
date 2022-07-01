package logger

import (
	"errors"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger - methods interface
type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

// APILogger ...
type APILogger struct {
	level          string
	sugarLogger    *zap.SugaredLogger
	serverMode     string
	loggerEncoding string
}

// NewAPILogger ...
func NewAPILogger(level string) *APILogger {
	return &APILogger{level: level}
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *APILogger) getLoggerLevel(lvl string) zapcore.Level {
	level, exist := loggerLevelMap[lvl]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// InitLogger ...
func (l *APILogger) InitLogger() {
	logLevel := l.getLoggerLevel(l.level)

	logWriter := zapcore.AddSync(os.Stderr)

	var encoderCfg zapcore.EncoderConfig
	if l.serverMode == "Development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	if l.loggerEncoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil && errors.Is(err, syscall.ENOTTY) {
		l.sugarLogger.Error(err)
	}
}

// Debug ...
func (l *APILogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Debugf ...
func (l *APILogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

// Info ...
func (l *APILogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof ...
func (l *APILogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Warn ...
func (l *APILogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// Warnf ...
func (l *APILogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

// Error ...
func (l *APILogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

// Errorf ...
func (l *APILogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

// DPanic ...
func (l *APILogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

// DPanicf ...
func (l *APILogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

// Panic ...
func (l *APILogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

// Panicf ...
func (l *APILogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

// Fatal ...
func (l *APILogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

// Fatalf ...
func (l *APILogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
