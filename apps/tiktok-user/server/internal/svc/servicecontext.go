package svc

import "github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
