package debug

import (
	"fmt"
	"log"
	"os"
	"time"
)

//type color = int
//
//const (
//	Black      color = iota // 黑色
//	Red                     // 红色
//	Green                   // 绿色
//	Yellow                  // 黄色
//	Blue                    // 蓝色
//	Fuchsia                 // 紫红色
//	BluishBlue              // 青蓝色
//	White                   // 白色
//)

// formatStr: logger format string
const formatStr = "%c[%d;%dm[%s] %v%c[0m"

var std = log.New(os.Stdout, "", log.Lmsgprefix)

func Debug(msg any) {
	_ = std.Output(2, fmt.Sprintf(formatStr, 0x1B, 0, 33, time.Now().Format("2006-01-02 15:04:05"), msg, 0x1B))
}

func Error(msg any) {
	_ = std.Output(2, fmt.Sprintf(formatStr, 0x1B, 0, 31, time.Now().Format("2006-01-02 15:04:05"), msg, 0x1B))
}

func ErrorE(msg any) {
	_ = std.Output(2, fmt.Sprintf(formatStr, 0x1B, 0, 31, time.Now().Format("2006-01-02 15:04:05"), msg, 0x1B))
	os.Exit(1)
}
