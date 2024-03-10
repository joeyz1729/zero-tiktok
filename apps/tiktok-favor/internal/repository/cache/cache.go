package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	ThumbupPrefix = "tiktok:thumbup:"

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

func (c *RedisImpl) AddThumbup(ctx context.Context, uid, vid int64) (err error) {
	uidStr, vidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(vid, 10)
	_, err = c.Client.SAdd(ctx, ThumbupPrefix+uidStr, vidStr).Result()
	return err
}

func (c *RedisImpl) DelThumbup(ctx context.Context, uid, vid int64) (err error) {
	uidStr, vidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(vid, 10)
	_, err = c.Client.SRem(ctx, ThumbupPrefix+uidStr, vidStr).Result()
	return err
}

func (c *RedisImpl) IsThumbup(ctx context.Context, userId int64, videoId int64) (bool, error) {
	uidStr, vidStr := strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10)
	ok, err := c.Client.SIsMember(ctx, ThumbupPrefix+uidStr, vidStr).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (c *RedisImpl) GetUserThumbupList(ctx context.Context, userId int64) ([]int64, error) {
	uidStr := strconv.FormatInt(userId, 10)
	result, err := c.Client.SMembers(ctx, ThumbupPrefix+uidStr).Result()
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

func (c *RedisImpl) IfExist(ctx context.Context, userId int64) (bool, error) {
	uidStr := strconv.FormatInt(userId, 10)
	count, err := c.Client.Exists(ctx, ThumbupPrefix+uidStr).Result()
	if err != nil {
		return false, err
	}
	return count == 1, nil
}
