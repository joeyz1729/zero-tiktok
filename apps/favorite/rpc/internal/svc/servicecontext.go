package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/dao"
)

type ServiceContext struct {
	Config config.Config

	FavoriteRepository *dao.Repo
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:             c,
		FavoriteRepository: dao.NewRepo(c),
	}
}
