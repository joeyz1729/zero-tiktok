package data

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/data/cache"
	"github.com/YiZou89/zero-tiktok/apps/follow/data/db"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ctx = context.Background()
)

type repo interface{}

type Repo struct {
	cache *cache.FollowCache
	mdb   *db.FollowDB
}

func NewRepo(mysqlDB *sqlx.DB, rdb *redis.Client) *Repo {
	return &Repo{
		cache: cache.NewFollowCache(rdb),
		mdb:   db.NewFollowDB(mysqlDB),
	}
}

func (r *Repo) CheckRelation(userId int64, toUserId int64) (ok bool, err error) {
	ok, err = r.cache.GetRelation(ctx, userId, toUserId)
	if err == nil && ok {

		return true, nil
	}
	// redis没查到或出错，走数据库
	ok, err = r.mdb.CheckRelation(ctx, userId, toUserId)
	if err != nil {
		logx.Error(err)
		return false, err
	}
	return ok, nil
}

func (r *Repo) AddRelation(userId int64, toUserId int64) (err error) {
	sqlStr := `insert into tiktok_follow.follow(user_id, follow_id) value(?, ?)`
	_, err = r.mdb.Exec(sqlStr, userId, toUserId)
	// 删除redis数据

	err = r.cache.RemRelation(ctx, userId, toUserId)
	return err
}

func (r *Repo) DelRelation(userId int64, toUserId int64) (err error) {
	sqlStr := `delete from tiktok_follow.follow where user_id = ? and follow_id = ? limit 1`
	_, err = r.mdb.Exec(sqlStr, userId, toUserId)
	err = r.cache.RemRelation(ctx, userId, toUserId)
	return err
}
