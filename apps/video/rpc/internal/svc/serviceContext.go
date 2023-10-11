package svc

import (
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/dao"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type ServiceContext struct {
	Config config.Config

	VideoModel model.VideoModel

	VideoCache *redis.Client

	VideoRepo dao.VideoRepo

	VideoDB *sqlx.DB
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
		//VideoModel: model.NewVideoModel(sqlConn, c.CacheRedis),
		VideoDB:    db,
		VideoCache: rdb,
		VideoRepo:  dao.NewRepoImpl(db, rdb),
	}
}
