package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"testing"
)

func InitRepo() *RepoImpl {
	dsn := "root:Zy_9908091729@tcp(localhost:3306)/tiktok_video?parseTime=true&charset=utf8"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
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
