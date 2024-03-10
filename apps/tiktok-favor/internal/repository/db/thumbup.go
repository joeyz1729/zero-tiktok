package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

const TableNameThumbup = "thumbup"

type Thumbup struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID     int64     `gorm:"column:user_id;not null" json:"user_id"`
	VideoID    int64     `gorm:"column:video_id;not null" json:"video_id"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

func (*Thumbup) TableName() string {
	return TableNameThumbup
}

type GormImpl struct {
	*gorm.DB
}

func NewGormImpl(dsn string) (*GormImpl, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &GormImpl{db}, nil
}

func (r *GormImpl) IsThumbup(userId, videoId int64) (bool, error) {
	var thumbup Thumbup
	err := r.Table(TableNameThumbup).Where("user_id = ? and video_id = ?", userId, videoId).Take(&thumbup).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *GormImpl) Add(userId int64, videoId int64) error {
	thumbup := Thumbup{
		UserID:  userId,
		VideoID: videoId,
	}
	return r.Table(TableNameThumbup).Create(&thumbup).Error
}

func (r *GormImpl) Delete(userId int64, videoId int64) error {
	thumbup := Thumbup{
		UserID:  userId,
		VideoID: videoId,
	}
	return r.Table(TableNameThumbup).Where("user_id = ? and video_id = ?", userId, videoId).Delete(&thumbup).Error
}

func (r *GormImpl) GetThumbupListByUserId(userId int64) ([]int64, error) {
	var list []*Thumbup
	err := r.DB.Table(TableNameThumbup).Select("video_id").
		Where("user_id = ?", userId).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(list))
	for i, f := range list {
		ids[i] = f.VideoID
	}
	return ids, nil
}

func (r *GormImpl) GetUserThumbupCount(userId int64) (int64, error) {
	var count int64
	err := r.DB.Table(TableNameThumbup).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormImpl) GetVideoThumbupCount(videoId int64) (int64, error) {
	var count int64
	err := r.DB.Table(TableNameThumbup).Where("video_id = ?", videoId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
