package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	cache2 "github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/cache"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/db"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type Repo struct {
	db  *gorm.DB
	rdb *redis.Client
}

var repo *Repo

func NewRepo(datasource, redisAddr string) *Repo {
	database, err := gorm.Open(mysql.Open(datasource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	rdb.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	repo = &Repo{
		db:  database,
		rdb: rdb,
	}
	return repo
}

func AddVideo(ctx context.Context, video *db.Video) (err error) {
	return db.NewVideoDao(db.GlobalClient()).AddVideo(video)
}

func GetVideoById(ctx context.Context, vid int64) (*db.Video, error) {
	// 1. 检查缓存
	key := cache2.VideoInfoPrefix + strconv.FormatInt(vid, 10)
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
func GetFavorLists(ctx context.Context, ids []int64) ([]*db.Video, error) {
	videos := make([]*db.Video, len(ids))
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
		logx.Error("get videoservice concurrency ", err)
		return nil, err
	default:
	}
	return videos, nil
}

// GetVideosByAuthorId 根据author id获取video列表
func GetVideosByAuthorId(ctx context.Context, uid int64) ([]*db.Video, error) {
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
func GetVideoIdsByAuthorId(ctx context.Context, uid int64) ([]int64, error) {
	uidStr := strconv.FormatInt(uid, 10)
	key := cache2.VideoPublishPrefix + uidStr
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

func FeedIds(ctx context.Context, lastTime int64) ([]int64, int64, error) {
	// TODO
	// 查询cache是否存在

	hit, err := r.cache.KeyExists(ctx, cache2.VideoFeedKey)
	if err == nil && hit {
		ids, nextTime, err := r.cache.GetFeedIds(ctx, lastTime)
		if err == nil {
			return ids, nextTime, nil
		}
	}
	// 如果不存在或者出错则走database
	videos, err := r.db.GetFeedIds(lastTime)
	if err != nil {
		return nil, 0, err
	}
	// 将查询结果添加到cache中
	nextTime := lastTime
	videoIds := make([]int64, len(videos))
	for i, v := range videos {
		videoIds[i] = v.VideoId
		vt := v.PublishTime.Unix()
		if vt < nextTime {
			nextTime = vt
		}
		if err := r.cache.AddFeedVideo(ctx, v.VideoId, vt); err != nil {
			return nil, 0, err
		}
	}
	return videoIds, nextTime, nil

}

func RefreshFeed(ctx context.Context, lastTime int64) error {
	// 可以换成查询refreshTime之后的，添加一个limit，然后返回下一个游标，更新refreshTime
	vs, err := r.db.GetFeedIds(lastTime)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(vs))
	errCh := make(chan error, len(vs))
	for _, v := range vs {
		go func(vid int64, publishTime int64) {
			defer wg.Done()
			if err := r.cache.AddFeedVideo(ctx, vid, publishTime); err != nil {
				errCh <- err
				return
			}
		}(v.VideoId, v.PublishTime.Unix())
	}
	wg.Wait()
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
