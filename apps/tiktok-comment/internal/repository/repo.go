package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/repository/cache"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/repository/db"
	"gorm.io/gorm"
)

type Repo struct {
	cache *cache.FollowCache
	mdb   *db.CommentDB
}

func NewRepo(mysqlDB *gorm.DB, rdb *redis.Client) *Repo {
	return &Repo{
		cache: cache.NewFollowCache(rdb),
		mdb:   db.NewCommentDB(mysqlDB),
	}
}

func (r *Repo) AddComment(ctx context.Context, comment *db.Comment) error {
	return r.mdb.Add(comment)
}

func (r *Repo) DelComment(ctx context.Context, commentId int64) error {
	return r.mdb.Delete(commentId)
}
