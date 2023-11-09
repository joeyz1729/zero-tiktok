package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func (r *Repo) UpdateFavorTx(userId, authorId int64, actionType bool) (err error) {
	// 更新数据库， 然后删除缓存
	var sqlStr1, sqlStr2 string
	if actionType {
		sqlStr1 = `update tiktok_user.user_count set favorite_count = favorite_count + 1 where user_id = ? limit 1`
		sqlStr2 = `update tiktok_user.user_count set total_favorited = total_favorited + 1 where user_id = ? limit 1`
	} else {
		sqlStr1 = `update tiktok_user.user_count set favorite_count = favorite_count - 1 where user_id = ? limit 1`
		sqlStr2 = `update tiktok_user.user_count set total_favorited = total_favorited - 1 where user_id = ? limit 1`
	}
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		logx.Error("db transaction begin ", err)
		return err
	}
	defer tx.Rollback()
	if _, err = tx.Exec(sqlStr1, userId); err != nil {
		return err
	}
	if _, err = tx.Exec(sqlStr2, authorId); err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		logx.Error("db transaction commit ", err)
		return err
	}

	// delete cache
	uidStr := strconv.FormatInt(userId, 10)
	aidStr := strconv.FormatInt(authorId, 10)
	pipeline := r.rdb.Pipeline()
	r.rdb.Del(context.Background(), UserCountPrefix+uidStr)
	r.rdb.Del(context.Background(), UserCountPrefix+aidStr)
	_, err = pipeline.Exec(context.Background())
	return err
}

func (r *Repo) UpdateFavoriteCount(userId int64, actionType bool) (err error) {
	// 更新数据库， 然后删除缓存
	var sqlStr string
	if actionType {
		sqlStr = `update tiktok_user.user_count set favorite_count = favorite_count + 1 where user_id = ? limit 1`
	} else {
		sqlStr = `update tiktok_user.user_count set favorite_count = favorite_count - 1 where user_id = ? limit 1`
	}
	_, err = r.db.Exec(sqlStr, userId)
	if err != nil {

		fmt.Println("db update favorite ", err)
		return err
	}

	// delete cache
	uidStr := strconv.FormatInt(userId, 10)
	_, err = r.rdb.Del(context.Background(), UserCountPrefix+uidStr).Result()
	// todo, 放入消息队列异步删除
	return err
}

func (r *Repo) UpdateTotalFavorited(authorId int64, actionType bool) (err error) {
	var sqlStr string
	if actionType {
		sqlStr = `update tiktok_user.user_count set total_favorited = total_favorited + 1 where user_id = ? limit 1`
	} else {
		sqlStr = `update tiktok_user.user_count set total_favorited = total_favorited - 1 where user_id = ? limit 1`
	}
	_, err = r.db.Exec(sqlStr, authorId)
	if err != nil {
		logx.Error("db update favorite ", err)
		return err
	}

	// delete cache
	uidStr := strconv.FormatInt(authorId, 10)
	_, err = r.rdb.Del(context.Background(), UserCountPrefix+uidStr).Result()
	// todo, 放入消息队列异步删除
	return err
}
