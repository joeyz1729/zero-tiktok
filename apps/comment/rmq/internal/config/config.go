package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Kafka    kq.KqConf
	UserRpc  zrpc.RpcClientConf
	VideoRpc zrpc.RpcClientConf
}
