package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("mysql", "root:Zy_9908091729@tcp(localhost:3306)/tiktok_user?parseTime=true&charset=utf8")
	if err != nil {
		panic(err)
	}
	var userIds []int64
	sqlStr := `select user_id from tiktok_user.user`
	err = db.Select(&userIds, sqlStr)
	if err != nil {
		panic(err)
	}
	for _, uid := range userIds {
		sqlStr = `insert into tiktok_user.user_count(user_id) value(?)`
		if _, err = db.Exec(sqlStr, uid); err != nil {
			fmt.Println(err)
			continue
		}
	}

}
