package repository

import (
	"context"
	cache2 "github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/repository/cache"
	db2 "github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/repository/db"
)

type ThumbupDao struct {
	db    *db2.GormImpl
	cache *cache2.RedisImpl
}

func NewRepo(dataSource, addr string) (*ThumbupDao, error) {
	dbImpl, err := db2.NewGormImpl(dataSource)
	if err != nil {
		return nil, err
	}
	rdbImpl, err := cache2.NewRedisImpl(addr)
	if err != nil {
		return nil, err
	}
	return &ThumbupDao{
		db:    dbImpl,
		cache: rdbImpl,
	}, nil
}

func (r *ThumbupDao) IsThumbup(c context.Context, userId, videoId int64) (bool, error) {
	ok, err := r.cache.IsThumbup(c, userId, videoId)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	// 查询数据库
	ok, err = r.db.IsThumbup(userId, videoId)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *ThumbupDao) AddThumbup(ctx context.Context, userId int64, videoId int64) error {
	err := r.db.Add(userId, videoId)
	if err != nil {
		return err
	}
	return r.cache.AddThumbup(ctx, userId, videoId)
}

func (r *ThumbupDao) DeleteThumbup(ctx context.Context, userId int64, videoId int64) error {
	err := r.db.Delete(userId, videoId)
	if err != nil {
		return err
	}
	return r.cache.DelThumbup(ctx, userId, videoId)
}

func (r *ThumbupDao) GerUserThumbupList(ctx context.Context, userId int64) ([]int64, error) {
	ok, err := r.cache.IfExist(ctx, userId)
	if err != nil {
		return nil, err
	}
	if ok {
		return r.cache.GetUserThumbupList(ctx, userId)
	}
	return r.db.GetThumbupListByUserId(userId)
}
