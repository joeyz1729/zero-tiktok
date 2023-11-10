package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrEmptySet = errors.New("empty set")
)

type FollowDB struct {
	*sqlx.DB
}

func NewFollowDB(db *sqlx.DB) *FollowDB {
	return &FollowDB{
		db,
	}
}

func (fd *FollowDB) CheckRelation(ctx context.Context, uid, tid int64) (ok bool, err error) {
	sqlStr := `select id from tiktok_follow.follow where user_id = ? and follow_id = ? limit 1`
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

func (fb *FollowDB) GetFollowerIds(ctx context.Context, uid int64) (ids []int64, err error) {
	sqlStr := `select user_id from tiktok_follow.follow where follow_id = ? `
	err = fb.Select(&ids, sqlStr, uid)
	if err != nil {
		logx.Errorw("mysql query failed",
			logx.Field("err", err),
		)
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrEmptySet
	}
	return ids, nil
}

func (fb *FollowDB) GetFollowedIds(ctx context.Context, uid int64) (ids []int64, err error) {
	sqlStr := `select follow_id from tiktok_follow.follow where user_id = ? `
	err = fb.Select(&ids, sqlStr, uid)
	if err != nil {
		logx.Errorw("mysql query failed",
			logx.Field("err", err),
		)
		return nil, err
	}
	if len(ids) == 0 {
		return nil, ErrEmptySet
	}
	return ids, nil
}
