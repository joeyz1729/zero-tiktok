package data

import (
	"fmt"
	"testing"
)

func TestRepo_GetUserInfo(t *testing.T) {
	repo := NewRepo(
		"root:root1234@tcp(localhost:13306)/tiktok_user?parseTime=true&charset=utf8",
		"127.0.0.1:16379")
	for i := 1; i <= 5; i++ {
		user, err := repo.GetUserInfo(int64(i))
		if err != nil {
			panic(err)
		} else {
			fmt.Println(user)
		}
	}
}

func TestRepo_GetCount(t *testing.T) {
	repo := NewRepo(
		"root:root1234@tcp(localhost:13306)/tiktok_user?parseTime=true&charset=utf8",
		"127.0.0.1:16379")
	user := new(UserInfo)
	if err := repo.GetCount(int64(1), user); err != nil {
		panic(err)
	}
	t.Log(user)
}

func TestRepo_GetDetail(t *testing.T) {
	repo := NewRepo(
		"root:root1234@tcp(localhost:13306)/tiktok_user?parseTime=true&charset=utf8",
		"127.0.0.1:16379")
	user := new(UserInfo)
	if err := repo.GetUserDetail(int64(1), user); err != nil {
		panic(err)
	}
	t.Log(user)
}
