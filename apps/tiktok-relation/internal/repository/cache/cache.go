package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var (
	FollowedPrefix = "tiktok:relation:followed:" // + uid, (tid) set
	FollowerPrefix = "tiktok:relation:follower:" // + uid, (tid) set
	//BloomPrefix    = "bloom:"

)

var (
	ErrMiss       = errors.New("cache miss")
	ErrExecFailed = errors.New("cache exec failed")
)

type FollowCache struct {
	*redis.Client
}

func NewFollowCache(client *redis.Client) *FollowCache {
	return &FollowCache{
		client,
	}

}

func (fc *FollowCache) AddRelation(ctx context.Context, uid, tid int64) (err error) {
	uidStr, tidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(tid, 10)
	pipeline := fc.Client.Pipeline()
	_, err = pipeline.SAdd(ctx, FollowedPrefix+uidStr, tid).Result()
	_, err = pipeline.SAdd(ctx, FollowerPrefix+tidStr, uid).Result()
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (fc *FollowCache) DelRelation(ctx context.Context, uid, tid int64) (err error) {
	uidStr, tidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(tid, 10)
	pipeline := fc.Client.Pipeline()
	_, err = pipeline.SRem(ctx, FollowedPrefix+uidStr, tid).Result()
	_, err = pipeline.SRem(ctx, FollowerPrefix+tidStr, uid).Result()
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (fc *FollowCache) GetRelation(ctx context.Context, uid, tid int64) (ok bool, err error) {
	uidStr := strconv.FormatInt(uid, 10)
	ok, err = fc.SIsMember(ctx, FollowedPrefix+uidStr, tid).Result()
	if err != nil {
		return false, ErrExecFailed
	}
	return ok, nil
}
