package cmd

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/core/db"
	"github.com/gsadism/open-admin/core/server"
	"github.com/gsadism/open-admin/osv"
	"github.com/gsadism/open-admin/pkg/crypto/open_aes"
	"github.com/gsadism/open-admin/pkg/crypto/open_rsa"
	"github.com/gsadism/open-admin/pkg/file"
	"github.com/spf13/viper"
	"log"
	"runtime"
)

func RSA(PublicKeyPath, PrivateKeyPath string) (*rsa.PublicKey, *rsa.PrivateKey) {
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

func Redis(host string, port int, password string) *redis.Client {
	if pwd, err := open_aes.DecryptECB(password, conf.SECRET_KEY); err != nil {
		log.Fatal(err)
		return nil
	} else {
		rds := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: pwd,
			DB:       0,
		})
		if _, err := rds.Ping(context.TODO()).Result(); err != nil {
			log.Fatal(err.Error())
			return nil
		} else {
			return rds
		}
	}
}

func Server(v *viper.Viper) {
	// db
	if client, err := db.Client(
		v.GetString("database.driver"),
		v.GetString("database.host"),
		v.GetInt("database.port"),
		v.GetString("database.username"),
		func() string {
			if pwd, err := open_aes.DecryptECB(v.GetString("database.password"), conf.SECRET_KEY); err != nil {
				log.Fatal(err)
				return ""
			} else {
				return pwd
			}
		}(),
		v.GetString("database.db"),
		v.GetString("database.charset"),
		v.GetInt("database.max-open"),
		v.GetInt("database.max-idle"),
		v.GetInt("database.max-idle-time"),
		v.GetInt("database.max-life-time"),
		conf.GORM,
	); err != nil {
		log.Fatal(err)
	} else {
		osv.DB.Init(client)
	}

	// redis
	osv.Redis.Init(Redis(v.GetString("redis.host"), v.GetInt("redis.port"), v.GetString("redis.password")))

	// rsa
	osv.Rsa.Init(RSA(v.GetString("secret.public_key"), v.GetString("secret.private_key")))

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
