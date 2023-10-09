package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/config"
	model2 "github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel model2.UserModel

	UserRepo *model2.Repo
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Repo.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model2.NewUserModel(sqlConn),
		UserRepo:  model2.NewRepo(c.Repo.DataSource, c.Repo.RedisAddr),
	}
}
