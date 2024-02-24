package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/cache"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/es"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
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
	ErrInvalidParams = errors.New("invalid params")
	ErrUserNotExist  = errors.New("tiktok-user not exist")
	ErrCacheMiss     = errors.New("cache miss")
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

func (r *Repo) GetUserDetail(userId int64) (*UserDetail, error) {
	var err error
	res, err := r.CacheGetUser(userId)
	if err != nil {
		logx.Errorw("[GetUserDetail] cache", logx.Field("err", err))
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	res, err = r.ESGetUser(userId)
	if err != nil {
		logx.Errorw("[GetUserDetail] es", logx.Field("err", err))
		return nil, err
	}
	go func() {
		err := r.CacheAddUser(res)
		if err != nil {
			logx.Errorw("[GetUserDetail] add cache", logx.Field("err", err))
		}
	}()
	return res, nil
}

func (r *Repo) DBCreateUser(userId int64, username, password string) error {
	user := db.User{ID: userId, Username: username, Password: password}
	return r.DB.Table(db.TableNameUser).Create(&user).Error
}

func (r *Repo) DBCreateCount(userId int64, createdTime time.Time) error {
	// todo 更新时间和创建时间需要吗，user和count的需要一致吗
	var count = db.UserCount{
		ID:         userId,
		CreateTime: createdTime,
		UpdateTime: createdTime,
	}
	return r.DB.Table(db.TableNameUserCount).Create(&count).Error
}

func (r *Repo) DBUpdateCount(userId, toUserId int64, incr int64) error {
	r.DB = r.DB.Table(db.TableNameUserCount)
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var userCount db.UserCount
		if err := tx.First(&userCount, userId).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", userId).
			Updates(map[string]interface{}{"follow_count": userCount.FollowCount + incr}).Error; err != nil {
			return err
		}

		var followedCount = db.UserCount{}
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

func (r *Repo) DBGetUserByName(username string) (*db.User, error) {
	var user db.User
	err := r.DB.Table(db.TableNameUser).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) DBGetUserById(userId int64) (*db.User, error) {
	var user db.User
	err := r.DB.Table(db.TableNameUser).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) CacheGetUser(userId int64) (*UserDetail, error) {
	b, err := r.RDB.Get(context.Background(), strconv.Itoa(int(userId))).Result()
	if err != nil {
		return nil, err
	}
	var detail UserDetail
	extra.RegisterFuzzyDecoders()
	err = jsoniter.Unmarshal([]byte(b), &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *Repo) CacheAddUser(detail *UserDetail) error {
	if detail == nil {
		return ErrInvalidParams
	}
	b, err := json.Marshal(detail)
	if err != nil {
		return err
	}
	key := cache.UserInfoPrefix + strconv.Itoa(int(detail.Id))
	_, err = r.RDB.Set(context.Background(), key, string(b), 24*time.Hour).Result()
	return err
}

func (r *Repo) ESGetUser(userId int64) (*UserDetail, error) {
	resp, err := repo.ES.Get(es.UserIndex, strconv.Itoa(int(userId))).Do(context.TODO())
	if err != nil {
		logx.Errorw("[ESGetUser]", logx.Field("err", err), logx.Field("userId", userId))
		return nil, err
	}
	b, err := resp.Source_.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var detail UserDetail
	extra.RegisterFuzzyDecoders()
	err = jsoniter.Unmarshal(b, &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *Repo) ESGetUserByName(username string) (*UserDetail, error) {
	resp, err := repo.ES.Search().
		Index(es.UserIndex).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Filter: []types.Query{
						{
							Term: map[string]types.TermQuery{
								"username": {Value: username},
							},
						},
					},
				},
			},
		}).Do(context.TODO())

	if err != nil {
		return nil, err
	}
	if len(resp.Hits.Hits) != 1 {
		return nil, errors.New("invalid record count")
	}
	str, err := resp.Hits.Hits[0].Source_.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var detail UserDetail
	extra.RegisterFuzzyDecoders()
	err = jsoniter.Unmarshal(str, &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}
