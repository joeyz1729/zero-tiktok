package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ FollowedModel = (*customFollowedModel)(nil)

type (
	// FollowedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowedModel.
	FollowedModel interface {
		followedModel
	}

	customFollowedModel struct {
		*defaultFollowedModel
	}
)

// NewFollowedModel returns a model for the database table.
func NewFollowedModel(conn sqlx.SqlConn) FollowedModel {
	return &customFollowedModel{
		defaultFollowedModel: newFollowedModel(conn),
	}
}
