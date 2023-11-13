package db

import (
	"fmt"

	"github.com/FogMeta/libra-os/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dc := config.Conf().DataBase
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dc.User, dc.Password, dc.Host, dc.Port, dc.Database)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if dc.Debug {
		DB = DB.Debug()
	}
}
