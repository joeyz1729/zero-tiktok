package data

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
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

func Test_AddVideos(t *testing.T) {
	repo := InitRepo()
	v := new(Video)
	for i := 10; i <= 20; i++ {
		v.VideoId = int64(i%5 + 1)
		v.AuthorId = int64(i%5 + 1)
		v.Title = fmt.Sprintf("test video %d", i%5+1)
		err := repo.AddVideo(context.Background(), v)
		if err != nil {
			panic(err)
		}
	}

}

func Test_UpdateCount(t *testing.T) {
	repo := InitRepo()
	for i := 1; i <= 5; i++ {
		err := repo.AddFavoriteCount(context.Background(), int64(i))
		assert.Nil(t, err)
	}
	for i := 1; i <= 5; i++ {
		err := repo.DelFavoriteCount(context.Background(), int64(i))
		assert.Nil(t, err)
	}
}
