package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	CacheRedis cache.CacheConf

	Consul consul.Conf

	Snowflake struct {
		StartTime string
		MachineId uint16
	}

	//RabbitMq struct {
	//	Username string
	//	Password string
	//	Host     string
	//	Port     int
	//}

	KafkaMq kq.KqConf
}
