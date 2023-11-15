package model

import "github.com/FogMeta/libra-os/module/db"

type Table interface {
	TableName() string
}

func AutoMigrateDBModel() {
	if err := db.DB.AutoMigrate(new(User)); err != nil {
		panic(err)
	}
}
