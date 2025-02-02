package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

type Level = zapcore.Level

const (
	DebugLevel Level = zapcore.DebugLevel
	InfoLevel  Level = zapcore.InfoLevel
	WarnLevel  Level = zapcore.WarnLevel
	ErrorLevel Level = zapcore.ErrorLevel
)

type Logger struct {
	log  *zap.Logger
	Skip int
	once sync.Once
}

func New(opt *Options) *Logger {
	l := new(Logger)
	l.Skip = 2

	l.load(opt)
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

func (l *Logger) stream(lv Level, formatter func(ent Entry) (string, error)) zapcore.Core {
	return zapcore.NewCore(&StreamHandler{
		Encoder:   zapcore.NewConsoleEncoder(l.config()),
		Formatter: formatter,
	}, zapcore.Lock(os.Stdout), lv)
}

func (l *Logger) file(
	FileName string,
	Lv Level,
	MaxSize int,
	MaxBackups int,
	MaxAge int,
	Compress bool,
) zapcore.Core {
	return zapcore.NewCore(zapcore.NewJSONEncoder(l.config()), zapcore.AddSync(&lumberjack.Logger{
		Filename:   FileName,
		MaxSize:    MaxSize,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		Compress:   Compress,
	}), Lv)
}

func (l *Logger) console() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(l.config()),
		zapcore.AddSync(os.Stdout),
		ErrorLevel,
	)
}

func (l *Logger) load(opt *Options) {
	cores := make([]zapcore.Core, 0)
	for _, h := range opt.Handler {
		if h == Stream {
			cores = append(cores, l.stream(opt.Stream.Level, opt.Stream.Formatter))
		} else if h == File {
			cores = append(cores, l.file(
				opt.File.FileName,
				opt.File.Level,
				opt.File.MaxSize,
				opt.File.MaxBackups,
				opt.File.MaxAge,
				opt.File.Compress,
			))
		}
	}

	if len(cores) == 0 {
		cores = append(cores, l.console())
	}

	l.once.Do(func() {
		l.log = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(l.Skip))
	})
}

func (l *Logger) Debug(msg any, fields ...zap.Field) { l.log.Debug(fmt.Sprint(msg), fields...) }
func (l *Logger) Info(msg any, fields ...zap.Field)  { l.log.Info(fmt.Sprint(msg), fields...) }
func (l *Logger) Warn(msg any, fields ...zap.Field)  { l.log.Warn(fmt.Sprint(msg), fields...) }
func (l *Logger) Error(msg any, fields ...zap.Field) { l.log.Error(fmt.Sprint(msg), fields...) }
func (l *Logger) Fatal(msg any, fields ...zap.Field) { l.log.Fatal(fmt.Sprint(msg), fields...) }
