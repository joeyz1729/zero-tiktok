package db

import (
	"time"
)

const TableNameRelationCount = "relation_count"

// RelationCount 用户关系计数表
type RelationCount struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Followed   int64     `gorm:"column:followed;not null" json:"followed"`
	Follower   int64     `gorm:"column:follower;not null" json:"follower"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName RelationCount's table name
func (*RelationCount) TableName() string {
	return TableNameRelationCount
}

func (fd *FollowDB) DecrCount(userId int64, toUserId int64) error {
	return nil
}

func (fd *FollowDB) IncrCount(userId int64, toUserId int64) error {
	return nil
}
