package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	sqlz "github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"time"
)

func (r *Repo) Register(userId int64, username, password string) error {
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlStr := `insert into tiktok_user.user(user_id, username, password, avatar, background_image, signature) value(?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(sqlStr, userId, username, password, defaultAvatar, defaultBackgroundImage, defaultSignature)
	if err != nil {
		return err
	}
	sqlStr = `insert into tiktok_user.user_count(user_id) value(?)`
	_, err = tx.Exec(sqlStr, userId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	uidStr := strconv.FormatInt(userId, 10)
	_, err = r.rdb.HSet(context.Background(), UserInfoPrefix+uidStr, FieldUsername, username).Result()
	if err != nil {
		return err
	}
	_, err = r.rdb.Expire(context.Background(), UserInfoPrefix+uidStr, time.Hour*168).Result()

	_, err = r.rdb.HSet(context.Background(), UserCountPrefix+uidStr, countMap).Result()
	if err != nil {
		return err
	}
	return err
}

type UserDetail struct {
	Id       int64  `db:"id"`
	UserId   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`

	FollowedCount int64 `db:"followed_count" json:"followedcount,string"` // 关注总数
	FollowerCount int64 `db:"follower_count" json:"followercount,string"` // 粉丝总数

	TotalFavorited int64     `db:"total_favorited" json:"totalfavorited,string"` //获赞数量
	WorkCount      int64     `db:"work_count" json:"workcount,string"`           //作品数量
	FavoriteCount  int64     `db:"favorite_count" json:"favoritecount,string"`   //点赞数量
	CreateTime     time.Time `db:"create_time"`
	UpdateTime     time.Time `db:"update_time"`
}

func (r *Repo) GetUserInfo(userId int64) (user *UserDetail, err error) {
	user = new(UserDetail)
	user.UserId = userId
	username, err := r.GetUsername(userId)
	if err != nil {
		logx.Error("get username failed", err)
		return nil, err
	}
	user.Username = username

	err = r.GetCount(userId, user)
	if err != nil {
		logx.Error("get count failed", err)
		return nil, err
	}
	return user, nil
}

func (r *Repo) GetUsername(userId int64) (username string, err error) {
	uidStr := strconv.FormatInt(userId, 10)
	username, err = r.rdb.HGet(context.Background(), UserInfoPrefix+uidStr, FieldUsername).Result()
	if err == nil {
		_, err = r.rdb.Expire(context.Background(), UserInfoPrefix+uidStr, time.Hour*168).Result()
		return username, nil
	}

	sqlStr := `select (username) from tiktok_user.user where user_id = ? limit 1`
	err = r.db.Get(&username, sqlStr, userId)
	//TODO
	if err != nil {
		return "", err
	}
	return username, nil
}

func (r *Repo) GetCount(userId int64, user *UserDetail) (err error) {
	uidStr := strconv.FormatInt(userId, 10)
	num, err := r.rdb.Exists(context.Background(), UserCountPrefix+uidStr).Result()
	if err != nil {
		logx.Error("check redis exists failed")
		return err
	}
	if num != 0 {
		cm, err := r.rdb.HGetAll(context.Background(), UserCountPrefix+uidStr).Result()
		if err == nil {
			b, err := json.Marshal(cm)
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, user)
			return err
		}
	}

	sqlStr := `select * from tiktok_user.user_count where user_id = ?`
	err = r.db.Get(user, sqlStr, userId)
	if err != nil {
		return
	}
	return r.AddCountCache(userId, user)

}

func (r *Repo) CheckUserValid(username string) (ok bool, err error) {
	var id int64
	sqlStr := `select id from tiktok_user.user where username = ? limit 1`
	err = r.db.Get(&id, sqlStr, username)
	if err != nil && err != sqlz.ErrNotFound {
		return false, err
	}
	return err != sqlz.ErrNotFound, nil
}

func (r *Repo) CheckLogin(username, password string) (userId int64, err error) {
	var user UserDetail
	sqlStr := `select user_id, password from tiktok_user.user where username = ? limit 1`
	err = r.db.Get(&user, sqlStr, username)
	if err != nil && err != sqlc.ErrNotFound {
		return 0, err
	}
	if err == sqlc.ErrNotFound {
		return 0, ErrUserNotExist
	}
	if password != user.Password {
		return 0, ErrInvalidPassword
	}
	return user.UserId, nil
}

func (r *Repo) AddCountCache(uid int64, user *UserDetail) (err error) {
	pipeline := r.rdb.TxPipeline()
	uidStr := strconv.FormatInt(uid, 10)
	pipeline.HSet(context.Background(), UserCountPrefix+uidStr, FieldFollowerCount, user.FollowerCount)
	pipeline.HSet(context.Background(), UserCountPrefix+uidStr, FieldFollowedCount, user.FollowedCount)
	pipeline.HSet(context.Background(), UserCountPrefix+uidStr, FieldTotalFavorited, user.TotalFavorited)
	pipeline.HSet(context.Background(), UserCountPrefix+uidStr, FieldWorkCount, user.WorkCount)
	pipeline.HSet(context.Background(), UserCountPrefix+uidStr, FieldFavoriteCount, user.FavoriteCount)
	_, err = pipeline.Exec(context.Background())
	return err
}
