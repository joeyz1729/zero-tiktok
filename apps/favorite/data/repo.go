package data

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"
	"strconv"
)

type RepoImpl struct {
	db    FavorDB
	cache FavorCache
}

func NewRepo(dataSource, addr string) (*RepoImpl, error) {
	dbImpl, err := NewDBImpl(dataSource)
	if err != nil {
		return nil, err
	}
	rdbImpl, err := NewCacheImpl(addr)
	if err != nil {
		return nil, err
	}
	return &RepoImpl{
		db:    dbImpl,
		cache: rdbImpl,
	}, nil
}

type FavorRepo interface {
	CheckFavor(c context.Context, userId, videoId int64) (bool, error)
	CreateFavoriteRecord(c context.Context, favor *model.Favorite) error
	DeleteFavoriteRecord(c context.Context, favor *model.Favorite) error

	GetFavorIds(c context.Context, userId int64) ([]int64, error)
}

var _ FavorRepo = (*RepoImpl)(nil)

// CheckFavor 查询用户点赞关系，并更新缓存
func (r *RepoImpl) CheckFavor(c context.Context, userId, videoId int64) (bool, error) {
	// 查询缓存
	uidStr := strconv.FormatInt(userId, 10)
	ok, err := r.cache.IsFavRecordExist(c, FavoriteSetPrefix+uidStr, videoId)
	if err == nil && ok {
		// cache hit
		return ok, nil
	}
	// 查询数据库
	ok, err = r.db.IsFavoriteRecordExist(userId, videoId)
	if err != nil {
		return false, err
	}
	go func() { _ = r.AddCacheFromDB(userId) }()
	return ok, nil
}

func (r *RepoImpl) CreateFavoriteRecord(c context.Context, favor *model.Favorite) error {
	// 添加数据库
	err := r.db.CreateFavoriteRecord(favor)
	if err != nil {
		return err
	}
	uidStr := strconv.FormatInt(favor.UserId, 10)
	// 删除缓存
	return r.cache.DelFavorite(c, FavoriteSetPrefix+uidStr, favor.VideoId)
}

func (r *RepoImpl) DeleteFavoriteRecord(c context.Context, favor *model.Favorite) error {
	// 删除数据库
	err := r.db.DeleteFavoriteRecord(favor)
	if err != nil {
		return err
	}
	uidStr := strconv.FormatInt(favor.UserId, 10)
	// 删除缓存
	return r.cache.DelFavorite(c, FavoriteSetPrefix+uidStr, favor.VideoId)
}

func (r *RepoImpl) GetFavorIds(c context.Context, userId int64) ([]int64, error) {
	//	check cache
	key := FavoriteSetPrefix + strconv.FormatInt(userId, 10)
	hit, err := r.cache.KeyExist(c, key)
	if err == nil && hit {
		// cache
		ids, err := r.cache.GetFavoriteVideoIds(c, key)
		if err == nil {
			return ids, nil
		}
	}
	// database
	ids, err := r.db.GetFavoriteIds(userId)
	if err != nil {
		return nil, err
	}

	// add cache
	err = r.cache.AddFavorites(c, key, ids...)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// AddCacheFromDB 从数据库中读取点赞列表，并添加到缓存
func (r *RepoImpl) AddCacheFromDB(userId int64) error {
	// 数据库读取
	ids, err := r.db.GetFavoriteIds(userId)
	if err != nil {
		return err
	}
	// 添加到数据
	key := FavoriteSetPrefix + strconv.FormatInt(userId, 10)
	return r.cache.AddFavorites(context.Background(), key, ids...)
}
