package cmd

import (
	"errors"
	"fmt"
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/core"
	"github.com/gsadism/open-admin/pkg/object"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var flags = new(struct {
	ConfFilePath string // 配置文件路径
	Log          string // 控制台日志输出级别
})

// readApplicationFile : 加载application.yml配置文件
func readApplicationFile(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//if errors.Is(err, new(viper.ConfigFileNotFoundError)) {
			return nil, errors.New(fmt.Sprintf("[%s] file not found", path))
		} else {
			return nil, err
		}
	}
	return v, nil
}

// bindFlags : 绑定命令行参数
func bindFlags(c *cobra.Command) *cobra.Command {
	c.PersistentFlags().StringVarP(&flags.ConfFilePath, "conf", "", "", "config file path")
	c.PersistentFlags().StringVarP(&flags.Log, "log", "", "", "console log output level")
	return c
}

func command() *cobra.Command {
	c := &cobra.Command{
		Use:  "server",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// 日志记录器初始化
			if object.In[string](strings.ToLower(flags.Log), []string{"debug", "info", "warn", "error"}) {
				_ = os.Setenv("OPEN_ADMIN_LOG_LEVEL", strings.ToLower(flags.Log))
			}

			if flags.ConfFilePath == "" {
				// 使用默认配置
				srv := core.Default().
					Middleware(conf.MIDDLEWARE...).
					Routers(conf.ROUTERS...)
				srv.ListenAndServer()
			} else {
				if v, err := readApplicationFile(flags.ConfFilePath); err != nil {
					core.Exit(err.Error())
				} else {
					srv := core.New(v).
						Middleware(conf.MIDDLEWARE...).
						Routers(conf.ROUTERS...)
					srv.ListenAndServer()
				}
			}
		},
	}
	return bindFlags(c)
}

func Execute(dir string) {
	_ = os.Setenv("OPEN_ADMIN_ROOT", dir)
	if err := command().Execute(); err != nil {
		core.Exit(err.Error())
	}
}
