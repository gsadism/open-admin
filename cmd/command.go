package cmd

import (
	"crypto/aes"
	"errors"
	"fmt"
	"github.com/gsadism/open-admin/core/debug"
	"github.com/gsadism/open-admin/pkg/crypto/open_aes"
	"github.com/gsadism/open-admin/pkg/crypto/open_rsa"
	"github.com/gsadism/open-admin/pkg/crypto/open_salt"
	"github.com/gsadism/open-admin/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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

// SuperAdminPasswordCommand : 生成 settings.go 文件下 SUPER_ADMIN_PASSWORD 超级管理员用户密码
func SuperAdminPasswordCommand() *cobra.Command {
	var Password string
	c := &cobra.Command{
		Use:  "super",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if Password == "" {
				debug.ErrorE("the password cannot be empty")
			} else {
				if pwd, err := open_salt.Encrypt(Password); err != nil {
					debug.ErrorE(err.Error())
				} else {
					if err := utils.WriteSettingsVal(filepath.Join(os.Getenv("OPEN_ADMIN_ROOT"), "conf", "settings.go"), "SUPER_ADMIN_PASSWORD", pwd); err != nil {
						debug.ErrorE(err.Error())
					} else {
						debug.Debug("writing the super admin password succeeds")
					}
				}
			}
		},
	}
	c.PersistentFlags().StringVarP(&Password, "password", "p", "", "super password")
	return c
}

// SecretCommand : 生成 settings.go 文件下 SECRET_KEY 秘钥
func SecretCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "secret",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if key, err := open_aes.Secret(2 * aes.BlockSize); err != nil {
				debug.Error(err.Error())
			} else {
				if err := utils.WriteSettingsVal(filepath.Join(os.Getenv("OPEN_ADMIN_ROOT"), "conf", "settings.go"), "SECRET_KEY", key); err != nil {
					debug.ErrorE(err.Error())
				} else {
					debug.Debug("writing the secret key succeeds")
				}
				//debug.Debug(fmt.Sprintf("SECRET_KEY: %s", key))
			}
		},
	}
	return c
}

// IVCommand : 生成AES偏移
func IVCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "iv",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if iv, err := open_aes.Secret(aes.BlockSize); err != nil {
				debug.Error(err.Error())
			} else {
				if err := utils.WriteSettingsVal(filepath.Join(os.Getenv("OPEN_ADMIN_ROOT"), "conf", "settings.go"), "SECRET_IV", iv); err != nil {
					debug.ErrorE(err.Error())
				} else {
					debug.Debug("writing the secret iv succeeds")
				}
			}
		},
	}
	return c
}

// RSACommand : 生成RSA 秘钥
func RSACommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "rsa",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			prv, pub := open_rsa.Generate(open_rsa.MEDIUM)
			PublicKey, err := open_rsa.PublicKeyWithString(pub)
			if err != nil {
				debug.ErrorE(err.Error())
			}
			PrivateKey, err := open_rsa.PrivateKeyWithString(prv)
			if err != nil {
				debug.ErrorE(err.Error())
			}
			if err := utils.WriteSettingsRsa(filepath.Join(os.Getenv("OPEN_ADMIN_ROOT"), "conf", "settings.go"), PrivateKey, PublicKey); err != nil {
				debug.ErrorE(err.Error())
			} else {
				debug.Debug("writing the rsa key succeeds")
			}
		},
	}
	return c
}

func Command() *cobra.Command {
	var ConfigFilePath string

	c := &cobra.Command{
		Use:  "server",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	c.AddCommand(
		SecretCommand(),
		SuperAdminPasswordCommand(),
		IVCommand(),
		RSACommand(),
	)
	c.PersistentFlags().StringVarP(&ConfigFilePath, "config", "c", "", "config file path")
	return c
}

func Execute(dir string) {
	_ = os.Setenv("OPEN_ADMIN_ROOT", dir)
	if err := Command().Execute(); err != nil {
		debug.ErrorE(err.Error())
	}
}
