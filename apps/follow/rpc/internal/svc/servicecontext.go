package svc

import "github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
