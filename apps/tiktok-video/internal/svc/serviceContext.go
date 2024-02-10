package svc

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joeyz1729/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	//VideoModel repository.VideoModel
	VideoCache *redis.Client
	VideoRepo  repository.VideoRepo
	VideoDB    *sqlx.DB

	UserRpc  userservice.UserService
	FavorRpc favorite.Favorite
}

func NewServiceContext(c config.Config) *ServiceContext {
	//sqlConn := sqlz.NewMysql(c.Mysql.DataSource)
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

	return &ServiceContext{
		Config: c,
		//VideoModel: repository.NewVideoModel(sqlConn, c.CacheRedis),
		VideoDB:    db,
		VideoCache: rdb,
		VideoRepo:  repository.NewRepoImpl(db, rdb),
		UserRpc:    userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		FavorRpc:   favorite.NewFavorite(zrpc.MustNewClient(c.FavorRpc)),
	}
}
