package main

import (
	"github.com/gsadism/open-admin/cmd"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	cmd.Execute(wd)
}
