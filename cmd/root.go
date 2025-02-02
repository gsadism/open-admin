package cmd

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/core"
	"github.com/gsadism/open-admin/depends/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func readApplicationFile(fileName string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(fileName)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New(fmt.Sprintf("%s file not found", fileName))
		} else {
			return nil, err
		}
	}
	return v, nil
}

func Execute() {
	var ConfigFilePath string
	root := &cobra.Command{
		Use:  "server",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if v, err := readApplicationFile(ConfigFilePath); err != nil {
				logging.Fatal(err.Error())
			} else {
				sv := core.NewServer(&core.Options{
					Debug: v.GetBool("server.debug"),
					Host:  v.GetString("server.host"),
					Port:  v.GetInt("server.port"),
				}).Middleware(conf.MIDDLEWARE...).Routers(func() []func(*gin.RouterGroup) {
					s := make([]func(*gin.RouterGroup), 0)
					for _, app := range conf.APPS {
						s = append(s, app.Routers...)
					}
					return s
				}()...)

				sv.ListenAndServe()
			}
		},
	}
	root.PersistentFlags().StringVarP(&ConfigFilePath, "config", "c", "", "config file path")

	if err := root.Execute(); err != nil {
		logging.Fatal(err.Error())
	}
}
