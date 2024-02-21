package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	FollowRpc zrpc.RpcClientConf

	Repo struct {
		DataSource  string
		RedisAddr   string
		EsAddresses []string
	}

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
