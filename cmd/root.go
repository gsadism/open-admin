package cmd

import (
	"github.com/gsadism/open-admin/core/debug"
	"os"
)

func Execute(dir string) {
	_ = os.Setenv("OPEN_ADMIN_ROOT", dir)
	if err := OpenAdminServer().Execute(); err != nil {
		debug.ErrorE(err.Error())
	}
}
