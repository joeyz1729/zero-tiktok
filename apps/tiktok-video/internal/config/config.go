package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Repo RepoConfig

	MinIO struct {
		Upload bool
	}

	UserRpc  zrpc.RpcClientConf
	FavorRpc zrpc.RpcClientConf

	Snowflake struct {
		StartTime string
		MachineId uint16
	}

	Kafka KafkaConfig

	Salt string
}

type KafkaConfig struct {
	Brokers     []string
	Topic       string
	Partition   int
	MaxBytes    int
	StartOffset int64
}

type RepoConfig struct {
	DataSource  string
	RedisAddr   string
	EsAddresses []string
	EsUsername  string
	EsPassword  string
	EsCACert    string
}
