package cmd

import (
	"errors"
	"fmt"
	"github.com/gsadism/open-admin/core"
	"github.com/gsadism/open-admin/logging"
	"github.com/gsadism/open-admin/pkg/object"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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

// folder : 获取路径
func folder(path string) string {
	if !filepath.IsAbs(path) {
		if d, err := filepath.Abs(path); err != nil {
			core.Exit(err.Error())
		} else {
			path = d
		}
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// 创建路径
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				core.Exit(err.Error())
			}
		}
	}
	return path
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
			} else {
				if v, err := readApplicationFile(flags.ConfFilePath); err != nil {
					core.Exit(err.Error())
				} else {
					logging.ReplaceGlobals(logging.New().SetSkip(2).File(
						folder(v.GetString("logger.file.path")),
						object.Default[string](v.GetString("logger.file.name"), "open-admin.log"),
						v.GetString("logger.file.level"),
						v.GetInt("logger.file.max-size"),
						v.GetInt("logger.file.max-age"),
						v.GetInt("logger.file.max-backups"),
						v.GetBool("logger.file.compress"),
					).R())

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
