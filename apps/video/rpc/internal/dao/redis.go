package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var (
	VideoFeedPrefix    = "tiktok:video::feed"    // +nil zset	(vid, timestamp)
	VideoInfoPrefix    = "tiktok:video:info:"    // +vid, hash	(info)
	VideoPublishPrefix = "tiktok:video:publish:" // +uid set (vid)

	FieldInfoTitle     = "title"
	FieldInfoPlayUrl   = "playurl"
	FieldInfoCoverUrl  = "coverurl"
	FieldInfoAuthorId  = "authorid"
	FieldCountFavorite = "favorite"
	FieldCountComment  = "comment"
)

type VideoCache interface {
	KeyExists(context.Context, string) (bool, error)

	DelVideo(context.Context, string) error
	AddVideo(context.Context, string, *Video) error

	GetVideoById(ctx context.Context, key string) (*Video, error)

	GetVideosByUser(ctx context.Context, key string) ([]*Video, error)

	AddPublishList(context.Context, string, []*Video) error
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

func (c *RedisImpl) GetVideoById(ctx context.Context, key string) (video *Video, err error) {
	cm, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	video = new(Video)
	video.AuthorId, err = strconv.ParseInt(cm[FieldInfoAuthorId], 10, 64)
	if err != nil {
		return nil, err
	}
	video.FavoriteCount, err = strconv.ParseInt(cm[FieldCountFavorite], 10, 64)
	if err != nil {
		return nil, err
	}
	video.CommentCount, err = strconv.ParseInt(cm[FieldCountComment], 10, 64)
	if err != nil {
		return nil, err
	}
	video.Title = cm[FieldInfoTitle]
	video.CoverUrl = cm[FieldInfoCoverUrl]
	video.PlayUrl = cm[FieldInfoPlayUrl]
	return video, nil
}

func (c *RedisImpl) AddVideo(ctx context.Context, key string, video *Video) (err error) {
	_, err = c.rdb.HSet(ctx, key,
		FieldInfoAuthorId, video.AuthorId,
		FieldInfoTitle, video.Title,
		FieldInfoPlayUrl, video.PlayUrl,
		FieldInfoCoverUrl, video.CoverUrl,
		FieldCountFavorite, video.FavoriteCount,
		FieldCountComment, video.CommentCount,
	).Result()
	return err
}

func (c *RedisImpl) DelVideo(ctx context.Context, key string) (err error) {
	_, err = c.rdb.Del(ctx, key).Result()
	return err
}

func (c *RedisImpl) GetVideosByUser(ctx context.Context, key string) ([]*Video, error) {
	return []*Video{}, nil
}

func (c *RedisImpl) DelPublishList(context.Context, string) error {
	return nil
}
func (c *RedisImpl) AddPublishList(ctx context.Context, uidStr string, videos []*Video) error {
	return nil
}
