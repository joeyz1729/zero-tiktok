package repository

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultAvatar          = "default avatar"
	defaultBackgroundImage = "default background image"
	defaultSignature       = "default signature"
)

var (
	ErrUserNotExist    = errors.New("tiktok-user not exist")
	ErrCacheMiss       = errors.New("cache miss")
	ErrInvalidPassword = errors.New("invalid password")
)

type Repo struct {
	db  *gorm.DB
	rdb *redis.Client
}

var repo *Repo

func NewRepo(datasource, redisAddr string) *Repo {
	database, err := gorm.Open(mysql.Open(datasource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	rdb.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	repo = &Repo{
		db:  database,
		rdb: rdb,
	}
	return repo
}

// Register 注册操作，事务添加用户信息，用户计数，不添加缓存
func (r *Repo) Register(userId int64, username, password string) error {
	user := db.User{ID: userId, Username: username, Password: password}
	return r.db.Create(&user).Error
}

func (r *Repo) Login(username, password string) (userId int64, err error) {
	var user db.User
	err = r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return 0, err
	}
	if password != user.Password {
		return 0, ErrInvalidPassword
	}
	return user.ID, nil
}

func (r *Repo) GetUserById(userId int64) (*db.User, error) {
	var user db.User
	err := r.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) GetUserByName(username string) (*db.User, error) {
	var user db.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
