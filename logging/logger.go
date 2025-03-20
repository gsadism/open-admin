package logging

import (
	"fmt"
	"github.com/gsadism/open-admin/logging/encoder"
	"github.com/gsadism/open-admin/pkg/object"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

type Logger struct {
	log *zap.Logger

	skip int
	once sync.Once

	cores []zapcore.Core
}

func New() *Logger {
	l := &Logger{
		cores: make([]zapcore.Core, 0),
	}
	l.reset()
	return l
}

func (l *Logger) SetSkip(skip int) *Logger {
	l.skip = skip
	return l
}

func (l *Logger) File(
	Dir string,
	Filename string,
	Level string,
	MaxSize int,
	MaxAge int,
	MaxBackup int,
	Compress bool,
) *Logger {
	l.cores = append(l.cores, zapcore.NewCore(zapcore.NewJSONEncoder(l.config()), zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(Dir, Filename),
		MaxSize:    MaxSize,
		MaxBackups: MaxBackup,
		MaxAge:     MaxAge,
		Compress:   Compress,
	}), parseLogLevel(object.Default[string](Level, "error"))))
	return l
}

func (l *Logger) config() zapcore.EncoderConfig {
	Config := zap.NewProductionEncoderConfig()
	Config.LineEnding = zapcore.DefaultLineEnding
	Config.EncodeDuration = zapcore.SecondsDurationEncoder
	Config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	Config.EncodeLevel = zapcore.CapitalLevelEncoder // 输出level序列化为全大写字符串，如 INFO DEBUG ERROR
	Config.EncodeCaller = zapcore.FullCallerEncoder
	return Config
}

func (l *Logger) reset() {
	l.cores = make([]zapcore.Core, 0)
	l.cores = append(l.cores, zapcore.NewCore(&encoder.Console{
		Encoder: zapcore.NewConsoleEncoder(l.config()),
	}, zapcore.Lock(os.Stdout), parseLogLevel(object.Default[string](os.Getenv("OPEN_ADMIN_LOG_LEVEL"), "debug"))))
}

func (l *Logger) gc() {
	l.cores = nil
}

func (l *Logger) R() *Logger {
	l.once.Do(func() {
		l.log = zap.New(zapcore.NewTee(l.cores...), zap.AddCaller(), zap.AddCallerSkip(l.skip))
	})
	return l
}

func (l *Logger) Debug(msg any, fields ...zap.Field) { l.log.Debug(fmt.Sprint(msg), fields...) }
func (l *Logger) Info(msg any, fields ...zap.Field)  { l.log.Info(fmt.Sprint(msg), fields...) }
func (l *Logger) Warn(msg any, fields ...zap.Field)  { l.log.Warn(fmt.Sprint(msg), fields...) }
func (l *Logger) Error(msg any, fields ...zap.Field) { l.log.Error(fmt.Sprint(msg), fields...) }

//func (l *Logger) Fatal(msg any, fields ...zap.Field) { l.log.Fatal(fmt.Sprint(msg), fields...) }
