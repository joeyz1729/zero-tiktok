package dto

import "time"

type Video struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AuthorID     int64     `gorm:"column:author_id;not null" json:"author_id"`
	Title        string    `gorm:"column:title;not null" json:"title"`
	PlayURL      string    `gorm:"column:play_url;not null" json:"play_url"`
	CoverURL     string    `gorm:"column:cover_url;not null" json:"cover_url"`
	ThumbupCount int64     `gorm:"column:thumbup_count;not null" json:"thumbup_count"`
	CommentCount int64     `gorm:"column:comment_count;not null" json:"comment_count"`
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}
