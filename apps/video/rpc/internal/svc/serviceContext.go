package svc

import (
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/data"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/config"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	//VideoModel data.VideoModel
	VideoCache *redis.Client
	VideoRepo  data.VideoRepo
	VideoDB    *sqlx.DB

	UserRpc  user.User
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
		//VideoModel: data.NewVideoModel(sqlConn, c.CacheRedis),
		VideoDB:    db,
		VideoCache: rdb,
		VideoRepo:  data.NewRepoImpl(db, rdb),
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		FavorRpc:   favorite.NewFavorite(zrpc.MustNewClient(c.FavorRpc)),
	}
}
