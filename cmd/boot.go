package cmd

import (
	"context"
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/core/server"
	"github.com/gsadism/open-admin/osv"
	"github.com/spf13/viper"
	"runtime"
)

func Server(v *viper.Viper) {

	OSV(v)

	srv := server.New(server.NewConfig().
		SetDebug(v.GetBool("server.debug")).
		SetHost(v.GetString("server.host")).
		SetPort(v.GetInt("server.port"))).
		SetPoolSize(v.GetInt("server.pool")).
		Middleware(conf.MIDDLEWARE...).
		Files("favicon.ico", v.GetString("website.favicon")).
		Files("robots.txt", v.GetString("website.robots")).
		Routers(conf.ROUTERS...).
		AutoMigrate(osv.DB.WithContext(context.WithValue(context.TODO(), "sudo", true)), v.GetBool("server.auto_migrate"), conf.MODELS...).
		Hoos(conf.HOOKS...)

	// 强行GC
	runtime.GC()
	srv.ListenAndServer()
}
