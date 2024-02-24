package svc

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/job"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	FollowRepo *repository.Repo

	Worker *job.Worker

	UserRpc userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
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

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.Kafka.Brokers,
		Topic:       c.Kafka.Topic,
		Partition:   0,
		MaxBytes:    10e6, // 10MB
		GroupID:     "tiktok-user",
		StartOffset: kafka.LastOffset,
	})
	repo := repository.NewRepo(db, rdb)

	return &ServiceContext{
		Config:     c,
		FollowRepo: repo,
		Worker: &job.Worker{
			Repo:        repo,
			KafkaReader: reader,
		},
		UserRpc: userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
	}
}
