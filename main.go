package main

import (
	"github.com/gsadism/open-admin/cmd"
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/depends/logging"
)

func main() {
	logging.ReplaceGlobals(logging.New(conf.LOGGING)) // 全局日志记录器初始化

	cmd.Execute()
}
