package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	RedisDB struct {
		Host     string
		Port     int
		DB       int
		PoolSize int
		Password string
	}

	Kafka struct {
		Brokers []string
		Topic   string
	}

	UserRpc zrpc.RpcClientConf
}
