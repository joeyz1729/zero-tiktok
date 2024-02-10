package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/cache"
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
	return db.AddVideo(video)
}

func Feed(ctx context.Context, lastTime int64) ([]*db.Video, int64, error) {
	return nil, 0, nil
}

func GetVideoById(ctx context.Context, vid int64) (*db.Video, error) {
	// 1. 检查缓存
	key := cache.VideoInfoPrefix + strconv.FormatInt(vid, 10)
	hit, err := cache.KeyExists(ctx, key)
	if err == nil && hit {
		// cache hit
		if v, err := cache.GetVideoById(ctx, key); err == nil {
			v.ID = vid
			return v, nil
		}
	}
	// 缓存命中，或者获取失败
	// 2. 从数据库读取
	v, err := db.GetVideoById(vid)
	if err != nil {
		return nil, err
	}
	_ = cache.AddVideo(ctx, key, v)
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
			video, err := cache.GetVideoById(ctx, "")
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
	//// 获取video ids
	//videoIds, err := cache.GetVideoIdsByAuthor(ctx, "")
	//if err != nil {
	//	return nil, err
	//}
	//// 根据video ids获取详细信息
	////videos, err := cache.GetFavorLists(ctx, videoIds)
	//if err != nil {
	//	return nil, err
	//}
	//return videos, nil
	return nil, nil
}

// GetVideoIdsByAuthorId 根据author id获取video id列表
func GetVideoIdsByAuthorId(ctx context.Context, uid int64) ([]int64, error) {
	uidStr := strconv.FormatInt(uid, 10)
	key := cache.VideoPublishPrefix + uidStr
	hit, err := cache.KeyExists(ctx, key)
	if err == nil && hit {
		ids, err := cache.GetVideoIdsByAuthor(ctx, key)
		if err == nil {
			return ids, err
		}
	}
	// 从数据库查询
	videos, err := db.GetVideosByAuthorId(uid)
	if err != nil {
		return nil, err
	}
	// 添加缓存
	var videoIds []int64
	for _, v := range videos {
		videoIds = append(videoIds, v.ID)
	}

	err = cache.AddPublishList(ctx, key, videoIds)
	if err != nil {
		return nil, err
	}
	return videoIds, nil
}
