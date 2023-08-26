package svc

import "github.com/YiZou89/zero-tiktok/apps/favorite/favorite/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
