package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Repo struct {
		DataSource  string
		RedisAddr   string
		EsAddresses []string
	}

	MinIO struct {
		Upload bool
	}

	UserRpc  zrpc.RpcClientConf
	FavorRpc zrpc.RpcClientConf

	Snowflake struct {
		StartTime string
		MachineId uint16
	}

	Kafka struct {
		Brokers []string
		Topic   string
	}

	Salt string
}
