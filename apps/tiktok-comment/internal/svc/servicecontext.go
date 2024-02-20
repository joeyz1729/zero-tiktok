package svc

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config config.Config

	Repo *repository.Repo

	KafkaPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
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
		Config:      c,
		Repo:        repository.NewRepo(db, rdb),
		KafkaPusher: kq.NewPusher(c.KafkaMq.Brokers, c.KafkaMq.Topic),
	}
}
