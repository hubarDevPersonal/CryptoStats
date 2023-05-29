package log

import (
	"fmt"

	"go.uber.org/zap"
)

type Field = zap.Field

type Logger struct {
	zLog *zap.Logger
}

var (
	Default *Logger
	Nop     *Logger
)

func init() {
	Default, _ = New("dev", "debug")
	Nop = &Logger{zLog: zap.NewNop()}
}

func New(name, level string) (*Logger, error) {
	var err error

	if name == "" {
		name = "prod"
	}

	var logCfg zap.Config

	switch name {
	case "dev":
		logCfg = zap.NewDevelopmentConfig()
	case "prod":
		logCfg = zap.NewProductionConfig()
	default:
		return nil, fmt.Errorf("invalid log configuration")
	}
	logCfg.OutputPaths = []string{"stdout"}

	al := zap.NewAtomicLevel()
	err = al.UnmarshalText([]byte(level))
	if err == nil {
		logCfg.Level = al
	}

	logger, err := logCfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("error building log: %w", err)
	}

	return &Logger{
		logger,
	}, nil
}

func (l *Logger) ZapLogger() *zap.Logger {
	return l.zLog
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.zLog.Fatal(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.zLog.Error(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.zLog.Info(msg, fields...)
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Err(err error) Field {
	return zap.Error(err)
}
