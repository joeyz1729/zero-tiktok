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

func (r *ThumbupDao) AddThumbup(c context.Context, userId int64, videoId int64) error {
	return r.db.Add(userId, videoId)
}

func (r *ThumbupDao) DeleteThumbup(ctx context.Context, userId int64, videoId int64) error {
	return r.db.Delete(userId, videoId)
}

func (r *ThumbupDao) GerUserThumbupList(ctx context.Context, userId int64) ([]int64, error) {
	return r.db.GetThumbupListByUserId(userId)
}

// AddCacheFromDB 从数据库中读取点赞列表，并添加到缓存
func (r *ThumbupDao) AddCacheFromDB(userId int64) error {
	//// 数据库读取
	//ids, err := r.db.GetThumbupListByUserId(userId)
	//if err != nil {
	//	return err
	//}
	//// 添加到数据
	//key := cache2.FavoriteSetPrefix + strconv.FormatInt(userId, 10)
	//return r.cache.AddFavorites(context.Background(), key, ids...)
	return nil
}

func (r *ThumbupDao) GetFavoriteIdsFromDB(userId int64) ([]int64, error) {
	//favors := []*pb.Favorite{}
	//err := r.db.Table("favorite").Select("video_id").Where("user_id = ?", userId).Find(&favors).Error
	//if err != nil {
	//	return nil, err
	//}
	//ids := make([]int64, len(favors))
	//for i, f := range favors {
	//	ids[i] = f.VideoId
	//}
	//return ids, nil
	return nil, nil
}
