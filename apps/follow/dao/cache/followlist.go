package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func (fc *FollowCache) GetFollowedIds(ctx context.Context, userId int64) ([]int64, error) {
	ctx = context.Background()
	uidStr := strconv.FormatInt(userId, 10)
	idStr, err := fc.SMembers(ctx, FollowedPrefix+uidStr).Result()
	if err != nil {
		return nil, err
	}
	if len(idStr) == 0 {
		return nil, ErrMiss
	}

	ids := make([]int64, len(idStr))
	for i, str := range idStr {
		id, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			logx.Errorw("parse int err",
				logx.Field("err", err))
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func (fc *FollowCache) GetFollowerIds(ctx context.Context, userId int64) ([]int64, error) {
	ctx = context.Background()
	uidStr := strconv.FormatInt(userId, 10)
	idStr, err := fc.SMembers(ctx, FollowerPrefix+uidStr).Result()
	if err != nil {
		return nil, err
	}
	if len(idStr) == 0 {
		return nil, ErrMiss
	}

	ids := make([]int64, len(idStr))
	for i, str := range idStr {
		id, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			logx.Errorw("parse int err",
				logx.Field("err", err))
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func (fc *FollowCache) AddFollower(ctx context.Context, userId int64, ids []int64) (err error) {
	ctx = context.Background()
	uidStr := strconv.FormatInt(userId, 10)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	pipeline := fc.TxPipeline()
	pipeline.Del(ctx, FollowerPrefix+uidStr)
	pipeline.SAdd(ctx, FollowerPrefix+uidStr, args...)
	pipeline.Expire(ctx, FollowerPrefix+uidStr, 5*time.Minute)
	_, err = pipeline.Exec(ctx)
	if err != nil {
		logx.Error("redis pipeline failed", err)
		return err
	}
	return nil
}

func (fc *FollowCache) AddFollowed(ctx context.Context, userId int64, ids []int64) (err error) {
	ctx = context.Background()
	uidStr := strconv.FormatInt(userId, 10)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	pipeline := fc.TxPipeline()
	pipeline.Del(ctx, FollowedPrefix+uidStr)
	pipeline.SAdd(ctx, FollowedPrefix+uidStr, args...)
	pipeline.Expire(ctx, FollowedPrefix+uidStr, time.Minute*5)
	_, err = pipeline.Exec(ctx)
	if err != nil {
		logx.Error("redis pipeline failed", err)
		return err
	}
	return nil
}
