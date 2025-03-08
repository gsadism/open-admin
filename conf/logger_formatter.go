package conf

import (
	"bytes"
	"fmt"
	"github.com/gsadism/open-admin/logging"
	"path/filepath"
	"strings"
	"text/template"
)

func Formatter(ent logging.Entry) (string, error) {
	type Data struct {
		Time    string
		Level   string
		Message string
		Pc      string
		Caller  string
	}
	// 按照 时间 [级别] 内容 <文件.函数> 调用栈 输出日志
	tmpl := "[{{.Time}}] |{{.Level}}| {{.Message}} {{.Pc}} {{.Caller}}" // 模板
	if t, err := template.New("stream").Parse(tmpl); err != nil {
		return "", err
	} else {
		wr := bytes.NewBuffer(nil)
		if err := t.Execute(wr, &Data{
			Time:    ent.Time.Format("2006-01-02 15:04:05"),
			Level:   strings.ToUpper(ent.Level.String() + strings.Repeat(" ", 5-len(ent.Level.CapitalString()))),
			Message: ent.Message,
			Pc:      fmt.Sprintf("<%s-%s>", filepath.Base(ent.Caller.File), strings.Split(ent.Caller.Function, ".")[len(strings.Split(ent.Caller.Function, "."))-1]),
			Caller:  ent.Caller.String(),
		}); err != nil {
			return "", err
		} else {
			tx := fmt.Sprintf("\033[%dm\033[%dm%v\033[0m\n",
				0, // 默认显示模式
				func() int { // 不同日志的显示颜色
					if ent.Level == logging.DebugLevel {
						return 32
					} else if ent.Level == logging.InfoLevel {
						return 37
					} else if ent.Level == logging.WarnLevel {
						return 33
					} else if ent.Level == logging.ErrorLevel {
						return 31
					} else if ent.Level == logging.FatalLevel {
						return 31
					} else {
						return 37
					}
				}(),
				wr.String(), // 日志内容
			)
			return tx, nil
		}
	}
}
