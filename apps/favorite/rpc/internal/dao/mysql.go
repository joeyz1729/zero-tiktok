package dao

import (
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	FavoriteTable = "favorites"
)

type DBImpl struct {
	*gorm.DB
}

func NewDBImpl(dsn string) (*DBImpl, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		return nil, err
	}
	return &DBImpl{db}, nil
}

type FavorDB interface {
	IsFavoriteVideo(userId, videoId int64) (bool, error)
	IsFavoriteRecordExist(userId, videoId int64) (bool, error)
	CreateFavoriteRecord(favoriteInfo *model.Favorite) error
	DeleteFavoriteRecord(favoriteInfo *model.Favorite) error
}

var _ FavorDB = (*DBImpl)(nil)

// IsFavoriteRecordExist 更新前先检查是否存在，查询数据库时不更新缓存
func (r *DBImpl) IsFavoriteRecordExist(userId, videoId int64) (bool, error) {
	var favor model.Favorite
	err := r.Table(FavoriteTable).Where("user_id = ? and video_id = ?", userId, videoId).Take(&favor).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *DBImpl) CreateFavoriteRecord(favor *model.Favorite) error {
	return r.Table(FavoriteTable).Create(favor).Error
}

func (r *DBImpl) DeleteFavoriteRecord(favor *model.Favorite) error {
	return r.Table(FavoriteTable).Where("user_id = ? and video_id = ?", favor.UserId, favor.VideoId).Delete(&favor).Error

}

// IsFavoriteVideo 查询关系
func (r *DBImpl) IsFavoriteVideo(userId, videoId int64) (bool, error) {

	return false, nil
}
