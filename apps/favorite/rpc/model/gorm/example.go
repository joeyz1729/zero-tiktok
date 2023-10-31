package main

import (
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(localhost:3306)/tiktok_favorite?parseTime=true&charset=utf8"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.Favorite{})
	if err != nil {
		panic(err)
	}
	fmt.Println("success")
}
