package data

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var ErrNotFound = sqlx.ErrNotFound

type UserInfo struct {
	Id              int64  `db:"id"`
	UserId          int64  `db:"user_id"`
	Username        string `db:"username"`
	Password        string `db:"password"`
	Avatar          string `db:"avatar"`
	BackgroundImage string `db:"background_image"`
	Signature       string `db:"signature"`

	FollowedCount int64 `db:"followed_count" json:"followedcount,string"` // 关注总数
	FollowerCount int64 `db:"follower_count" json:"followercount,string"` // 粉丝总数

	TotalFavorited int64     `db:"total_favorited" json:"totalfavorited,string"` //获赞数量
	WorkCount      int64     `db:"work_count" json:"workcount,string"`           //作品数量
	FavoriteCount  int64     `db:"favorite_count" json:"favoritecount,string"`   //点赞数量
	CreateTime     time.Time `db:"create_time"`
	UpdateTime     time.Time `db:"update_time"`
}

type UserCount struct {
	UserId         int64 `db:"user_id"`
	FollowedCount  int64 `db:"followed_count" json:"followedcount,string"`   // 关注总数
	FollowerCount  int64 `db:"follower_count" json:"followercount,string"`   // 粉丝总数
	TotalFavorited int64 `db:"total_favorited" json:"totalfavorited,string"` //获赞数量
	WorkCount      int64 `db:"work_count" json:"workcount,string"`           //作品数量
	FavoriteCount  int64 `db:"favorite_count" json:"favoritecount,string"`   //点赞数量
}
