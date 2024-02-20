package svc

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/cache"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/db"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	UserRpc  userservice.UserService
	FavorRpc favorite.Favorite
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.RedisDB.Host, c.RedisDB.Port),
		Password: c.RedisDB.Password,
		DB:       c.RedisDB.DB,
		PoolSize: c.RedisDB.PoolSize,
	})
	rdb.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	db.InitDB(database)
	cache.InitRdb(rdb)

	return &ServiceContext{
		Config:   c,
		UserRpc:  userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		FavorRpc: favorite.NewFavorite(zrpc.MustNewClient(c.FavorRpc)),
	}
}
