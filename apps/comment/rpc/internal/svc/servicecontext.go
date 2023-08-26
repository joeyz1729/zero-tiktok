package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/model"
	"github.com/jmoiron/sqlx"
	sqlz "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	CommentModel model.CommentModel

	CommentDB *sqlx.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlz.NewMysql(c.Mysql.DataSource)
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:       c,
		CommentModel: model.NewCommentModel(sqlConn),
		CommentDB:    db,
	}
}
