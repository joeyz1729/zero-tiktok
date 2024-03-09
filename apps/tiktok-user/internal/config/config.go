package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	FollowRpc zrpc.RpcClientConf

	Repo RepoConfig

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

type RepoConfig struct {
	DataSource  string
	RedisAddr   string
	EsAddresses []string
	EsUsername  string
	EsPassword  string
	EsCACert    string
}
