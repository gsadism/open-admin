package logging

import (
	"fmt"
	"github.com/gsadism/open-admin/pkg/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
)

type Entry = zapcore.Entry

type Level = zapcore.Level

const (
	DebugLevel Level = zapcore.DebugLevel
	InfoLevel  Level = zapcore.InfoLevel
	WarnLevel  Level = zapcore.WarnLevel
	ErrorLevel Level = zapcore.ErrorLevel
	FatalLevel Level = zapcore.FatalLevel
)

type Logger struct {
	log  *zap.Logger
	skip int
	once sync.Once

	cores map[string]zapcore.Core
	dlt   zapcore.Core
}

func New() *Logger {
	l := &Logger{
		skip:  1,
		cores: make(map[string]zapcore.Core),
	}

	return l
}

func (l *Logger) SetSkip(skip int) *Logger {
	l.skip = skip
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

func (l *Logger) Default(Lv Level) *Logger {
	l.dlt = zapcore.NewCore(
		zapcore.NewConsoleEncoder(l.config()),
		zapcore.AddSync(os.Stdout),
		Lv,
	)
	return l
}

func (l *Logger) Console(Lv Level, Formatter func(ent Entry) (string, error)) *Logger {
	l.cores["console"] = zapcore.NewCore(&Console{
		Encoder:   zapcore.NewConsoleEncoder(l.config()),
		Formatter: Formatter,
	}, zapcore.Lock(os.Stdout), Lv)
	return l
}

func (l *Logger) File(
	Dir string,
	FileName string,
	Lv Level,
	MaxSize int,
	MaxBackups int,
	MaxAge int,
	Compress bool,
) *Logger {
	if file.Exists(Dir) {
		if err := os.MkdirAll(Dir, os.ModeDir); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}
	l.cores["file"] = zapcore.NewCore(zapcore.NewJSONEncoder(l.config()), zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(Dir, FileName),
		MaxSize:    MaxSize,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		Compress:   Compress,
	}), Lv)
	return l
}

func (l *Logger) gc() {
	l.cores = nil
}

func (l *Logger) R() *Logger {
	if len(l.cores) == 0 {
		if l.dlt == nil {
			l.Default(ErrorLevel)
			l.cores["console"] = l.dlt
		} else {
			l.cores["console"] = l.dlt
		}
	}
	cores := func() []zapcore.Core {
		sc := make([]zapcore.Core, 0)
		for _, v := range l.cores {
			sc = append(sc, v)
		}
		return sc
	}()

	l.once.Do(func() {
		l.log = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(l.skip))
	})

	l.gc()

	return l
}

func (l *Logger) Debug(msg any, fields ...zap.Field) { l.log.Debug(fmt.Sprint(msg), fields...) }
func (l *Logger) Info(msg any, fields ...zap.Field)  { l.log.Info(fmt.Sprint(msg), fields...) }
func (l *Logger) Warn(msg any, fields ...zap.Field)  { l.log.Warn(fmt.Sprint(msg), fields...) }
func (l *Logger) Error(msg any, fields ...zap.Field) { l.log.Error(fmt.Sprint(msg), fields...) }
func (l *Logger) Fatal(msg any, fields ...zap.Field) { l.log.Fatal(fmt.Sprint(msg), fields...) }
