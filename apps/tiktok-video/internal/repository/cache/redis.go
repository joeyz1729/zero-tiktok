package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/db"
	"strconv"
)

var (
	VideoFeedRefreshKey = "tiktok:video::refresh"
	VideoFeedKey        = "tiktok:video::feed"    // +nil zset	(vid, timestamp)
	VideoInfoPrefix     = "tiktok:video:info:"    // +vid, hash	(info)
	VideoPublishPrefix  = "tiktok:video:publish:" // +uid set (vid)
	FieldInfoTitle      = "title"

	FieldInfoPlayUrl   = "playurl"
	FieldInfoCoverUrl  = "coverurl"
	FieldInfoAuthorId  = "authorid"
	FieldCountFavorite = "favorite"

	ErrInvalidType = errors.New("invalid type")
	ErrEmptySet    = errors.New("empty set")
)

var globalRdb *redis.Client

func InitRdb(rdb *redis.Client) {
	globalRdb = rdb
}

func DelKey(ctx context.Context, key string) error {
	_, err := globalRdb.Del(ctx, key).Result()
	return err
}

func KeyExists(ctx context.Context, key string) (ok bool, err error) {
	num, err := globalRdb.Exists(ctx, key).Result()
	return num == 1, err
}

func GetVideoById(ctx context.Context, key string) (video *db.Video, err error) {
	cm, err := globalRdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(cm) != 6 {
		return nil, ErrInvalidType
	}
	video = new(db.Video)
	video.AuthorID, err = strconv.ParseInt(cm[FieldInfoAuthorId], 10, 64)
	if err != nil {
		return nil, err
	}

	video.Title = cm[FieldInfoTitle]
	video.CoverURL = cm[FieldInfoCoverUrl]
	video.PlayURL = cm[FieldInfoPlayUrl]
	return video, nil
}

func AddVideo(ctx context.Context, key string, video *db.Video) (err error) {
	_, err = globalRdb.HSet(ctx, key,
		FieldInfoAuthorId, video.AuthorID,
		FieldInfoTitle, video.Title,
		FieldInfoPlayUrl, video.PlayURL,
		FieldInfoCoverUrl, video.CoverURL,
	).Result()
	return err
}

func DelVideo(ctx context.Context, key string) (err error) {
	_, err = globalRdb.Del(ctx, key).Result()
	return err
}

func GetVideosByUser(ctx context.Context, key string) ([]*db.Video, error) {
	return []*db.Video{}, nil
}

func DelPublishList(context.Context, string) error {
	return nil
}

// GetVideoIdsByAuthor 根据作者id获取video id列表
func GetVideoIdsByAuthor(ctx context.Context, key string) ([]int64, error) {
	ids, err := globalRdb.SMembers(ctx, key).Result()
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

func AddPublishList(ctx context.Context, uidStr string, videoIds []int64) error {
	return nil
}

// GetFeedIds 通过cache获取vid列表，
func GetFeedIds(ctx context.Context, lastTime int64) ([]int64, int64, error) {
	// limit 30，以及获取结果中最小的时间戳
	zs, err := globalRdb.ZRevRangeByScoreWithScores(ctx, VideoFeedKey,
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

func AddFeedVideo(ctx context.Context, vid int64, stamp int64) error {
	_, err := globalRdb.ZAdd(ctx, VideoFeedKey, &redis.Z{Member: strconv.FormatInt(vid, 10), Score: float64(stamp)}).Result()
	return err
}
