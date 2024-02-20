package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	ThumbupFormat = "tiktok:thumbup:%d:%d" // favorite:userId videoId, set

	NormalTTL = 24 * time.Hour
	NilTTL    = 5 * time.Minute
)

type RedisImpl struct {
	*redis.Client
}

func NewRedisImpl(addr string) (*RedisImpl, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	return &RedisImpl{rdb}, nil
}

func (c *RedisImpl) IsThumbup(ctx context.Context, userId int64, videoId int64) (bool, error) {
	count, err := c.Exists(ctx, fmt.Sprintf(ThumbupFormat, userId, videoId)).Result()
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (c *RedisImpl) AddThumbup(ctx context.Context, userId int64, videoId int64) (err error) {
	_, err = c.Set(ctx, fmt.Sprintf(ThumbupFormat, userId, videoId), 1, NormalTTL).Result()
	return err
}

func (c *RedisImpl) DelThumbup(ctx context.Context, userId int64, videoId int64) (err error) {
	_, err = c.Del(ctx, fmt.Sprintf(ThumbupFormat, userId, videoId)).Result()
	return err
}

func (c *RedisImpl) GetKeys(ctx context.Context, keyPattern string) ([]string, error) {
	return c.GetKeys(ctx, keyPattern)
}

func (c *RedisImpl) IsExist(ctx context.Context, userId int64) (bool, error) {
	count, err := c.Exists(ctx, fmt.Sprintf(ThumbupFormat, userId, 0)).Result()
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (c *RedisImpl) GetUserThumbupList(ctx context.Context, userId int64) ([]int64, error) {
	result, err := c.SMembers(ctx, ThumbupFormat+strconv.FormatInt(userId, 10)).Result()
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
