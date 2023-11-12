package data

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func InitRepo() *RepoImpl {
	dsn := "root:root1234@tcp(localhost:13306)/tiktok_video?parseTime=true&charset=utf8"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:16379",
	})
	return NewRepoImpl(db, rdb)
}

func Test_GetVideoById(t *testing.T) {
	repo := InitRepo()
	for i := 1; i <= 5; i++ {
		vid := int64(i)
		video, err := repo.GetVideoById(context.Background(), vid)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(video)
		}
	}
}

func TestImpl_GetVideoIdsByAuthor(t *testing.T) {
	repo := InitRepo()
	vid := int64(479264886078055299)
	fmt.Println(repo.GetVideoIdsByAuthorId(context.Background(), vid))
	for i := 1; i <= 5; i++ {
		vid := int64(i)
		videoIds, err := repo.GetVideoIdsByAuthorId(context.Background(), vid)
		if err != nil {
			fmt.Println("err: ", err)
		} else {
			fmt.Println(videoIds)
		}
	}
}

func TestRepoImpl_GetVideosByAuthorId(t *testing.T) {
	repo := InitRepo()
	vid := int64(2)
	videos, err := repo.GetVideosByAuthorId(context.Background(), vid)
	if err != nil {
		panic(err)
	}
	for _, v := range videos {
		fmt.Println(*v)
	}
	//for i := 1; i <= 5; i++ {
	//	vid := int64(i)
	//	videos, err := repo.GetVideosByAuthorId(context.Background(), vid)
	//	if err != nil {
	//		fmt.Println("err: ", err)
	//	} else {
	//		for _, v := range videos {
	//			fmt.Println(*v)
	//		}
	//	}
	//}
}

func TestRepoImpl_FeedIds(t *testing.T) {
	repo := InitRepo()
	lastTime := time.Now().Unix()
	fmt.Println(repo.FeedIds(context.Background(), lastTime))
}

func TestRepoImpl_RefreshFeed(t *testing.T) {
	repo := InitRepo()
	lastTime := time.Now().Unix()
	fmt.Println(repo.RefreshFeed(context.Background(), lastTime))
}

func TestRepoImpl_UpdateFavoriteCount(t *testing.T) {
	repo := InitRepo()
	num := 50
	var wg sync.WaitGroup
	wg.Add(num * 2)
	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			err := repo.AddFavoriteCount(context.Background(), int64(i%5+1))
			assert.Nil(t, err)
		}(i)
	}
	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			err := repo.DelFavoriteCount(context.Background(), int64(i%5+1))
			assert.Nil(t, err)
		}(i)
	}
	wg.Wait()

}
