package logging

var (
	global = New()
)

func ReplaceGlobals(logger *Logger) { global = logger }
func Debug(msg any)                 { global.Debug(msg) }
func Info(msg any)                  { global.Info(msg) }
func Warn(msg any)                  { global.Warn(msg) }
func Error(msg any)                 { global.Error(msg) }
func Fatal(msg any)                 { global.Fatal(msg) }
