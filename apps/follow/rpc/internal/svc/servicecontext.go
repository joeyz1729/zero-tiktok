package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	sqlz "github.com/zeromicro/go-zero/core/stores/sqlx"

	sqlx "github.com/jmoiron/sqlx"
)

type ServiceContext struct {
	Config config.Config

	FollowModel model.FollowModel

	FollowDB *sqlx.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlz.NewMysql(c.Mysql.DataSource)
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:      c,
		FollowModel: model.NewFollowModel(sqlConn, c.CacheRedis),
		FollowDB:    db,
	}
}
