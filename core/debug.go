package core

import (
	"fmt"
	"log"
	"os"
	"time"
)

// formatStr: logger format string
const formatStr = "%c[%d;%dm[%s] %v%c[0m"

var std = log.New(os.Stdout, "", log.Lmsgprefix)

func Debug(msg any) {
	_ = std.Output(2, fmt.Sprintf(formatStr, 0x1B, 0, 31, time.Now().Format("2006-01-02 15:04:05"), msg, 0x1B))
}

// Exit : 打印内容并中断程序
func Exit(msg string) {
	_ = std.Output(2, fmt.Sprintf(formatStr, 0x1B, 0, 31, time.Now().Format("2006-01-02 15:04:05"), msg, 0x1B))
	os.Exit(1)
}
