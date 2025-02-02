package logging

import "go.uber.org/zap/zapcore"

type Handler string

const (
	Stream Handler = "stream"
	File   Handler = "file"
)

type Entry = zapcore.Entry

type StreamOptions struct {
	Level
	Formatter func(ent Entry) (string, error)
}

type FileOptions struct {
	FileName string // FileName 日志文件路径
	Level    Level  // Level 日志记录级别
	// MaxSize 日志文件达到的最大大小(以MB为单位)
	// 当日志文件大小达到此值时,自动将日志文件备份并创建一个新的日志文件继续记录
	// 默认值为100MB
	MaxSize int
	// MaxBackups 保留的旧日志文件的最大数量
	// 当达到此数量时,最旧的日志文件将被删除
	// 默认值为0,表示不限制备份数量‌
	MaxBackups int
	// MaxAge 旧日志文件的最大保存天数
	// 超过此天数的旧日志文件将被删除
	// 默认值为0, 表示不限制保存天数‌
	MaxAge int
	// Compress 否对备份的日志文件进行压缩
	// 默认值为false,表示不压缩备份文
	Compress bool
}

type Options struct {
	Handler []Handler
	Stream  StreamOptions
	File    FileOptions
}
