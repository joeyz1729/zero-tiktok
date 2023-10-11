package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

const (
	FavoriteSetPrefix = "tiktok:favorite:" // favorite:userId videoId, set
)

type CacheImpl struct {
	*redis.Client
}

func NewCacheImpl(addr string) (*CacheImpl, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	return &CacheImpl{rdb}, nil
}

type FavorCache interface {
	GetKeys(c context.Context, keyPattern string) ([]string, error)
	KeyExist(c context.Context, key string) (bool, error)
	CreateFavorite(c context.Context, key string, video int64) error
	DelFavorite(c context.Context, key string, videoId int64) error
	IsFavRecordExist(c context.Context, key string, videoId int64) (bool, error)
	GetFavoriteVideoIds(c context.Context, key string) ([]int64, error)
}

var _ FavorCache = (*CacheImpl)(nil)

func (c *CacheImpl) IsFavRecordExist(ctx context.Context, key string, videoId int64) (bool, error) {
	return c.SIsMember(ctx, key, videoId).Result()
}

func (c *CacheImpl) CreateFavorite(ctx context.Context, key string, vid int64) (err error) {
	_, err = c.SAdd(ctx, key, vid).Result()
	return err
}

func (c *CacheImpl) DelFavorite(ctx context.Context, key string, videoId int64) (err error) {
	_, err = c.SRem(ctx, key, videoId).Result()

	return err
}
func (c *CacheImpl) KeyExist(ctx context.Context, key string) (bool, error) {
	num, err := c.Exists(ctx, key).Result()
	return num == 0, err
}

func (c *CacheImpl) GetKeys(ctx context.Context, keyPattern string) ([]string, error) {
	return c.GetKeys(ctx, keyPattern)
}

func (c *CacheImpl) GetFavoriteVideoIds(ctx context.Context, key string) ([]int64, error) {
	result, err := c.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(result))
	for i, idStr := range result {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}
