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

func (fd *FollowDB) AddRelation(ctx context.Context, uid, tid int64) (err error) {
	tx, err := fd.DB.Beginx()
	if err != nil {
		logx.Errorw("[mysql] begin transaction failed",
			logx.Field("err", err))
		return err
	}

	sqlStr := `insert into tiktok_follow.followed(user_id, followed_id) value(?, ?)`
	_, err = tx.Exec(sqlStr, uid, tid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logx.Errorw("add followed failed, and rollback failed",
				logx.Field("err", err))
			return rollbackErr
		}
		logx.Errorw("add followed failed, rollback success",
			logx.Field("err", err))
		return err
	}

	sqlStr = `insert into tiktok_follow.follower(user_id, follower_id) value(?, ?)`
	_, err = tx.Exec(sqlStr, tid, uid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logx.Errorw("add follower failed, and rollback failed")
			return rollbackErr
		}
		logx.Errorw("add follower failed, rollback success")
		return err
	}

	var followedCount, followerCount int32
	sqlStr = `select followed from tiktok_follow.follow_count where user_id = ?`
	err = tx.Get(&followedCount, sqlStr, uid)
	sqlStr = `select follower from tiktok_follow.follow_count where user_id = ?`
	err = tx.Get(&followerCount, sqlStr, tid)

	sqlStr = `update tiktok_follow.follow_count set followed = ? where user_id = ?`
	_, err = tx.Exec(sqlStr, followedCount+1, uid)
	sqlStr = `update tiktok_follow.follow_count set follower = ? where user_id = ?`
	_, err = tx.Exec(sqlStr, followerCount+1, tid)

	return nil
}

func (fd *FollowDB) CheckRelation(ctx context.Context, uid, tid int64) (ok bool, err error) {
	sqlStr := `select id from tiktok_follow.followed where user_id = ? and followed_id = ?`
	_, err = fd.DB.Query(sqlStr, uid, tid)
	// failed
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil

}
