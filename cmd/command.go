package cmd

import (
	"errors"
	"fmt"
	"github.com/gsadism/open-admin/pkg/crypto/open_rsa"
	"github.com/gsadism/open-admin/pkg/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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

func OpenAdminRsaGenerateKeyCommand() *cobra.Command {
	var Dir string
	c := &cobra.Command{
		Use:  "rsa-generate-key",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if dir, err := file.Folder(Dir); err != nil {
				log.Fatal(err.Error())
			} else {
				prv, pub := open_rsa.Generate(open_rsa.MEDIUM)
				// 写入pem
				if err := open_rsa.WritePKIXPublicKeyFile(pub, filepath.Join(dir, "public_key.pem"), nil); err != nil {
					log.Fatal(err.Error())
				}
				if err := open_rsa.WritePKCS8PrivateKeyFile(prv, filepath.Join(dir, "private_key.pem"), nil); err != nil {
					log.Fatal(err.Error())
				}
				log.Println("rsa generate key success")
			}
		},
	}
	c.PersistentFlags().StringVarP(&Dir, "dir", "d", ".", "Directory to serve files from")
	return c
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

	c.AddCommand(OpenAdminRsaGenerateKeyCommand())

	return c
}

func Execute(dir string) {
	_ = os.Setenv("OPEN_ADMIN_ROOT", dir)

	if err := OpenAdminServerCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
