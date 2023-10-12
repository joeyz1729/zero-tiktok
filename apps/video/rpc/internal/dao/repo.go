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

func (r *RepoImpl) GetVideosByUserId(ctx context.Context, uid int64) ([]Video, error) {
	uidStr := strconv.FormatInt(uid, 10)
	hit, err := r.cache.KeyExists(ctx, VideoPublishPrefix+uidStr)
	if err == nil && hit {
		// cache hit
		//return r.cache.GetVideosByUser(ctx, model.VideoPublishPrefix+uidStr)
	}

	// 从数据库读取
	videos, err := r.db.GetVideoByUser(uid)
	if err != nil {
		return nil, err
	}
	// 添加用户的发布列表缓存
	// 添加视频信息的缓存
	_ = r.cache.AddPublishList(ctx, uidStr, videos)
	//return videos, nil
	return []Video{}, nil
}

func (r *RepoImpl) AddVideoInfo(ctx context.Context, video *Video) (err error) {
	return
}

func (r *RepoImpl) AddVideo(ctx context.Context, video *Video) (err error) {
	return
}

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
