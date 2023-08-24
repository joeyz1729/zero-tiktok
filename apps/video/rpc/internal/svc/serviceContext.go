package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	VideoModel model.VideoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		VideoModel: model.NewVideoModel(sqlConn),
	}
}
