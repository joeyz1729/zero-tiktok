package data

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var (
	VideoFeedRefreshKey      = "tiktok:video::refresh"
	VideoFeedKey             = "tiktok:video::feed"     // +nil zset	(vid, timestamp)
	VideoInfoPrefix          = "tiktok:video:info:"     // +vid, hash	(info)
	VideoPublishPrefix       = "tiktok:video:publish:"  // +uid set (vid)
	VideoFavoriteCountPrefix = "tiktok:video:favorite:" // +vid, (favorite count)
	VideoCommentCountPrefix  = "tiktok:video:comment:"
	FieldInfoTitle           = "title"

	FieldInfoPlayUrl   = "playurl"
	FieldInfoCoverUrl  = "coverurl"
	FieldInfoAuthorId  = "authorid"
	FieldCountFavorite = "favorite"

	ErrInvalidType = errors.New("invalid type")
	ErrEmptySet    = errors.New("empty set")
)

type VideoCache interface {
	KeyExists(context.Context, string) (bool, error)

	DelVideo(context.Context, string) error
	AddVideo(context.Context, string, *Video) error

	DelKey(context.Context, string) error

	GetVideoById(ctx context.Context, key string) (*Video, error)

	GetVideosByUser(ctx context.Context, key string) ([]*Video, error)

	AddPublishList(context.Context, string, []int64) error
	DelPublishList(context.Context, string) error
	AddFeedVideo(context.Context, int64, int64) error

	GetVideoIdsByAuthor(context.Context, string) ([]int64, error)
	GetFeedIds(context.Context, int64) ([]int64, int64, error)
}

type RedisImpl struct {
	rdb *redis.Client
}

func NewRedisImpl(rdb *redis.Client) *RedisImpl {
	return &RedisImpl{rdb}
}

var _ VideoCache = (*RedisImpl)(nil)

func (c *RedisImpl) DelKey(ctx context.Context, key string) error {
	_, err := c.rdb.Del(ctx, key).Result()
	return err
}

func (c *RedisImpl) KeyExists(ctx context.Context, key string) (ok bool, err error) {
	num, err := c.rdb.Exists(ctx, key).Result()
	return num == 1, err
}

func (c *RedisImpl) GetVideoById(ctx context.Context, key string) (video *Video, err error) {
	cm, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(cm) != 6 {
		return nil, ErrInvalidType
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
	//video.CommentCount, err = strconv.ParseInt(cm[Field], 10, 64)
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
		//FieldCountComment, video.CommentCount,
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

// GetVideoIdsByAuthor 根据作者id获取video id列表
func (c *RedisImpl) GetVideoIdsByAuthor(ctx context.Context, key string) ([]int64, error) {
	ids, err := c.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	res := make([]int64, len(ids))
	for i, idStr := range ids {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = id
	}
	return res, nil
}

func (c *RedisImpl) AddPublishList(ctx context.Context, uidStr string, videoIds []int64) error {
	return nil
}

// GetFeedIds 通过cache获取vid列表，
func (c *RedisImpl) GetFeedIds(ctx context.Context, lastTime int64) ([]int64, int64, error) {
	// limit 30，以及获取结果中最小的时间戳
	zs, err := c.rdb.ZRevRangeByScoreWithScores(ctx, VideoFeedKey,
		&redis.ZRangeBy{
			Max:    strconv.FormatInt(lastTime, 10),
			Offset: 0,
			Count:  30,
		}).Result()
	if err != nil {
		return nil, 0, err
	}
	if len(zs) == 0 {
		return nil, 0, ErrEmptySet
	}
	nextTime := float64(lastTime)
	vids := make([]int64, len(zs))
	for i, z := range zs {
		id, err := strconv.ParseInt(z.Member.(string), 10, 64)
		if err != nil {
			return nil, 0, err
		}
		vids[i] = id
		if z.Score < nextTime {
			nextTime = z.Score
		}
	}
	return vids, int64(nextTime), nil
}

func (r *RedisImpl) AddFeedVideo(ctx context.Context, vid int64, stamp int64) error {
	//TODO
	_, err := r.rdb.ZAdd(ctx, VideoFeedKey, &redis.Z{Member: strconv.FormatInt(vid, 10), Score: float64(stamp)}).Result()
	return err
}
