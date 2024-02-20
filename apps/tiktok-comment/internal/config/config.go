package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	CacheRedis struct {
		Addr string
	}

	Snowflake struct {
		StartTime string
		MachineId uint16
	}

	KafkaMq kq.KqConf
}
