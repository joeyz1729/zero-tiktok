package data

import (
	"context"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/config"
	"strconv"
	"testing"
)

func InitRepo() *RepoImpl {
	repo, err := NewRepo(config.Config{
		Mysql:      struct{ DataSource string }{DataSource: "root:root@tcp(localhost:3306)/tiktok_favorite?parseTime=true&charset=utf8"},
		CacheRedis: struct{ Addr string }{Addr: "127.0.0.1:6379"},
	})
	if err != nil {
		panic(err)
	}
	return repo
}

func TestDBImpl_GetFavoriteIds(t *testing.T) {
	repo := InitRepo()
	ids, err := repo.db.GetFavoriteIds(int64(482313611805467523))
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)

	for i := 1; i <= 5; i++ {
		ids, err := repo.db.GetFavoriteIds(int64(i))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ids)
	}
}

func TestCacheImpl_GetFavoriteIds(t *testing.T) {
	repo := InitRepo()
	uid := strconv.FormatInt(int64(482313611805467523), 10)
	hit, err := repo.cache.KeyExist(context.Background(), FavoriteSetPrefix+uid)
	fmt.Println(hit, err)
	ids, err := repo.cache.GetFavoriteVideoIds(context.Background(), FavoriteSetPrefix+uid)
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)

	for i := 1; i <= 5; i++ {
		uid := strconv.FormatInt(int64(i), 10)
		ids, err := repo.cache.GetFavoriteVideoIds(context.Background(), FavoriteSetPrefix+uid)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ids)
	}
}

func TestRepoImpl_GetFavorIds(t *testing.T) {
	repo := InitRepo()
	ids, err := repo.GetFavorIds(context.Background(), int64(482313611805467523))
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)

	for i := 1; i <= 5; i++ {
		ids, err := repo.GetFavorIds(context.Background(), int64(i))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ids)
	}
}

func TestRepoImpl_IsFavoriteRecordExist(t *testing.T) {
	repo := InitRepo()
	pairs := [][2]int64{
		{1, 1},
		{1, 2},
		{2, 1},
		{2, 2},
		{482313611805467523, 1},
		{482313611805467523, 2},
		{482313611805467523, 3},
		{482313611805467523, 479275561387041668},
		{482313611805467523, 479275563937178500},
	}
	for _, p := range pairs {
		fmt.Println(repo.IsFavoriteRecordExist(context.Background(), p[0], p[1]))
	}
}
