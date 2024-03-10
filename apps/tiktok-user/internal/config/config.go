package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	FollowRpc zrpc.RpcClientConf

	Repo RepoConfig

	Kafka     KafkaConfig
	Snowflake struct {
		StartTime string
		MachineId uint16
	}

	Salt string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type RepoConfig struct {
	DataSource  string
	RedisAddr   string
	EsAddresses []string
	EsUsername  string
	EsPassword  string
	EsCACert    string
}
