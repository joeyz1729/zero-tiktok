package model

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

const (
	UserInfoPrefix  = "tiktok:user:info:"  // hash
	UserCountPrefix = "tiktok:user:count:" // hash
)

const (
	FieldUsername       = "username"
	FieldFollowedCount  = "followedcount"
	FieldFollowerCount  = "followercount"
	FieldTotalFavorited = "totalfavorited"
	FieldWorkCount      = "workcount"
	FieldFavoriteCount  = "favoritecount"
)

const (
	defaultAvatar          = "default avatar"
	defaultBackgroundImage = "default background image"
	defaultSignature       = "default signature"
)

var (
	ErrUserNotExist    = errors.New("user not exist")
	ErrCacheMiss       = errors.New("cache miss")
	ErrInvalidPassword = errors.New("invalid password")
)

var (
	countMap = map[string]interface{}{
		FieldFollowedCount:  0,
		FieldFollowerCount:  0,
		FieldTotalFavorited: 0,
		FieldWorkCount:      0,
		FieldFavoriteCount:  0,
	}
)

type Repo struct {
	db  *sqlx.DB
	rdb *redis.Client
}

var repo *Repo

func NewRepo(datasource, redisAddr string) *Repo {
	db, err := sqlx.Connect("mysql", datasource)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	rdb.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	repo = &Repo{
		db:  db,
		rdb: rdb,
	}
	return repo
}
