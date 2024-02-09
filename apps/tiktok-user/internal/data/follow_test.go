package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepo_AddFollow(t *testing.T) {
	repo := NewRepo(
		"root:root1234@tcp(localhost:13306)/tiktok_user?parseTime=true&charset=utf8",
		"127.0.0.1:16379")
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			err := repo.AddFollow(int64(i), int64(j))
			if err != nil {
				panic(err)
			}
		}
	}
	for i := 1; i <= 5; i++ {
		user, err := repo.GetUserInfo(int64(i))
		assert.Nil(t, err)
		t.Log(user.FollowedCount, user.FollowerCount)
	}
}

func TestRepo_DelFollow(t *testing.T) {
	repo := NewRepo(
		"root:root1234@tcp(localhost:13306)/tiktok_user?parseTime=true&charset=utf8",
		"127.0.0.1:16379")
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			err := repo.DelFollow(int64(i), int64(j))
			if err != nil {
				panic(err)
			}
		}
	}
	for i := 1; i <= 5; i++ {
		user, err := repo.GetUserInfo(int64(i))
		assert.Nil(t, err)
		t.Log(user.FollowedCount, user.FollowerCount)
	}
}
