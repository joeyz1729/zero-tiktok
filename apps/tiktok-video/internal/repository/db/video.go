package db

import (
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

type VideoDao struct {
	db *gorm.DB
}

func NewVideoDao(db *gorm.DB) *VideoDao {
	return &VideoDao{db: db}
}

var globalDBClient *gorm.DB

func GlobalClient() *gorm.DB {
	return globalDBClient
}

func (s *VideoDao) AddVideo(v *Video) error {
	return s.db.Create(v).Error
}

func (s *VideoDao) GetVideoById(vid int64) (*Video, error) {
	var v Video
	err := s.db.Where("id = ?", vid).First(&v).Error
	return &v, err
}

func (s *VideoDao) GetVideosByAuthorId(uid int64) ([]*Video, error) {
	var videoList []*Video
	err := s.db.Where("author_id = ?", uid).Find(&videoList).Error
	return videoList, err
}

func (s *VideoDao) Feeds(lastTime int64) ([]*Video, error) {
	var videoList []*Video
	err := s.db.Where("create_time <= ?", lastTime).Find(videoList).Error
	return videoList, err
}
