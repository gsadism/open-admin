package cmd

import (
	"crypto/aes"
	"fmt"
	"github.com/gsadism/open-admin/core/debug"
	"github.com/gsadism/open-admin/pkg/crypto/_aes"
	"github.com/gsadism/open-admin/pkg/crypto/passlib"
	"github.com/spf13/cobra"
)

func OpenAdminSuper() *cobra.Command {
	var Password string
	c := &cobra.Command{
		Use:  "super",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if Password == "" {
				debug.Error("the password cannot be empty")
			} else {
				if pwd, err := passlib.Encrypt(Password); err != nil {
					debug.Error(err.Error())
				} else {
					debug.Debug(fmt.Sprintf("super password: %s", pwd))
				}
			}
		},
	}
	c.PersistentFlags().StringVarP(&Password, "password", "p", "", "super password")
	return c
}

func OpenAdminSecret() *cobra.Command {
	c := &cobra.Command{
		Use:  "secret",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if key, err := _aes.Secret(2 * aes.BlockSize); err != nil {
				debug.Error(err.Error())
			} else {
				debug.Debug(fmt.Sprintf("SECRET_KEY: %s", key))
			}
		},
	}
	return c
}

func OpenAdminServer() *cobra.Command {
	var ConfigFilePath string

	c := &cobra.Command{
		Use:  "server",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	c.PersistentFlags().StringVarP(&ConfigFilePath, "config", "c", "", "config file path")

	c.AddCommand(OpenAdminSuper(), OpenAdminSecret())

	return c
}
