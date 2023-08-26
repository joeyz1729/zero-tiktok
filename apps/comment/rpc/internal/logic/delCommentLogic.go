package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DelCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCommentLogic {
	return &DelCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelCommentLogic) DelComment(in *model.DelCommentRequest) (*model.DelCommentResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.DelCommentResponse)
	var err error

	getStr := `select id from tiktok_comment.comment where video_id = ? and user_id = ? and comment_id = ?`
	_, err = l.svcCtx.CommentDB.Query(getStr, in.VideoId, in.UserId, in.CommentId)

	// 1. query failed
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorw("get comment record failed",
			logx.Field("err", err),
		)
		resp.Msg = err.Error()
		return resp, err
	}

	// 2. no comment record
	if err != nil && err == sqlx.ErrNotFound {
		resp.Msg = "no comment record"
		return resp, nil
	}

	// 3. query success
	delStr := `delete from tiktok_comment.comment where user_id = ? and video_id = ? and comment_id = ?`
	_, err = l.svcCtx.CommentDB.Exec(delStr, in.UserId, in.VideoId, in.CommentId)
	if err != nil {
		logx.Errorw("del comment record failed",
			logx.Field("err", err),
		)
		resp.Msg = err.Error()
		return resp, err
	}

	resp.Msg = "success"
	return resp, nil
}
