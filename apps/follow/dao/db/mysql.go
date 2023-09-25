package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)
import "github.com/jmoiron/sqlx"

type FollowDB struct {
	*sqlx.DB
}

func NewFollowDB(db *sqlx.DB) *FollowDB {
	return &FollowDB{
		db,
	}
}

func (fd *FollowDB) AddRelation(ctx context.Context, uid, tid int64) (followedCount, followerCount int32, err error) {
	tx, err := fd.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		logx.Errorw("[mysql] begin transaction failed",
			logx.Field("err", err))
		return 0, 0, err
	}

	sqlStr := `insert into tiktok_follow.followed(user_id, followed_id) value(?, ?)`
	_, err = tx.Exec(sqlStr, uid, tid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logx.Errorw("add followed failed, and rollback failed",
				logx.Field("err", err))
			return 0, 0, rollbackErr
		}
		logx.Errorw("add followed failed, rollback success",
			logx.Field("err", err))
		return 0, 0, err
	}

	sqlStr = `insert into tiktok_follow.follower(user_id, follower_id) value(?, ?)`
	_, err = tx.Exec(sqlStr, tid, uid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logx.Errorw("add follower failed, and rollback failed")
			return 0, 0, rollbackErr
		}
		logx.Errorw("add follower failed, rollback success")
		return 0, 0, err
	}

	sqlStr = `select followed from tiktok_follow.follow_count where user_id = ?`
	err = tx.Get(&followedCount, sqlStr, uid)
	if err == sql.ErrNoRows {
		// 计数表中没有user的数据
		sqlStr = `insert into tiktok_follow.follow_count(user_id,followed) value(?, ?)`
		res, err := tx.Exec(sqlStr, uid, 1)
		fmt.Println(res, err)
	} else if err == nil {
		// 计数表中有数据，并且已经查到计数
		sqlStr = `update tiktok_follow.follow_count set followed = ? where user_id = ?`
		_, err = tx.Exec(sqlStr, followedCount+1, uid)
	}

	sqlStr = `select follower from tiktok_follow.follow_count where user_id = ?`
	err = tx.Get(&followerCount, sqlStr, tid)
	if err == sql.ErrNoRows {
		// 计数表中没有user的数据
		sqlStr = `insert into tiktok_follow.follow_count(user_id,follower) value(?, ?)`
		tx.Exec(sqlStr, tid, 1)
	} else if err == nil {
		// 计数表中有数据，并且已经查到计数
		sqlStr = `update tiktok_follow.follow_count set follower = ? where user_id = ?`
		tx.Exec(sqlStr, followerCount+1, tid)
	}
	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}
	return followedCount + 1, followerCount + 1, nil
}

func (fd *FollowDB) DelRelation(ctx context.Context, uid, tid int64) (followedCount, followerCount int32, err error) {
	tx, err := fd.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		logx.Errorw("[mysql] begin transaction failed",
			logx.Field("err", err))
		return 0, 0, err
	}

	// 查询uid的关注数，和tid的粉丝数
	sqlStr := `select followed from tiktok_follow.follow_count where user_id = ?`
	err = tx.Get(&followedCount, sqlStr, uid)
	if err != nil {
		// 不论是查询出错还是没有数据
		return 0, 0, err
	}
	// tid同理
	sqlStr = `select follower from tiktok_follow.follow_count where user_id = ?`
	err = tx.Get(&followerCount, sqlStr, tid)
	if err != nil {
		// 计数表中没有user的数据
		return 0, 0, err
	}
	// 查询成功，并记录数量

	// 删除关系
	sqlStr = `delete from tiktok_follow.followed where user_id = ? and followed_id = ?`
	_, err = tx.Exec(sqlStr, uid, tid)
	if err != nil {
		return 0, 0, err
	}
	sqlStr = `delete from tiktok_follow.follower where user_id = ? and follower_id = ?`
	_, err = tx.Exec(sqlStr, tid, uid)
	if err != nil {
		return 0, 0, err
	}

	// 修改计数
	followerCount--
	followedCount--
	sqlStr = `update tiktok_follow.follow_count set followed = ? where user_id = ?`
	_, err = tx.Exec(sqlStr, followedCount, uid)
	if err != nil {
		return 0, 0, err
	}
	sqlStr = `update tiktok_follow.follow_count set follower = ? where user_id = ?`
	_, err = tx.Exec(sqlStr, followerCount, tid)
	if err != nil {
		return 0, 0, err
	}
	err = tx.Commit()
	if err != nil {
		logx.Errorw("tx commit failed",
			logx.Field("err", err))
		return 0, 0, err
	}
	return followedCount, followerCount, nil
}

func (fd *FollowDB) CheckRelation(ctx context.Context, uid, tid int64) (ok bool, err error) {
	sqlStr := `select id from tiktok_follow.followed where user_id = ? and followed_id = ? limit 1`
	row := fd.DB.QueryRowx(sqlStr, uid, tid)
	// failed
	var id int64
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil

}
