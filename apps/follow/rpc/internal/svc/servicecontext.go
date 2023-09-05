package svc

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/bloom"

	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"github.com/go-redis/redis/v8"
	sqlx "github.com/jmoiron/sqlx"
	redisz "github.com/zeromicro/go-zero/core/stores/redis"
	sqlz "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	FollowModel model.FollowModel

	FollowDB *sqlx.DB

	FollowCache *redis.Client

	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlz.NewMysql(c.Mysql.DataSource)
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.RedisDB.Host, c.RedisDB.Port),
		Password: c.RedisDB.Password,
		DB:       c.RedisDB.DB,
		PoolSize: c.RedisDB.PoolSize,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	store := redisz.New(c.CacheRedis[0].Host, func(r *redisz.Redis) {
		r.Type = redisz.NodeType
	})
	filter := bloom.New(store, "tiktok:follow:bloom", 20*1<<20)

	return &ServiceContext{
		Config:      c,
		FollowModel: model.NewFollowModel(sqlConn, c.CacheRedis),
		FollowDB:    db,
		FollowCache: rdb,
		Filter:      filter,
	}
}
