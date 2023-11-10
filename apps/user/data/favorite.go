package data

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
)

// AddFavoriteRelation 点赞操作的时候，事务更新用户和作者的计数
func (r *Repo) AddFavoriteRelation(userId, authorId int64) (err error) {
	// 更新数据库， 然后删除缓存
	var sqlStr1, sqlStr2 string
	sqlStr1 = `update tiktok_user.user_count set favorite_count = favorite_count + 1 where user_id = ? limit 1`
	sqlStr2 = `update tiktok_user.user_count set total_favorited = total_favorited + 1 where user_id = ? limit 1`
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		logx.Error("db transaction begin ", err)
		return err
	}
	// 出现错误返回的时候回滚
	defer tx.Rollback()
	if _, err = tx.Exec(sqlStr1, userId); err != nil {
		return err
	}
	if _, err = tx.Exec(sqlStr2, authorId); err != nil {
		return err
	}
	// 提交失败返回
	if err = tx.Commit(); err != nil {
		logx.Error("db transaction commit ", err)
		return err
	}

	return r.DelCountCache(userId, authorId)
}

// DelFavoriteRelation 取消点赞操作的时候，事务更新用户和作者的计数
func (r *Repo) DelFavoriteRelation(userId, authorId int64, actionType bool) (err error) {
	// 更新数据库， 然后删除缓存
	var sqlStr1, sqlStr2 string
	sqlStr1 = `update tiktok_user.user_count set favorite_count = favorite_count - 1 where user_id = ? limit 1`
	sqlStr2 = `update tiktok_user.user_count set total_favorited = total_favorited - 1 where user_id = ? limit 1`

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

	return r.DelCountCache(userId, authorId)
}
