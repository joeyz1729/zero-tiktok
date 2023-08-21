package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRpc user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
