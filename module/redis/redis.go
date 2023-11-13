package redis

import (
	"context"
	"fmt"

	"github.com/FogMeta/libra-os/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Init() {
	rc := config.Conf().Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rc.Host, rc.Port),
		Password: rc.Password,
		DB:       rc.DB,
	})
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
