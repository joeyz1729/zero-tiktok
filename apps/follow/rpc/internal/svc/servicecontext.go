package svc

import (
	"context"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/follow/dao/cache"
	datadb "github.com/YiZou89/zero-tiktok/apps/follow/dao/db"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/bloom"

	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/config"
	"github.com/go-redis/redis/v8"
	sqlx "github.com/jmoiron/sqlx"
	redisz "github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config

	FollowDB *datadb.FollowDB

	FollowCache *cache.FollowCache

	KqPusher *kq.Pusher
	//KqWriter *kafka.Writer

	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	redisAddr := fmt.Sprintf("%s:%d", c.RedisDB.Host, c.RedisDB.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: c.RedisDB.Password,
		DB:       c.RedisDB.DB,
		PoolSize: c.RedisDB.PoolSize,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	store := redisz.MustNewRedis(
		redisz.RedisConf{
			Host: redisAddr,
			Type: redisz.NodeType},
	)
	filter := bloom.New(store, "tiktok:follow:bloom", 20*1<<20)

	pusher := kq.NewPusher(c.KafkaMq.Brokers, c.KafkaMq.Topic)

	return &ServiceContext{
		Config:      c,
		FollowDB:    db.NewFollowDB(db),
		FollowCache: cache.NewFollowCache(rdb),
		Filter:      filter,
		//KqWriter:    w,
		KqPusher: pusher,
	}
}
