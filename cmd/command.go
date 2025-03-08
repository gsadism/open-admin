package cmd

import (
	"crypto/aes"
	"errors"
	"fmt"
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/pkg/crypto/open_aes"
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

func OpenAdminSecretGenerateKeyCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "secret-generate-key",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if key, err := open_aes.Secret(2 * aes.BlockSize); err != nil {
				log.Fatal(err.Error())
			} else {
				log.Println(fmt.Sprintf("SECRET_KEY: %s", key))
			}
		},
	}
	return c
}

func OpenAdminEncipherCommand() *cobra.Command {
	var Password string
	c := &cobra.Command{
		Use:  "encipher",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if Password == "" {
				log.Fatal("lease configure the text to be encrypted correctly. <-c Encrypted text>")
			}
			if cipher, err := open_aes.EncryptECB(Password, conf.SECRET_KEY); err != nil {
				log.Fatal(err.Error())
			} else {
				fmt.Println(fmt.Sprintf("Encrypted Text: %s", cipher))
			}
		},
	}
	c.PersistentFlags().StringVarP(&Password, "password", "p", "", " encrypt text.")
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

	c.AddCommand(OpenAdminRsaGenerateKeyCommand(), OpenAdminSecretGenerateKeyCommand(), OpenAdminEncipherCommand())

	return c
}

func Execute(dir string) {
	_ = os.Setenv("OPEN_ADMIN_ROOT", dir)

	if err := OpenAdminServerCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
