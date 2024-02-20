package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/repository/cache"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/repository/db"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type Repo struct {
	cache *cache.FollowCache
	mdb   *db.FollowDB
}

func NewRepo(mysqlDB *gorm.DB, rdb *redis.Client) *Repo {
	return &Repo{
		cache: cache.NewFollowCache(rdb),
		mdb:   db.NewFollowDB(mysqlDB),
	}
}

func (r *Repo) CheckRelation(ctx context.Context, userId int64, toUserId int64) (ok bool, err error) {
	ok, err = r.cache.GetRelation(ctx, userId, toUserId)
	if err == nil && ok {

		return true, nil
	}
	// redis没查到或出错，走数据库
	ok, err = r.mdb.IsFollow(userId, toUserId)
	if err != nil {
		logx.Error(err)
		return false, err
	}
	return ok, nil
}

func (r *Repo) AddRelation(ctx context.Context, userId int64, toUserId int64) (err error) {
	return r.mdb.Add(userId, toUserId)
}

func (r *Repo) DelRelation(ctx context.Context, userId int64, toUserId int64) (err error) {
	return r.mdb.Delete(userId, toUserId)
}

func (r *Repo) GetFollowedIds(ctx context.Context, userId int64) (ids []int64, err error) {
	//ids, err = r.cache.GetFollowerIds(ctx, userId)
	//if err == nil {
	//	return ids, nil
	//}
	//// mysql
	//ids, err = r.GetFollowerIds(userId)
	//if err != nil {
	//	return nil, err
	//}
	//// 异步添加数据库，无视错误
	//if len(ids) != 0 {
	//	go r.cache.AddFollower(ctx, userId, ids)
	//}
	return ids, nil
}

// GetFollowerIds 根据用户获取他的粉丝列表
func (r *Repo) GetFollowerIds(ctx context.Context, userId int64) (ids []int64, err error) {
	//ids, err = r.cache.GetFollowerIds(ctx, userId)
	//if err == nil {
	//	return ids, nil
	//}
	//// mysql
	//ids, err = r.mdb.GetFollowerIds(ctx, userId)
	//if err != nil {
	//	return nil, err
	//}
	//// 异步添加数据库，无视错误
	//if len(ids) != 0 {
	//	go r.cache.AddFollower(ctx, userId, ids)
	//}

	return ids, nil
}
