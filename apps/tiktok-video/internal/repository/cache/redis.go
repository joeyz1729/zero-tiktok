package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/dto"
	"strconv"
	"time"
)

var (
	VideoFeedKey       = "tiktok:video:feed"     // +nil zset	(vid, timestamp)
	VideoInfoPrefix    = "tiktok:video:info:"    // +vid, hash	(info)
	VideoPublishPrefix = "tiktok:video:publish:" // +uid set (vid)

	NormalTTL = time.Hour * 24
)

func getPublishListKey(uid int64) string {
	return VideoPublishPrefix + strconv.FormatInt(uid, 10)
}

func getVideoKey(vid int64) string {
	return VideoInfoPrefix + strconv.FormatInt(vid, 10)
}

func AddVideo(ctx context.Context, video *dto.Video, rdb *redis.Client) error {
	data, err := json.Marshal(video)
	if err != nil {
		return err
	}
	_, err = rdb.Set(ctx, getVideoKey(video.ID), data, NormalTTL).Result()
	return err
}

func AddFeed(ctx context.Context, video map[string]interface{}, rdb *redis.Client) error {
	layout := "2006-01-02 15:04:05"
	createdTime, err := time.Parse(layout, video["create_time"].(string))
	if err != nil {
		return err
	}
	_, err = rdb.ZAdd(ctx, VideoFeedKey, &redis.Z{
		Score:  float64(createdTime.Unix()),
		Member: video["id"].(string),
	}).Result()
	return err
}

func AddPublishList(ctx context.Context, uid int64, videoIds []int64, rdb *redis.Client) error {
	_, err := rdb.SAdd(ctx, getPublishListKey(uid), videoIds).Result()
	if err != nil {
		return err
	}
	_, err = rdb.Expire(ctx, getPublishListKey(uid), NormalTTL).Result()
	return err
}

func GetVideoIdsByAuthor(ctx context.Context, userId int64, rdb *redis.Client) ([]string, error) {
	ids, err := rdb.SMembers(ctx, getPublishListKey(userId)).Result()
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func GetVideo(ctx context.Context, vid int64, rdb *redis.Client) (*dto.Video, error) {
	data, err := rdb.Get(ctx, getVideoKey(vid)).Result()
	if err != nil {
		return nil, err
	}
	var video dto.Video
	if err = json.Unmarshal([]byte(data), &video); err != nil {
		return nil, err
	}
	return &video, nil
}

func GetFeedIds(ctx context.Context, lastTime int64, rdb *redis.Client) ([]int64, int64, error) {
	// limit 30，以及获取结果中最小的时间戳
	zSlice, err := rdb.ZRevRangeByScoreWithScores(ctx, VideoFeedKey,
		&redis.ZRangeBy{
			Max:    strconv.FormatInt(lastTime, 10),
			Offset: 0,
			Count:  30,
		}).Result()
	if err != nil {
		return nil, 0, err
	}

	nextTime := float64(lastTime)

	ids := make([]int64, len(zSlice))
	for i, z := range zSlice {
		id, err := strconv.ParseInt(z.Member.(string), 10, 64)
		if err != nil {
			return nil, 0, err
		}
		ids[i] = id
		if z.Score < nextTime {
			nextTime = z.Score
		}
	}
	return ids, int64(nextTime), nil
}
