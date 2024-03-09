package repository

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/cache"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/db"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/dto"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/es"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type Repo struct {
	DB  *gorm.DB
	RDB *redis.Client
	ES  *elasticsearch.TypedClient
}

var repo *Repo

func NewRepo(c config.RepoConfig) (*Repo, error) {
	database, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: c.RedisAddr,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	crt, err := os.ReadFile(c.EsCACert)
	if err != nil {
		return nil, err
	}

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: c.EsAddresses,
		Username:  c.EsUsername,
		Password:  c.EsPassword,
		CACert:    crt,
	})
	if err != nil {
		return nil, err
	}
	repo = &Repo{
		DB:  database,
		RDB: rdb,
		ES:  client,
	}
	return repo, nil
}

func (repo *Repo) GetVideoById(ctx context.Context, vid int64) (*dto.Video, error) {
	// 1. 从redis获取
	video, err := cache.GetVideo(ctx, vid, repo.RDB)
	if err != nil && err != redis.Nil {
		logx.Errorw("cache.GetVideo", logx.Field("err", err))
		return nil, err
	}
	if video != nil {
		return video, nil
	}
	// 2. 从es获取
	video, err = es.GetVideoById(ctx, vid, repo.ES)
	if err != nil {
		logx.Errorw("es.GetVideoById", logx.Field("err", err))
		return nil, err
	}
	// 3. 添加到redis
	go func() {
		if err := cache.AddVideo(ctx, video, repo.RDB); err != nil {
			logx.Errorw("cache.AddVideo", logx.Field("err", err))
		}
	}()
	return video, nil
}

func (repo *Repo) GetVideosByAuthor(ctx context.Context, uid int64) ([]*dto.Video, error) {
	// todo
	// 1. redis是否存有authorId-->set(videoId)的key
	videoIds, err := cache.GetVideoIdsByAuthor(ctx, uid, repo.RDB)
	if err != nil {
		return nil, err
	}
	var videos []*dto.Video
	if len(videoIds) == 0 {
		videos, err = es.GetVideosByAuthor(ctx, uid, repo.ES)
		if err != nil {
			logx.Errorw("es.GetVideosByAuthor", logx.Field("err", err))
			return nil, err
		}
		go func() {
			ids := make([]int64, len(videos))
			for i, video := range videos {
				ids[i] = video.ID
			}
			if err := cache.AddPublishList(ctx, uid, ids, repo.RDB); err != nil {
				logx.Errorw("cache.AddPublishList", logx.Field("err", err))
			}
		}()
		return videos, nil
	}
	videos = make([]*dto.Video, len(videoIds))
	for i := range videos {
		vid, err := strconv.ParseInt(videoIds[i], 10, 64)
		if err != nil {
			logx.Errorw("parse video id", logx.Field("err", err))
			return nil, err
		}
		video, err := repo.GetVideoById(ctx, vid)
		if err != nil {
			logx.Errorw("repo.GetVideoById", logx.Field("err", err))
			return nil, err
		}
		videos[i] = video
		go func() {
			if err := cache.AddVideo(ctx, video, repo.RDB); err != nil {
				logx.Errorw("cache.AddVideo", logx.Field("err", err))
			}
		}()
	}
	return videos, nil
}

func Feed(ctx context.Context, lastTime int64) ([]*db.Video, int64, error) {
	return nil, 0, nil
}
