package dao

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type VideoRepo interface {
	GetFavorLists(ctx context.Context, ids []int64) ([]*Video, error)
	GetVideoById(ctx context.Context, vid int64) (Video, error)
	GetVideosByUserId(ctx context.Context, uid int64) ([]Video, error)
	AddVideoInfo(ctx context.Context, video *model.Video) error
	AddVideo(ctx context.Context, video *model.Video) error
}

type RepoImpl struct {
	db    VideoDB
	cache VideoCache
}

func NewRepoImpl(db *sqlx.DB, rdb *redis.Client) *RepoImpl {
	return &RepoImpl{
		NewMysqlImpl(db),
		NewRedisImpl(rdb),
	}
}

var _ VideoRepo = (*RepoImpl)(nil)

func (r *RepoImpl) GetVideoById(ctx context.Context, vid int64) (Video, error) {
	vidStr := strconv.FormatInt(vid, 10)
	hit, err := r.cache.KeyExists(ctx, model.VideoInfoPrefix+vidStr)
	if err == nil && hit {
		// cache hit
		return r.cache.GetVideoById(ctx, model.VideoInfoPrefix+vidStr)
	}
	// cache miss
	defer func() {
		_ = r.cache.DelVideo(ctx, model.VideoInfoPrefix+vidStr)
	}()
	return r.db.GetVideoById(vid)
}

func (r *RepoImpl) GetVideosByUserId(ctx context.Context, uid int64) ([]Video, error) {
	uidStr := strconv.FormatInt(uid, 10)
	hit, err := r.cache.KeyExists(ctx, model.VideoPublishPrefix+uidStr)
	if err == nil && hit {
		// cache hit
		return r.cache.GetVideosByUser(ctx, model.VideoPublishPrefix+uidStr)
	}

	// 从数据库读取
	videos, err := r.db.GetVideoByUser(uid)
	if err != nil {
		return nil, err
	}
	// 添加用户的发布列表缓存
	// 添加视频信息的缓存
	_ = r.cache.AddPublishList(ctx, uidStr, videos)
	return videos, nil
}

func (r *RepoImpl) AddVideoInfo(ctx context.Context, video *model.Video) (err error) {
	return
}

func (r *RepoImpl) AddVideo(ctx context.Context, video *model.Video) (err error) {
	return
}

func (r *RepoImpl) GetFavorLists(ctx context.Context, ids []int64) ([]*Video, error) {
	return []*Video{}, nil
}
