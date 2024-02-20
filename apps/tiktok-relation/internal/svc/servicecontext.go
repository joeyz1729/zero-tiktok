package svc

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	FollowRepo *repository.Repo

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

	return &ServiceContext{
		Config:     c,
		FollowRepo: repository.NewRepo(db, rdb),
		UserRpc:    userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
	}
}
