package cmd

import (
	"crypto/rsa"
	"fmt"
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/core/server"
	"github.com/gsadism/open-admin/osv"
	"github.com/gsadism/open-admin/pkg/crypto/open_rsa"
	"github.com/gsadism/open-admin/pkg/file"
	"github.com/spf13/viper"
	"log"
	"runtime"
)

func OsRSA(PublicKeyPath, PrivateKeyPath string) (*rsa.PublicKey, *rsa.PrivateKey) {
	if file.Exists(PublicKeyPath) && file.Exists(PrivateKeyPath) {
		pub, err := open_rsa.LoadPKIXPublicKeyWithFile(PublicKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		prv, err := open_rsa.LoadPKCS8PrivateKeyWithFile(PrivateKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		// 检查 密钥对是否匹配
		text := "hello world"
		if Encrypt, err := open_rsa.EncryptWithPKCS1v15(pub, text); err != nil {
			log.Fatal(err)
		} else {
			// 解密
			if _, err := open_rsa.DecryptWithPKCS1v15(prv, Encrypt); err != nil {
				log.Fatal(err)
			}
		}
		return pub, prv
	} else {
		log.Fatal(fmt.Errorf("secret path ont exists"))
		return nil, nil
	}
}

func Server(v *viper.Viper) {
	// rsa
	osv.Rsa.Init(OsRSA(v.GetString("secret.public_key"), v.GetString("secret.private_key")))

	srv := server.New(server.NewConfig().
		SetDebug(v.GetBool("server.debug")).
		SetHost(v.GetString("server.host")).
		SetPort(v.GetInt("server.port"))).
		Middleware(conf.MIDDLEWARE...).
		Files("favicon.ico", v.GetString("website.favicon")).
		Files("robots.txt", v.GetString("website.robots"))

	// 强行GC
	runtime.GC()
	srv.ListenAndServer()
}
