package cmd

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gsadism/open-admin/conf"
	"github.com/gsadism/open-admin/core/db"
	"github.com/gsadism/open-admin/osv"
	"github.com/gsadism/open-admin/pkg/crypto/open_aes"
	"github.com/gsadism/open-admin/pkg/crypto/open_rsa"
	"github.com/gsadism/open-admin/pkg/file"
	"github.com/gsadism/open-admin/pkg/image"
	"github.com/gsadism/open-admin/pkg/next/snowflake"
	"github.com/gsadism/open-admin/pkg/storage"
	"github.com/spf13/viper"
	"log"
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

func Minio(web, host string, port int, username, password string, token string, ssl bool) (*storage.Minio, string) {
	if pwd, err := open_aes.DecryptECB(password, conf.SECRET_KEY); err != nil {
		log.Fatal(err)
		return nil, ""
	} else {
		if client, err := storage.NewMinioClient(fmt.Sprintf("%s:%d", host, port), username, pwd, token, ssl); err != nil {
			log.Fatal(err)
			return nil, ""
		} else {
			return client, web
		}
	}
}

func GLOBAL(v *viper.Viper) {
	osv.STATIC = func() string {
		if dir, err := file.Abs(v.GetString("website.static_dir")); err != nil {
			return ""
		} else {
			return dir
		}
	}()

}

func OSV(v *viper.Viper) {
	GLOBAL(v)

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

	// snow flake
	osv.SnowFlake.Init(snowflake.New(v.GetInt64("snowflake.machine"), v.GetInt64("snowflake.service")))

	// gg
	osv.Image.Init(image.NewGG(func() string {
		if dir, err := file.Abs(v.GetString("website.tft")); err != nil {
			log.Fatal(err)
			return ""
		} else {
			return dir
		}
	}()), osv.Redis.Client())

	// minio
	osv.Minio.Init(Minio(
		v.GetString("minio.web"),
		v.GetString("minio.host"),
		v.GetInt("minio.port"),
		v.GetString("minio.username"),
		v.GetString("minio.password"),
		v.GetString("minio.token"),
		v.GetBool("minio.ssl"),
	))
}
