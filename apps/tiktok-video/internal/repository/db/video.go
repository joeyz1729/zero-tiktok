package db

import (
	"context"
	"gorm.io/gorm"
	"time"
)

const TableNameVideo = "video"

// Video 视频表
type Video struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AuthorID   int64     `gorm:"column:author_id;not null" json:"author_id"`
	Title      string    `gorm:"column:title;not null" json:"title"`
	PlayURL    string    `gorm:"column:play_url;not null" json:"play_url"`
	CoverURL   string    `gorm:"column:cover_url;not null" json:"cover_url"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName Video's table name
func (*Video) TableName() string {
	return TableNameVideo
}

func CreateVideo(ctx context.Context, video *Video, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(TableNameVideo).Create(video).Error; err != nil {
			return err
		}
		var videoCount = VideoCount{
			ID: video.ID,
		}
		if err := tx.Table(TableNameVideoCount).Create(&videoCount).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
