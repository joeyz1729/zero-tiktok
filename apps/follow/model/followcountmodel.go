package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ FollowCountModel = (*customFollowCountModel)(nil)

type (
	// FollowCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowCountModel.
	FollowCountModel interface {
		followCountModel
	}

	customFollowCountModel struct {
		*defaultFollowCountModel
	}
)

// NewFollowCountModel returns a data for the database table.
func NewFollowCountModel(conn sqlx.SqlConn) FollowCountModel {
	return &customFollowCountModel{
		defaultFollowCountModel: newFollowCountModel(conn),
	}
}
