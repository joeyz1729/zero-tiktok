package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

var (
	ErrEmptySet = errors.New("empty set")
)

const TableNameRelation = "relation"

// Relation 关注表
type Relation struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID     int64     `gorm:"column:user_id;not null" json:"user_id"`
	FollowedID int64     `gorm:"column:followed_id;not null" json:"followed_id"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName Relation's table name
func (*Relation) TableName() string {
	return TableNameRelation
}

type FollowDB struct {
	*gorm.DB
}

func NewFollowDB(db *gorm.DB) *FollowDB {
	return &FollowDB{
		db,
	}
}

func (fd *FollowDB) IsFollow(userId, toUserId int64) (bool, error) {
	var relation Relation
	err := fd.Table(TableNameRelation).Where("user_id = ? and followed_id = ?", userId, toUserId).Take(&relation).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fd *FollowDB) Add(userId int64, followedId int64) error {
	relation := Relation{
		UserID:     userId,
		FollowedID: followedId,
	}
	return fd.Table(TableNameRelation).Create(&relation).Error
}

func (fd *FollowDB) Delete(userId int64, followedId int64) error {
	relation := Relation{
		UserID:     userId,
		FollowedID: followedId,
	}
	return fd.Table(TableNameRelation).Where("user_id = ? and followed_id = ?", userId, followedId).Delete(&relation).Error
}

func (fd *FollowDB) GetFollowerIds(ctx context.Context, uid int64) (ids []int64, err error) {
	return nil, nil
}

func (fd *FollowDB) GetFollowedIds(ctx context.Context, uid int64) (ids []int64, err error) {
	return nil, nil
}
