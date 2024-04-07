package entity

import (
	"github.com/qiushenglei/gin-skeleton/pkg/xdb"
)

type Configure struct {
	Env    string       `json:"env"`
	Domain string       `json:"domain"`
	Source int          `json:"source"`
	Port   int          `json:"port"`
	DB1    *xdb.RDB     `json:"db1"`
	RD1    *xdb.RedisDB `json:"rd1"`
}

type MysqlConf struct {
	DSN string `json:"dsn"`
}

type RedisConf struct {
	DSN string `json:"dsn"`
}
