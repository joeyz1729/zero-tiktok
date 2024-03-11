package cache

import (
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"strconv"
	"time"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/dto"

	"github.com/go-redis/redis/v8"
	"github.com/json-iterator/go/extra"
)

const (
	UserInfoPrefix = "tiktok:user:info:"
	NormalTTL      = 24 * time.Hour
	InvalidTTL     = 5 * time.Minute
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
	rd := time.Duration(rand.Intn(5*60)) * time.Second
	_, err = RDB.Set(context.Background(), key, string(b), NormalTTL+rd).Result()
	return err
}
