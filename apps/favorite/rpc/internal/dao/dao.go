package dao

import (
	"context"
	"sync"

	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"
	sqlz "github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	FavoriteModel model.FavoriteModel

	delCache sync.Map
	addCache sync.Map

	FavoriteDB *sqlx.DB

	FavoriteCache *redis.Client
}

var repo *Repo

func NewRepo(c config.Config) *Repo {
	sqlConn := sqlz.NewMysql(c.Mysql.DataSource)
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)

	rdb := redis.NewClient(&redis.Options{
		Addr: c.CacheRedis.Addr,
	})
	rdb.Ping(context.Background())

	if err != nil {
		panic(err)
	}

	repo = &Repo{
		FavoriteModel: model.NewFavoriteModel(sqlConn),
		FavoriteDB:    db,
		FavoriteCache: rdb,
		addCache:      sync.Map{},
		delCache:      sync.Map{},
	}

	return repo

}

func (r *Repo) AddFavorite(uid, vid string) {
	r.addCache.Store(uid+vid, struct{}{})
}

func (r *Repo) DelFavorite(uid, vid string) {
	r.delCache.Store(uid+vid, struct{}{})
}
