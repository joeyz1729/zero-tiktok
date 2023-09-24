package db

import (
	"context"
	"database/sql"
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
	sqlStr = `select follower from tiktok_follow.follow_count where user_id = ?`
	err = tx.Get(&followerCount, sqlStr, tid)
	followedCount++
	followerCount++
	sqlStr = `update tiktok_follow.follow_count set followed = ? where user_id = ?`
	_, err = tx.Exec(sqlStr, followedCount, uid)
	sqlStr = `update tiktok_follow.follow_count set follower = ? where user_id = ?`
	_, err = tx.Exec(sqlStr, followerCount, tid)
	err = tx.Commit()
	if err != nil {
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
