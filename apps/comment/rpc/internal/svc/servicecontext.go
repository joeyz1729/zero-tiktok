package svc

import "github.com/YiZou89/zero-tiktok/apps/comment/rpc/comment/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
