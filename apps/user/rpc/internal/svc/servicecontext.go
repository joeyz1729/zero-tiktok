package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel model.UserModel

	// TODO, 添加redis

}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlConn),
	}
}
