package svc

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/config"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-queue/kq"

	//"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/kmq"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/model"
	"github.com/jmoiron/sqlx"
	sqlz "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	CommentModel model.CommentModel

	CommentDB *sqlx.DB

	CommentCache *redis.Client

	KafkaPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlz.NewMysql(c.Mysql.DataSource)
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: c.CacheRedis.Addr,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:       c,
		CommentModel: model.NewCommentModel(sqlConn),
		CommentDB:    db,
		CommentCache: rdb,
		KafkaPusher:  kq.NewPusher(c.KafkaMq.Brokers, c.KafkaMq.Topic),
	}
}
