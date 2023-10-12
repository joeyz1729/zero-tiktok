package dao

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"
	"strconv"
)

type RepoImpl struct {
	db    FavorDB
	cache FavorCache
}

func NewRepo(c config.Config) (*RepoImpl, error) {
	dbImpl, err := NewDBImpl(c.Mysql.DataSource)
	if err != nil {
		return nil, err
	}
	rdbImpl, err := NewCacheImpl(c.CacheRedis.Addr)
	if err != nil {
		return nil, err
	}
	return &RepoImpl{
		db:    dbImpl,
		cache: rdbImpl,
	}, nil
}

type FavorRepo interface {
	IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error)
	IsFavoriteRecordExist(c context.Context, userId, videoId int64) (bool, error)
	CreateFavoriteRecord(c context.Context, favor *model.Favorite) error
	DeleteFavoriteRecord(c context.Context, favor *model.Favorite) error

	GetFavorIds(c context.Context, userId int64) ([]int64, error)
}

var _ FavorRepo = (*RepoImpl)(nil)

// IsFavoriteRecordExist 更新前先检查是否存在，查询数据库时不更新缓存
func (r *RepoImpl) IsFavoriteRecordExist(c context.Context, userId, videoId int64) (bool, error) {
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
	return ok, nil
}

// IsFavoriteVideo 查询关系，并更新缓存
func (r *RepoImpl) IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error) {
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
	// 如果存在就添加缓存
	if ok {
		err = r.cache.CreateFavorite(c, FavoriteSetPrefix+uidStr, videoId)
	}
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
