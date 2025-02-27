package debug

import (
	"fmt"
	"github.com/gsadism/open-admin/pkg/color"
	"log"
	"os"
	"time"
)

// formatStr: logger format string
const formatStr = "%c[%d;%dm[%s] %v%c[0m"

var std = log.New(os.Stdout, "", log.Lmsgprefix)

func Debug(msg any) {
	_ = std.Output(2, fmt.Sprintf(formatStr, 0x1B, 0, 30+color.Yellow, time.Now().Format("2006-01-02 15:04:05"), msg, 0x1B))
}

func Error(msg any) {
	_ = std.Output(2, fmt.Sprintf(formatStr, 0x1B, 0, 30+color.Red, time.Now().Format("2006-01-02 15:04:05"), msg, 0x1B))
}

func ErrorE(msg any) {
	_ = std.Output(2, fmt.Sprintf(formatStr, 0x1B, 0, 30+color.Red, time.Now().Format("2006-01-02 15:04:05"), msg, 0x1B))
	os.Exit(1)
}
