package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	dataSource = "root:root1234@tcp(localhost:13306)/tiktok_follow?parseTime=true&charset=utf8"
	addr       = "127.0.0.1:16379"
	r          *Repo
)

func init() {
	db, err := sqlx.Connect("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	r = NewRepo(db, rdb)
}

func TestRepo_CheckRelation(t *testing.T) {
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			ok, err := r.CheckRelation(1, 1)
			assert.Nil(t, err)
			assert.False(t, ok)
		}
	}

}
