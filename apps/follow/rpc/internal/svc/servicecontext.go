package svc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joeyz1729/zero-tiktok/apps/follow/data"
	"github.com/joeyz1729/zero-tiktok/apps/follow/data/cache"
	mdb "github.com/joeyz1729/zero-tiktok/apps/follow/data/db"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/user"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	FollowRepo *data.Repo

	FollowDB *mdb.FollowDB

	FollowCache *cache.FollowCache

	UserRpc user.User

	BarrierDB *sql.DB
	DtmServer string
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	barrier, err := sql.Open("mysql", "root:root1234@tcp(localhost:13306)/tiktok_follow?parseTime=true&charset=utf8")
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
		Config:      c,
		FollowRepo:  data.NewRepo(db, rdb),
		UserRpc:     user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		FollowCache: cache.NewFollowCache(rdb),
		FollowDB:    mdb.NewFollowDB(db),
		DtmServer:   c.DtmServer,
		BarrierDB:   barrier,
	}
}
