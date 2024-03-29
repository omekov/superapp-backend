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
	SugarLogger    *zap.SugaredLogger
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

	l.SugarLogger = logger.Sugar()
	if err := l.SugarLogger.Sync(); err != nil && errors.Is(err, syscall.ENOTTY) {
		l.SugarLogger.Error(err)
	}
}

// Debug ...
func (l *APILogger) Debug(args ...interface{}) {
	l.SugarLogger.Debug(args...)
}

// Debugf ...
func (l *APILogger) Debugf(template string, args ...interface{}) {
	l.SugarLogger.Debugf(template, args...)
}

// Info ...
func (l *APILogger) Info(args ...interface{}) {
	l.SugarLogger.Info(args...)
}

// Infof ...
func (l *APILogger) Infof(template string, args ...interface{}) {
	l.SugarLogger.Infof(template, args...)
}

// Warn ...
func (l *APILogger) Warn(args ...interface{}) {
	l.SugarLogger.Warn(args...)
}

// Warnf ...
func (l *APILogger) Warnf(template string, args ...interface{}) {
	l.SugarLogger.Warnf(template, args...)
}

// Error ...
func (l *APILogger) Error(args ...interface{}) {
	l.SugarLogger.Error(args...)
}

// Errorf ...
func (l *APILogger) Errorf(template string, args ...interface{}) {
	l.SugarLogger.Errorf(template, args...)
}

// DPanic ...
func (l *APILogger) DPanic(args ...interface{}) {
	l.SugarLogger.DPanic(args...)
}

// DPanicf ...
func (l *APILogger) DPanicf(template string, args ...interface{}) {
	l.SugarLogger.DPanicf(template, args...)
}

// Panic ...
func (l *APILogger) Panic(args ...interface{}) {
	l.SugarLogger.Panic(args...)
}

// Panicf ...
func (l *APILogger) Panicf(template string, args ...interface{}) {
	l.SugarLogger.Panicf(template, args...)
}

// Fatal ...
func (l *APILogger) Fatal(args ...interface{}) {
	l.SugarLogger.Fatal(args...)
}

// Fatalf ...
func (l *APILogger) Fatalf(template string, args ...interface{}) {
	l.SugarLogger.Fatalf(template, args...)
}
