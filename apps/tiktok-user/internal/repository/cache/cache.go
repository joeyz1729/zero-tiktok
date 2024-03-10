package cache

import (
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/dto"

	"github.com/go-redis/redis/v8"
	"github.com/json-iterator/go/extra"
)

const (
	UserInfoPrefix  = "tiktok:user:info:"
	UserCountPrefix = "tiktok:user:count:"
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
	ErrInvalidParams = errors.New("invalid params")
	ErrUserNotExist  = errors.New("not exist")
	ErrCacheMiss     = errors.New("cache miss")
)

func GetUser(userId int64, RDB *redis.Client) (*dto.User, error) {
	b, err := RDB.Get(context.Background(), strconv.Itoa(int(userId))).Result()
	if err != nil {
		return nil, err
	}
	var detail dto.User
	extra.RegisterFuzzyDecoders()
	err = jsoniter.Unmarshal([]byte(b), &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func AddUser(detail *dto.User, RDB *redis.Client) error {
	if detail == nil {
		return errors.New("invalid user struct")
	}
	b, err := jsoniter.Marshal(detail)
	if err != nil {
		return err
	}
	key := UserInfoPrefix + strconv.Itoa(int(detail.ID))
	_, err = RDB.Set(context.Background(), key, string(b), 24*time.Hour).Result()
	return err
}
