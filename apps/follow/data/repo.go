package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joeyz1729/zero-tiktok/apps/follow/data/cache"
	"github.com/joeyz1729/zero-tiktok/apps/follow/data/db"
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

// CheckRelation 检查是否存在关系
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
	// 问题，缓存是关注列表的集合，如果添加缓存是添加所有的吗？
	return ok, nil
}

// AddRelation 修改关系
func (r *Repo) AddRelation(userId int64, toUserId int64) (err error) {
	sqlStr := `insert into tiktok_follow.follow(user_id, follow_id) value(?, ?)`
	_, err = r.mdb.Exec(sqlStr, userId, toUserId)
	// 删除redis数据

	err = r.cache.RemRelation(ctx, userId, toUserId)
	return err
}

// DelRelation 修改关系
func (r *Repo) DelRelation(userId int64, toUserId int64) (err error) {
	sqlStr := `delete from tiktok_follow.follow where user_id = ? and follow_id = ? limit 1`
	_, err = r.mdb.Exec(sqlStr, userId, toUserId)
	err = r.cache.RemRelation(ctx, userId, toUserId)
	return err
}

// GetFollowedIds 根据用户获取他的关注列表
func (r *Repo) GetFollowedIds(userId int64) (ids []int64, err error) {
	ids, err = r.cache.GetFollowerIds(ctx, userId)
	if err == nil {
		return ids, nil
	}
	// mysql
	ids, err = r.GetFollowerIds(userId)
	if err != nil {
		return nil, err
	}
	// 异步添加数据库，无视错误
	if len(ids) != 0 {
		go r.cache.AddFollower(ctx, userId, ids)
	}
	return ids, nil
}

// GetFollowerIds 根据用户获取他的粉丝列表
func (r *Repo) GetFollowerIds(userId int64) (ids []int64, err error) {
	ids, err = r.cache.GetFollowerIds(ctx, userId)
	if err == nil {
		return ids, nil
	}
	// mysql
	ids, err = r.mdb.GetFollowerIds(ctx, userId)
	if err != nil {
		return nil, err
	}
	// 异步添加数据库，无视错误
	if len(ids) != 0 {
		go r.cache.AddFollower(ctx, userId, ids)
	}

	return ids, nil
}
