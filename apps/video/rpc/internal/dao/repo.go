package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"sync"
)

type VideoRepo interface {
	GetVideoById(ctx context.Context, vid int64) (*Video, error)

	GetVideosByAuthorId(context.Context, int64) ([]*Video, error)
	GetVideoIdsByAuthorId(context.Context, int64) ([]int64, error)

	//GetFavorLists(ctx context.Context, ids []int64) ([]*Video, error)
	//GetVideosByUserId(ctx context.Context, uid int64) ([]Video, error)
	//AddVideoInfo(ctx context.Context, video *model.Video) error
	AddVideo(ctx context.Context, video *Video) error
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

func (r *RepoImpl) AddVideoInfo(ctx context.Context, video *Video) (err error) {
	return
}

func (r *RepoImpl) AddVideo(ctx context.Context, video *Video) (err error) {
	return
}

// GetVideoById 根据video  id查询
func (r *RepoImpl) GetVideoById(ctx context.Context, vid int64) (*Video, error) {
	// 1. 检查缓存
	key := VideoInfoPrefix + strconv.FormatInt(vid, 10)
	hit, err := r.cache.KeyExists(ctx, key)
	if err == nil && hit {
		// cache hit
		if v, err := r.cache.GetVideoById(ctx, key); err == nil {
			v.VideoId = vid
			return v, nil
		}
	}
	// 缓存命中，或者获取失败
	// 2. 从数据库读取
	v, err := r.db.GetVideoById(vid)
	if err != nil {
		return nil, err
	}
	_ = r.cache.AddVideo(ctx, key, v)
	return v, nil
}

// GetFavorLists 根据用户点赞的video id列表获取详细信息，不查询is favorite
func (r *RepoImpl) GetFavorLists(ctx context.Context, ids []int64) ([]*Video, error) {
	videos := make([]*Video, len(ids))
	var wg sync.WaitGroup
	var errCh = make(chan error, len(ids))
	wg.Add(len(ids))
	for i := range videos {
		go func(i int) {
			defer wg.Done()
			video, err := r.GetVideoById(ctx, ids[i])
			if err != nil {
				errCh <- err
				return
			}
			videos[i] = video
		}(i)
	}
	wg.Wait()
	select {
	case err := <-errCh:
		logx.Error("get video concurrency ", err)
		return nil, err
	default:
	}
	return videos, nil
}

// GetVideosByAuthorId 根据author id获取video列表
func (r *RepoImpl) GetVideosByAuthorId(ctx context.Context, uid int64) ([]*Video, error) {
	// 获取video ids
	videoIds, err := r.GetVideoIdsByAuthorId(ctx, uid)
	if err != nil {
		return nil, err
	}
	// 根据video ids获取详细信息
	videos, err := r.GetFavorLists(ctx, videoIds)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// GetVideoIdsByAuthorId 根据author id获取video id列表
func (r *RepoImpl) GetVideoIdsByAuthorId(ctx context.Context, uid int64) ([]int64, error) {
	uidStr := strconv.FormatInt(uid, 10)
	key := VideoPublishPrefix + uidStr
	hit, err := r.cache.KeyExists(ctx, key)
	if err == nil && hit {
		ids, err := r.cache.GetVideoIdsByAuthor(ctx, key)
		if err == nil {
			return ids, err
		}
	}
	// 从数据库查询
	videoIds, err := r.db.GetVideoIdsByAuthorId(uid)
	if err != nil {
		return nil, err
	}
	// 添加缓存
	err = r.cache.AddPublishList(ctx, key, videoIds)
	if err != nil {
		return nil, err
	}
	return videoIds, nil
}
