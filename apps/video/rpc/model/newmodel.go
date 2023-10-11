package model

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

const (
	VideoInfoPrefix    = "tiktok:video:info:"    // +vid, hash	(info)
	VideoCountPrefix   = "tiktok:video:count:"   // +vid, hash	(count)
	VideoFeedPrefix    = "tiktok:video:feed::"   // +nil zset	(vid, timestamp)
	VideoPublishPrefix = "tiktok:video:publish:" // +uid set (vid)
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
