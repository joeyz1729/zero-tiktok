package data

import (
	"github.com/jmoiron/sqlx"
	"testing"
)

var (
	dsn  = "root:root1234@tcp(localhost:13306)/tiktok_comment?parseTime=true&charset=utf8"
	addr = "127.0.0.1:16379"
)

func init() {

}

func TestDel(t *testing.T) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	count := 0
	var vid, cid, uid int64 = 1, 1, 2
	getStr := `select count(*) from tiktok_comment.comment where video_id = ?  and comment_id = ? and user_id = ? limit 1`
	err = db.Get(&count, getStr, vid, cid, uid)
	t.Log(count, err)
}
