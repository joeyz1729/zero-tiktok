package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Repo struct {
		DataSource string
		RedisAddr  string
	}

	CacheRedis cache.CacheConf

	Snowflake struct {
		StartTime string
		MachineId uint16
	}

	Salt string
}
