package model

import (
	"fmt"
	"testing"
)

func TestRepo_GetUserInfo(t *testing.T) {
	repo := NewRepo("root:Zy_9908091729@tcp(localhost:3306)/tiktok_user?parseTime=true&charset=utf8", "127.0.0.1:6379")
	for i := 1; i <= 5; i++ {
		user, err := repo.GetUserInfo(int64(i))
		if err != nil {
			panic(err)
		} else {
			fmt.Println(user)
		}
	}
}
