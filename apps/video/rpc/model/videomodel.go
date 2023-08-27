package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VideoModel = (*customVideoModel)(nil)

type (
	// VideoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideoModel.
	VideoModel interface {
		videoModel
		FindVideosByUserId(ctx context.Context, userId int64) ([]*Video, error)
	}

	customVideoModel struct {
		*defaultVideoModel
	}
)

// NewVideoModel returns a model for the database table.
func NewVideoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) VideoModel {
	return &customVideoModel{
		defaultVideoModel: newVideoModel(conn, c, opts...),
	}
}

func (m *customVideoModel) FindVideosByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	videoIdKey := fmt.Sprintf("%s%v", cacheVideoIdPrefix, userId)
	var resp []*Video
	err := m.QueryRowsNoCacheCtx(ctx, &resp, videoIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from video where `user_id` = ?", videoRows)
		return conn.QueryRowsCtx(ctx, &resp, query, userId)
	})
	logx.Errorw("find videos by user id failed",
		logx.Field("err", err),
	)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

type VideoDetail struct {
	VideoId  int64  `db:"video_id"`
	AuthorId int64  `db:"author_id"`
	Title    string `db:"title"`
	PlayUrl  string `db:"play_url"`
	CoverUrl string `db:"cover_url"`
}
