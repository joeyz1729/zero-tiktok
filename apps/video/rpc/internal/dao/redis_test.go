package dao

import (
	"context"
	"fmt"
	"strconv"
	"testing"
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
