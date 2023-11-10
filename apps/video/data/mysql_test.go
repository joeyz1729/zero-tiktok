package data

import (
	"fmt"
	"testing"
	"time"
)

func TestMysqlImpl_GetFeedIds(t *testing.T) {
	repo := InitRepo()
	lastTime := time.Now().Unix()
	videos, err := repo.db.GetFeedIds(lastTime)
	if err != nil {
		panic(err)
	}
	for _, v := range videos {
		fmt.Println(v.VideoId, v.PublishTime, v.PublishTime.Unix())
	}
}
