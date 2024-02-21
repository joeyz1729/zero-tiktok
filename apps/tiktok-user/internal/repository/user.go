package repository

import (
	"context"
	"errors"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
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
	DB  *gorm.DB
	RDB *redis.Client
	ES  *elasticsearch.TypedClient
}

var repo *Repo

func NewRepo(datasource, redisAddr string, esAddresses []string) (*Repo, error) {
	database, err := gorm.Open(mysql.Open(datasource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: esAddresses,
	})
	if err != nil {
		return nil, err
	}

	repo = &Repo{
		DB:  database,
		RDB: rdb,
		ES:  client,
	}

	return repo, nil
}

// Register 注册操作，事务添加用户信息，用户计数，不添加缓存
func (r *Repo) Register(userId int64, username, password string) error {
	user := db.User{ID: userId, Username: username, Password: password}
	return r.DB.Table(db.TableNameUser).Create(&user).Error
}

func (r *Repo) Login(username, password string) (userId int64, err error) {
	var user db.User
	err = r.DB.Table(db.TableNameUser).Where("username = ?", username).First(&user).Error
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
	err := r.DB.Table(db.TableNameUser).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) GetUserByName(username string) (*db.User, error) {
	var user db.User
	err := r.DB.Table(db.TableNameUser).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) CreateCount(userId int64, createdTime time.Time) error {
	var count = db.UserCount{
		ID:         userId,
		CreateTime: createdTime,
		UpdateTime: createdTime,
	}
	return r.DB.Table(db.TableNameUserCount).Create(&count).Error
}
