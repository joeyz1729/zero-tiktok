package data

import (
	"context"
	"fmt"
	"github.com/joeyz1729/zero-tiktok/apps/favorite/rpc/model"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

var (
	dsn       = "root:root1234@tcp(localhost:13307)/tiktok_favorite?parseTime=true&charset=utf8"
	redisAddr = "127.0.0.1:16379"
)

func InitRepo() *RepoImpl {
	repo, err := NewRepo(
		dsn, redisAddr)
	if err != nil {
		panic(err)
	}
	return repo
}

func Test_Connect(t *testing.T) {
	repo, err := NewRepo(dsn, redisAddr)
	assert.Nil(t, err)
	assert.NotNil(t, repo)
}

func TestDBImpl_GetFavoriteIds(t *testing.T) {
	repo := InitRepo()
	userId := 486525991745758083
	ids, err := repo.GetFavoriteIdsFromDB(int64(userId))
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)

	ids, err = repo.db.GetFavoriteIds(486525991745758083)
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)

}
func TestRepoImpl_GetFavorIds(t *testing.T) {
	repo := InitRepo()
	userId := 482313611805467523
	ids, err := repo.GetFavorIds(context.Background(), int64(userId))
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)

	ids, err = repo.db.GetFavoriteIds(int64(482313611805467523))
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)

}

func TestDBGet(t *testing.T) {
	db, err := NewDBImpl(dsn)
	if err != nil {
		panic(err)
	}
	userId := 486525991745758083
	var favors []*model.Favorite
	err = db.Table("favorite").Select("video_id").Where("user_id = ?", userId).Find(&favors).Error
	if err != nil {
		t.Log(err)
	}
	ids := make([]int64, len(favors))
	for i, f := range favors {
		ids[i] = f.VideoId
	}
	t.Log(ids)
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
		fmt.Println(repo.CheckFavor(context.Background(), p[0], p[1]))
	}
}

func TestRepoImpl_CreateFavoriteRecord(t *testing.T) {
	repo := InitRepo()
	num := 5 * 5
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			go func(i, j int) {
				defer wg.Done()
				err := repo.CreateFavoriteRecord(context.Background(), &model.Favorite{UserId: int64(i), VideoId: int64(j)})
				assert.Nil(t, err)
			}(i, j)
		}
	}
	wg.Wait()
	wg.Add(num)
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			go func(i, j int) {
				defer wg.Done()
				err := repo.DeleteFavoriteRecord(context.Background(), &model.Favorite{UserId: int64(i), VideoId: int64(j)})
				assert.Nil(t, err)
			}(i, j)
		}
	}
	wg.Wait()

}
