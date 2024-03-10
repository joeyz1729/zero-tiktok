package db

import (
	"context"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const TableNameVideoCount = "video_count"

// VideoCount 视频信息计数表
type VideoCount struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ThumbupCount int64     `gorm:"column:thumbup_count;not null" json:"thumbup_count"`
	CommentCount int64     `gorm:"column:comment_count;not null" json:"comment_count"`
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName VideoCount's table name
func (*VideoCount) TableName() string {
	return TableNameVideoCount
}

func CreateVideoCount(ctx context.Context, data map[string]interface{}, db *gorm.DB) error {
	videoId, err := strconv.ParseInt(data["id"].(string), 10, 64)
	if err != nil {
		return err
	}
	layout := "2006-01-02 15:04:05"
	createdTime, err := time.Parse(layout, data["create_time"].(string))
	if err != nil {
		return err
	}
	var count = VideoCount{
		ID:         videoId,
		CreateTime: createdTime,
		UpdateTime: createdTime,
	}
	return db.Table(TableNameVideoCount).Create(&count).Error
}

func UpdateThumbupCount(videoId int64, incr int64, DB *gorm.DB) error {
	err := DB.Table(TableNameVideoCount).Transaction(func(tx *gorm.DB) error {
		var videoCount VideoCount
		if err := tx.First(&videoCount, videoId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", videoId).
			Updates(map[string]interface{}{"thumbup_count": videoCount.ThumbupCount + incr}).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return err
}

func UpdateCommentCount(videoId int64, incr int64, DB *gorm.DB) error {
	err := DB.Table(TableNameVideoCount).Transaction(func(tx *gorm.DB) error {
		var videoCount VideoCount
		if err := tx.First(&videoCount, videoId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", videoId).
			Updates(map[string]interface{}{"comment_count": videoCount.CommentCount + incr}).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return err
}
