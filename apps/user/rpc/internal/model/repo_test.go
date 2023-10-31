package model

import (
	"fmt"
	"testing"
)

func TestRepo_GetUserInfo(t *testing.T) {
	repo := NewRepo("root:root@tcp(localhost:3306)/tiktok_user?parseTime=true&charset=utf8", "127.0.0.1:6379")
	fmt.Println(repo.GetUserInfo(int64(479003589444903811)))
	fmt.Println(repo.GetUserInfo(int64(482472807905633155)))
	for i := 1; i <= 5; i++ {
		user, err := repo.GetUserInfo(int64(i))
		if err != nil {
			panic(err)
		} else {
			fmt.Println(user)
		}
	}
}
