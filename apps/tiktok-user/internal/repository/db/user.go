package db

import (
	"gorm.io/gorm"
	"time"
)

const TableNameUser = "user"

// User 用户表
type User struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username        string    `gorm:"column:username;not null" json:"username"`
	Password        string    `gorm:"column:password;not null" json:"password"`
	Avatar          string    `gorm:"column:avatar;not null" json:"avatar"`
	BackgroundImage string    `gorm:"column:background_image;not null" json:"background_image"`
	Signature       string    `gorm:"column:signature;not null" json:"signature"`
	CreateTime      time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

func CreateUser(userId int64, username, password string, DB *gorm.DB) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		var user = User{
			ID:       userId,
			Username: username,
			Password: password,
		}
		if err := tx.Table(TableNameUser).Create(&user).Error; err != nil {
			return err
		}
		var userCount = UserCount{
			ID: userId,
		}
		if err := tx.Table(TableNameUserCount).Create(&userCount).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetUserByName(username string, DB *gorm.DB) (*User, error) {
	var user User
	err := DB.Table(TableNameUser).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserById(userId int64, DB *gorm.DB) (*User, error) {
	var user User
	err := DB.Table(TableNameUser).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateRelationCount(userId, toUserId int64, incr int64, DB *gorm.DB) error {
	err := DB.Table(TableNameUserCount).Transaction(func(tx *gorm.DB) error {
		var userCount UserCount
		if err := tx.First(&userCount, userId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", userId).
			Updates(map[string]interface{}{"follow_count": userCount.FollowCount + incr}).Error; err != nil {
			return err
		}

		var followedCount = UserCount{}
		if err := tx.First(&followedCount, toUserId).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", toUserId).
			Updates(map[string]interface{}{"follower_count": followedCount.FollowerCount + incr}).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
	return err
}

func UpdateWorkCount(userId int64, incr int64, DB *gorm.DB) error {
	err := DB.Table(TableNameUserCount).Transaction(func(tx *gorm.DB) error {
		var userCount UserCount
		if err := tx.First(&userCount, userId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", userId).
			Updates(map[string]interface{}{"work_count": userCount.WorkCount + incr}).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return err
}

func UpdateFavorCount(userId, authorId int64, incr int64, DB *gorm.DB) error {
	DB = DB.Table(TableNameUserCount)
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 更新用户点赞数
		var userCount UserCount
		if err := tx.First(&userCount, userId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", userId).
			Updates(map[string]interface{}{"favorite_count": userCount.FavoriteCount + incr}).Error; err != nil {
			return err
		}
		// 更新作者获赞数
		var authorCount = UserCount{}
		if err := tx.First(&authorCount, authorId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", authorId).
			Updates(map[string]interface{}{"total_favorited": authorCount.TotalFavorited + incr}).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return err
}

func UpdateFavoriteCount(userId int64, incr int64, DB *gorm.DB) error {
	DB = DB.Table(TableNameUserCount)
	err := DB.Transaction(func(tx *gorm.DB) error {
		var userCount UserCount
		if err := tx.First(&userCount, userId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", userId).
			Updates(map[string]interface{}{"favorite_count": userCount.FavoriteCount + incr}).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return err
}

func UpdateTotalFavorited(userId int64, incr int64, DB *gorm.DB) error {
	DB = DB.Table(TableNameUserCount)
	err := DB.Transaction(func(tx *gorm.DB) error {
		var userCount = UserCount{}
		if err := tx.First(&userCount, userId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", userId).
			Updates(map[string]interface{}{"total_favorited": userCount.TotalFavorited + incr}).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return err
}
