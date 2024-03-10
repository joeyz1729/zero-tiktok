package dto

type User struct {
	ID              int64  `json:"id"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`

	TotalFavorited int64 `json:"total_favorited"`
	WorkCount      int64 `json:"work_count"`
	FavoriteCount  int64 `json:"favorite_count"`
	FollowCount    int64 `json:"follow_count"`
	FollowerCount  int64 `json:"follower_count"`

	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
