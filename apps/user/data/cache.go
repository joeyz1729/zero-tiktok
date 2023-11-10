package data

import (
	"context"
	"strconv"
	"time"
)

const (
	UserInfoPrefix  = "tiktok:user:info:"  // hash
	UserCountPrefix = "tiktok:user:count:" // hash
)

const (
	FieldUsername        = "username"
	FieldAvatar          = "avatar"
	FieldBackgroundImage = "bgimg"
	FieldSignature       = "sign"
	FieldFollowedCount   = "followedcount"
	FieldFollowerCount   = "followercount"
	FieldTotalFavorited  = "totalfavorited"
	FieldWorkCount       = "workcount"
	FieldFavoriteCount   = "favoritecount"
)

var (
	ctx = context.Background()
)

func (r *Repo) AddCountCache(uid int64, user *UserInfo) (err error) {
	pipeline := r.rdb.TxPipeline()
	uidStr := strconv.FormatInt(uid, 10)
	key := UserCountPrefix + uidStr
	pipeline.HSet(ctx, key, FieldFollowerCount, user.FollowerCount)
	pipeline.HSet(ctx, key, FieldFollowedCount, user.FollowedCount)
	pipeline.HSet(ctx, key, FieldTotalFavorited, user.TotalFavorited)
	pipeline.HSet(ctx, key, FieldWorkCount, user.WorkCount)
	pipeline.HSet(ctx, key, FieldFavoriteCount, user.FavoriteCount)
	pipeline.Expire(ctx, key, 5*time.Minute)
	_, err = pipeline.Exec(ctx)
	return err
}

func (r *Repo) AddDetailCache(uid int64, user *UserInfo) (err error) {
	pipeline := r.rdb.TxPipeline()
	uidStr := strconv.FormatInt(uid, 10)
	pipeline.HSet(ctx, UserInfoPrefix+uidStr, FieldUsername, user.Username)
	pipeline.HSet(ctx, UserInfoPrefix+uidStr, FieldAvatar, user.Avatar)
	pipeline.HSet(ctx, UserInfoPrefix+uidStr, FieldBackgroundImage, user.BackgroundImage)
	pipeline.HSet(ctx, UserInfoPrefix+uidStr, FieldSignature, user.Signature)
	_, err = pipeline.Exec(ctx)
	return err
}

// DelCountCache 更新计数时，删除对应的缓存，保证一致性
func (r *Repo) DelCountCache(userId, authorId int64) (err error) {
	uidStr := strconv.FormatInt(userId, 10)
	aidStr := strconv.FormatInt(authorId, 10)
	pipeline := r.rdb.Pipeline()
	r.rdb.Del(context.Background(), UserCountPrefix+uidStr)
	r.rdb.Del(context.Background(), UserCountPrefix+aidStr)
	_, err = pipeline.Exec(context.Background())
	return
}
