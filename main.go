package main

import (
	"fmt"

	"github.com/FogMeta/libra-os/config"
	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/module/db"
	"github.com/FogMeta/libra-os/module/log"
	"github.com/FogMeta/libra-os/module/redis"
	"github.com/FogMeta/libra-os/router"
)

func main() {
	config.Init()
	log.Init(config.Conf().Log.Level)
	db.Init()
	model.AutoMigrateDBModel()
	redis.Init()

	router.Router.SetTrustedProxies(nil)
	router.Router.Run(fmt.Sprintf(":%d", config.Conf().Port))
}
