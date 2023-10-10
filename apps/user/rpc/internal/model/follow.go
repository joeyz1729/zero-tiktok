package model

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func (r *Repo) AddFollow(userId, toUserId int64) (err error) {
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlStr := `update tiktok_user.user_count set followed_count = followed_count + 1 where user_id = ?`
	if _, err = tx.Exec(sqlStr, userId); err != nil {
		return err
	}
	sqlStr = `update tiktok_user.user_count set follower_count = follower_count + 1 where user_id = ?`
	if _, err = tx.Exec(sqlStr, toUserId); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		logx.Errorw("mysql transaction commit failed", logx.Field("err", err))
		return err
	}

	uidStr := strconv.FormatInt(userId, 10)
	tidStr := strconv.FormatInt(toUserId, 10)
	pipeline := r.rdb.TxPipeline()
	pipeline.HIncrBy(context.Background(), UserCountPrefix+uidStr, FieldFollowedCount, 1)
	pipeline.HIncrBy(context.Background(), UserCountPrefix+tidStr, FieldFollowerCount, 1)
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		logx.Errorw("redis pipeline exec failed",
			logx.Field("err", err))
		return err
	}

	return nil
}

func (r *Repo) DelFollow(userId, toUserId int64) (err error) {
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlStr := `update tiktok_user.user_count set followed_count = followed_count - 1 where user_id = ?`
	if _, err = tx.Exec(sqlStr, userId); err != nil {
		return err
	}
	sqlStr = `update tiktok_user.user_count set follower_count = follower_count - 1 where user_id = ?`
	if _, err = tx.Exec(sqlStr, toUserId); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		logx.Errorw("mysql transaction commit failed", logx.Field("err", err))
		return err
	}

	uidStr := strconv.FormatInt(userId, 10)
	tidStr := strconv.FormatInt(toUserId, 10)
	pipeline := r.rdb.TxPipeline()
	pipeline.HIncrBy(context.Background(), UserCountPrefix+uidStr, FieldFollowedCount, -1)
	pipeline.HIncrBy(context.Background(), UserCountPrefix+tidStr, FieldFollowerCount, -1)
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		logx.Errorw("redis pipeline exec failed",
			logx.Field("err", err))
		return err
	}

	return nil
}
