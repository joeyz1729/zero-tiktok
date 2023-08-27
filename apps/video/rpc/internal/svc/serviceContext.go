package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/jmoiron/sqlx"
	sqlz "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	VideoModel model.VideoModel

	VideoDB *sqlx.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlz.NewMysql(c.Mysql.DataSource)
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		VideoModel: model.NewVideoModel(sqlConn, c.CacheRedis),
		VideoDB:    db,
	}
}
