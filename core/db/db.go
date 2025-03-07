package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Client(
	Driver string,
	Host string,
	Port int,
	UserName string,
	Password string,
	DBName string,
	Charset string,
	MaxOpenConn int,
	MaxIdleConn int,
	MaxIdleLifeTime int,
	MaxLifetime int,
	Config *gorm.Config,
) (*gorm.DB, error) {
	var Dialector gorm.Dialector

	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		UserName,
		Password,
		Host,
		Port,
		DBName,
		Charset,
	)
	switch strings.ToLower(Driver) {
	case "mysql":
		Dialector = mysql.Open(dns)
	default:
		return nil, fmt.Errorf("driver %s not supported", Driver)
	}

	if client, err := gorm.Open(Dialector, Config); err != nil {
		return nil, err
	} else {
		if client != nil {

			_db, _ := client.DB()
			_db.SetMaxOpenConns(MaxOpenConn)
			_db.SetMaxIdleConns(MaxIdleConn)
			_db.SetConnMaxIdleTime(time.Duration(MaxIdleLifeTime))
			_db.SetConnMaxLifetime(time.Duration(MaxLifetime))
			return client, nil
		} else {
			return nil, errors.New("db is nil")
		}
	}
}
