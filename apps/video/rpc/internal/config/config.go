package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	CacheRedis cache.CacheConf

	UserRpc zrpc.RpcClientConf

	Snowflake struct {
		StartTime string
		MachineId uint16
	}

	RedisDB struct {
		Host     string
		Port     int
		Password string
		DB       int
		PoolSize int
	}

	Salt string
}
