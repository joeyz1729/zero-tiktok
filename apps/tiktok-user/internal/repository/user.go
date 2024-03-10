package repository

import (
	"context"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/cache"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/dto"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/es"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

const (
	defaultAvatar          = "default avatar"
	defaultBackgroundImage = "default background image"
	defaultSignature       = "default signature"
)

type Repo struct {
	DB  *gorm.DB
	RDB *redis.Client
	ES  *elasticsearch.TypedClient
}

var repo *Repo

func NewRepo(c config.RepoConfig) (*Repo, error) {
	database, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: c.RedisAddr,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	crt, err := os.ReadFile(c.EsCACert)
	if err != nil {
		return nil, err
	}

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: c.EsAddresses,
		Username:  c.EsUsername,
		Password:  c.EsPassword,
		CACert:    crt,
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

func (r *Repo) GetUserDetail(userId int64) (*dto.User, error) {
	var err error
	res, err := cache.GetUser(userId, repo.RDB)
	if err != nil && err != redis.Nil {
		// 查询错误
		logx.Errorw("[GetUserDetail] cache", logx.Field("err", err))
		return nil, err
	}
	// 查询成功
	if res != nil {
		return res, nil
	}
	// redis不存在
	res, err = es.GetUser(userId, repo.ES)
	if err != nil {
		logx.Errorw("[GetUserDetail] es", logx.Field("err", err))
		return nil, err
	}
	go func() {
		err := cache.AddUser(res, repo.RDB)
		if err != nil {
			logx.Errorw("[GetUserDetail] add cache", logx.Field("err", err))
		}
	}()
	return res, nil
}

func (r *Repo) GetUserByName(username string) (*dto.User, error) {
	return es.GetUserByName(username, r.ES)

}

func (r *Repo) CreateUser(userId int64, username, password string) error {
	return db.CreateUser(userId, username, password, r.DB)
}
