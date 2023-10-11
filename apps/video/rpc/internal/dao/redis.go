package dao

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/go-redis/redis/v8"
)

type VideoCache interface {
	KeyExists(context.Context, string) (bool, error)
	GetVideoById(ctx context.Context, key string) (Video, error)
	GetVideosByUser(ctx context.Context, key string) ([]Video, error)
	AddVideoInfo(context.Context, model.Video) error
	DelVideo(context.Context, string) error
	AddPublishList(context.Context, string, []Video) error
	DelPublishList(context.Context, string) error
}

type RedisImpl struct {
	rdb *redis.Client
}

func NewRedisImpl(rdb *redis.Client) *RedisImpl {
	return &RedisImpl{rdb}
}

var _ VideoCache = (*RedisImpl)(nil)

func (c *RedisImpl) KeyExists(ctx context.Context, key string) (ok bool, err error) {
	num, err := c.rdb.Exists(ctx, key).Result()
	return num == 1, err
}

func (c *RedisImpl) GetVideoById(ctx context.Context, key string) (video Video, err error) {
	//var aidStr, fcStr, ccStr string
	//pipeline := c.rdb.Pipeline()
	//aidStr, err = pipeline.HGet(ctx, model.VideoInfoPrefix+vidStr, model.FieldInfoAuthorId).Result()
	//video.Title, err = pipeline.HGet(ctx, model.VideoInfoPrefix+vidStr, model.FieldInfoTitle).Result()
	//video.PlayUrl, err = pipeline.HGet(ctx, model.VideoInfoPrefix+vidStr, model.FieldInfoPlayUrl).Result()
	//video.CoverUrl, err = pipeline.HGet(ctx, model.VideoInfoPrefix+vidStr, model.FieldInfoCoverUrl).Result()
	//video.CoverUrl, err = pipeline.HGet(ctx, model.VideoInfoPrefix+vidStr, model.FieldInfoCoverUrl).Result()
	//fcStr, err = pipeline.HGet(ctx, model.VideoInfoPrefix+vidStr, model.FieldCountFavorite).Result()
	//ccStr, err = pipeline.HGet(ctx, model.VideoInfoPrefix+vidStr, model.FieldCountComment).Result()
	//_, err = pipeline.Exec(ctx)
	//if err != nil {
	//	return model.VideoDetail{}, err
	//}
	//if aid, err := strconv.ParseInt(aidStr, 10, 64); err != nil {
	//	return model.VideoDetail{}, err
	//} else {
	//	video.AuthorId = aid
	//}
	//
	//if fc, err := strconv.ParseInt(fcStr, 10, 64); err != nil {
	//	return model.Video{}, err
	//}else {video.}
	//video.AuthorId = aid
	return video, nil
}

func (c *RedisImpl) AddVideoInfo(ctx context.Context, video model.Video) (err error) {

	return
}

func (c *RedisImpl) DelVideo(ctx context.Context, key string) (err error) {
	_, err = c.rdb.Del(ctx, key).Result()
	return err
}

func (c *RedisImpl) GetVideosByUser(ctx context.Context, key string) ([]Video, error) {
	return []Video{}, nil
}

func (c *RedisImpl) DelPublishList(context.Context, string) error {
	return nil
}
func (c *RedisImpl) AddPublishList(ctx context.Context, uidStr string, videos []Video) error {
	return nil
}
