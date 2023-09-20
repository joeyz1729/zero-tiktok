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

	//Consul consul.Conf
	//Auth struct { // JWT 认证需要的密钥和过期时间配置
	//	AccessSecret string
	//	AccessExpire int64
	//}

	Snowflake struct {
		StartTime string
		MachineId uint16
	}

	Salt string
}
