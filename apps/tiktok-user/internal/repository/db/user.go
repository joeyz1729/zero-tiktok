package db

import (
	"time"
)

const TableNameUser = "user"

// User 用户表
type User struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username        string    `gorm:"column:username;not null" json:"username"`
	Password        string    `gorm:"column:password;not null" json:"password"`
	Avatar          string    `gorm:"column:avatar;not null" json:"avatar"`
	BackgroundImage string    `gorm:"column:background_image;not null" json:"background_image"`
	Signature       string    `gorm:"column:signature;not null" json:"signature"`
	CreateTime      time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
