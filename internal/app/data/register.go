package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"vod/internal/app/Init"
)

var VODRedis *redis.Client
var VODMysql *gorm.DB

func InitData(c context.Context) {
	var err error

	// mysql
	//if VODMysql, err = Init.Configure.DB1.RegisterMySQL(); err != nil {
	//	panic(err)
	//}

	// redis
	if VODRedis, err = Init.Configure.RD1.RegisterRDBClient(c); err != nil {
		panic(err)
	}

}
