package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

func ReadApplicationFile(fileName string) (*viper.Viper, error) {
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

func OpenAdminServerCommand() *cobra.Command {
	var ConfigFilePath string
	c := &cobra.Command{
		Use:  "server",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if v, err := ReadApplicationFile(ConfigFilePath); err != nil {
				log.Fatal(err)
			} else {
				Server(v)
			}
		},
	}
	c.PersistentFlags().StringVarP(&ConfigFilePath, "config", "c", "", "config file path")

	c.AddCommand()

	return c
}

func Execute(dir string) {
	_ = os.Setenv("OPEN_ADMIN_ROOT", dir)

	if err := OpenAdminServerCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
