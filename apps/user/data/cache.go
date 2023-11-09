package data

import (
	"context"
	"strconv"
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
	pipeline.HSet(ctx, UserCountPrefix+uidStr, FieldFollowerCount, user.FollowerCount)
	pipeline.HSet(ctx, UserCountPrefix+uidStr, FieldFollowedCount, user.FollowedCount)
	pipeline.HSet(ctx, UserCountPrefix+uidStr, FieldTotalFavorited, user.TotalFavorited)
	pipeline.HSet(ctx, UserCountPrefix+uidStr, FieldWorkCount, user.WorkCount)
	pipeline.HSet(ctx, UserCountPrefix+uidStr, FieldFavoriteCount, user.FavoriteCount)
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
