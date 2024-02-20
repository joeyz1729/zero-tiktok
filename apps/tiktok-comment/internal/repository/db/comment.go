package db

import (
	"gorm.io/gorm"
	"time"
)

const TableNameComment = "comment"

// Comment 用户表
type Comment struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	VideoID    int64     `gorm:"column:video_id;not null" json:"video_id"`
	UserID     int64     `gorm:"column:user_id;not null" json:"user_id"`
	Content    string    `gorm:"column:content;not null" json:"content"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}

type CommentDB struct {
	*gorm.DB
}

func NewCommentDB(db *gorm.DB) *CommentDB {
	return &CommentDB{
		db,
	}
}

func (db *CommentDB) Add(comment *Comment) error {
	return db.Table(TableNameComment).Create(comment).Error
}

func (db *CommentDB) Delete(commentId int64) error {
	return db.Table(TableNameComment).Where("id = ?", commentId).Delete(&Comment{}).Error
}

func (db *CommentDB) Get(commentId int64) (*Comment, error) {
	var comment Comment
	err := db.Table(TableNameComment).Where("id = ?", commentId).First(&Comment{}).Error
	return &comment, err
}
