package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ FollowerModel = (*customFollowerModel)(nil)

type (
	// FollowerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowerModel.
	FollowerModel interface {
		followerModel
	}

	customFollowerModel struct {
		*defaultFollowerModel
	}
)

// NewFollowerModel returns a model for the database table.
func NewFollowerModel(conn sqlx.SqlConn) FollowerModel {
	return &customFollowerModel{
		defaultFollowerModel: newFollowerModel(conn),
	}
}
