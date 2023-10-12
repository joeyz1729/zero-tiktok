package dao

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRedisImpl_KeyExists(t *testing.T) {
	repo := InitRepo()
	for i := 1; i <= 6; i++ {
		vid := strconv.FormatInt(int64(i), 10)
		video, err := repo.cache.KeyExists(context.Background(), VideoInfoPrefix+vid)
		if err != nil {
			fmt.Println("err: ", err)
		} else {
			fmt.Println(video)
		}
	}
}

func TestRedisImpl_GetVideoById(t *testing.T) {
	repo := InitRepo()
	for i := 1; i <= 6; i++ {
		vid := strconv.FormatInt(int64(i), 10)
		video, err := repo.cache.GetVideoById(context.Background(), VideoInfoPrefix+vid)
		if err != nil {
			fmt.Println("err: ", err)
		} else {
			fmt.Println(video)
		}
	}
}

func TestRedisImpl_GetFeedIds(t *testing.T) {
	repo := InitRepo()
	lastTime := time.Now().Unix()
	//lastTime = 1697119044
	hit, err := repo.cache.KeyExists(context.Background(), VideoFeedKey)
	if err == nil && hit {
		ids, nextTime, err := repo.cache.GetFeedIds(context.Background(), lastTime)
		if err == nil {
			fmt.Println("[cache hit] ", ids, nextTime)
		} else {
			panic(err)
		}

	}
	if err != nil {
		panic(err)
	}
}

func TestRedisImpl_AddFeedVideo(t *testing.T) {
	repo := InitRepo()
	var uid, stamp int64
	uid = time.Now().Unix()
	stamp = time.Now().Unix()
	if err := repo.cache.AddFeedVideo(context.Background(), uid, stamp); err != nil {
		fmt.Println(err)
	}
	fmt.Println("success")
}
