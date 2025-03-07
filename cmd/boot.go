package cmd

import (
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/core/server"
	"github.com/spf13/viper"
	"runtime"
)

func Server(v *viper.Viper) {

	srv := server.New(server.NewConfig().
		SetDebug(v.GetBool("server.debug")).
		SetHost(v.GetString("server.host")).
		SetPort(v.GetInt("server.port"))).
		Middleware(conf.MIDDLEWARE...).
		Files("favicon.ico", v.GetString("website.favicon")).
		Files("robots.txt", v.GetString("website.robots"))

	// 强行GC
	runtime.GC()
	srv.ListenAndServer()
}
