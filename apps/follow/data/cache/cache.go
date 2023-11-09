package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

var (
	RelationPrefix = "tiktok:relation:"
	BloomPrefix    = RelationPrefix + "bloom:"
	FollowedPrefix = RelationPrefix + "followed:" // + uid, (tid) set
	FollowerPrefix = RelationPrefix + "follower:" // + uid, (tid) set
	CountPrefix    = "tiktok:relation:cnt:"       // + uid, (follower, followed) hash
	FollowerField  = "follower"
	FollowedField  = "followed"
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

func (fc *FollowCache) RemRelation(ctx context.Context, uid, tid int64) (err error) {
	uidStr, tidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(tid, 10)
	pipeline := fc.Client.TxPipeline()
	_, err = pipeline.SRem(ctx, FollowedPrefix+uidStr, tid).Result()
	_, err = pipeline.SRem(ctx, FollowerPrefix+tidStr, uid).Result()
	_, err = pipeline.HDel(ctx, CountPrefix+uidStr, "*").Result()
	_, err = pipeline.HDel(ctx, CountPrefix+tidStr, "*").Result()
	_, err = pipeline.Exec(ctx)
	if err != nil {
		logx.Errorw("[redis] transaction pipeline failed",
			logx.Field("err", err),
		)
		return err
	}
	logx.Info("[redis] del relation success")
	return nil
}

func (fc *FollowCache) UpdateAll(ctx context.Context, uid, tid int64) (err error) {
	uidStr, tidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(tid, 10)
	pipeline := fc.Client.TxPipeline()
	_, err = pipeline.SRem(ctx, FollowedPrefix+uidStr, tid).Result()
	_, err = pipeline.SRem(ctx, FollowerPrefix+tidStr, uid).Result()
	_, err = pipeline.HDel(ctx, CountPrefix+uidStr).Result()
	_, err = pipeline.HDel(ctx, CountPrefix+tidStr).Result()
	_, err = pipeline.Exec(ctx)
	if err != nil {
		logx.Errorw("[redis] transaction pipeline failed",
			logx.Field("err", err),
		)
		return err
	}
	logx.Info("[redis] del relation success")
	return nil
}

func (fc *FollowCache) AddRelation(ctx context.Context, uid, tid int64, cnt bool, fedCount, ferCount int32) (err error) {
	ctx = context.Background()
	uidStr, tidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(tid, 10)
	pipeline := fc.Client.TxPipeline()
	_, err = pipeline.SAdd(ctx, FollowedPrefix+uidStr, tid).Result()
	_, err = pipeline.SAdd(ctx, FollowerPrefix+tidStr, uid).Result()
	//_, err = pipeline.ExpireNX(ctx, FollowedPrefix+uidStr, time.Minute*5).Result()
	//_, err = pipeline.ExpireNX(ctx, FollowerPrefix+tidStr, time.Minute*5).Result()
	if cnt {
		_, err = pipeline.HSet(ctx, CountPrefix+uidStr, FollowedField, fedCount).Result()
		_, err = pipeline.HSet(ctx, CountPrefix+tidStr, FollowerField, ferCount).Result()
		//_, err = pipeline.ExpireNX(ctx, CountPrefix+uidStr, time.Minute*5).Result()
		//_, err = pipeline.ExpireNX(ctx, CountPrefix+tidStr, time.Minute*5).Result()
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		logx.Errorw("[redis] transaction pipeline failed",
			logx.Field("err", err),
		)
		return err
	}
	logx.Info("[redis] del relation success")
	return nil
}

func (fc *FollowCache) AddFollow(ctx context.Context, uid, tid int64) (err error) {
	ctx = context.Background()
	uidStr, tidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(tid, 10)
	pipeline := fc.Client.TxPipeline()
	_, err = pipeline.SAdd(ctx, FollowedPrefix+uidStr, tid).Result()
	_, err = pipeline.SAdd(ctx, FollowerPrefix+tidStr, uid).Result()
	_, err = pipeline.Exec(ctx)
	if err != nil {
		logx.Errorw("[redis] transaction pipeline failed",
			logx.Field("err", err),
		)
		return err
	}
	logx.Info("[redis] del relation success")
	return nil
}

func (fc *FollowCache) UpdateCount(ctx context.Context, uid, tid int64) (err error) {
	uidStr, tidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(tid, 10)
	pipeline := fc.Client.TxPipeline()
	_, err = pipeline.HDel(ctx, CountPrefix+uidStr).Result()
	_, err = pipeline.HDel(ctx, CountPrefix+tidStr).Result()
	_, err = pipeline.Exec(ctx)
	if err != nil {
		logx.Errorw("[redis] transaction pipeline failed",
			logx.Field("err", err),
		)
		return err
	}
	logx.Info("[redis] del relation success")
	return nil
}

func (fc *FollowCache) GetRelation(ctx context.Context, uid, tid int64) (ok bool, err error) {
	uidStr := strconv.FormatInt(uid, 10)
	redisKey := FollowedPrefix + uidStr
	ok, err = fc.SIsMember(ctx, redisKey, tid).Result()
	if err != nil {
		return false, ErrExecFailed
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (fc *FollowCache) GetCount(ctx context.Context, uid int64) (follower, followed int32, err error) {
	ctx = context.Background()
	uidStr := strconv.FormatInt(uid, 10)
	key := CountPrefix + uidStr
	result, err := fc.HMGet(ctx, key, FollowedField, FollowerField).Result()
	if err != nil {
		return 0, 0, ErrExecFailed
	}
	if result == nil || len(result) != 2 {
		return 0, 0, ErrMiss
	}
	return result[0].(int32), result[1].(int32), nil
}

func (fc *FollowCache) DelRelation(ctx context.Context, uid, tid int64) (err error) {
	ctx = context.Background()
	uidStr, tidStr := strconv.FormatInt(uid, 10), strconv.FormatInt(tid, 10)
	pipeline := fc.Client.TxPipeline()
	_, err = pipeline.Del(ctx, FollowedPrefix+uidStr, FollowerPrefix+tidStr).Result()
	_, err = pipeline.Del(ctx, CountPrefix+uidStr, CountPrefix+tidStr).Result()
	_, err = pipeline.Exec(ctx)
	if err != nil {
		logx.Errorw("[redis] del relation transaction pipeline failed",
			logx.Field("err", err),
		)
		return err
	}
	logx.Info("[redis] del relation success")
	return nil
}
