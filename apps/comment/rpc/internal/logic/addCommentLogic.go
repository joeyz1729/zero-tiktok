package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/model"
	"github.com/YiZou89/zero-tiktok/pkg/snowflake"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCommentLogic) AddComment(in *model.AddCommentRequest) (*model.AddCommentResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.AddCommentResponse)
	var err error

	cid, err := snowflake.GenID()
	if err != nil {
		logx.Errorw("snowflake gen comment id failed",
			logx.Field("err", err),
		)
		resp.Msg = err.Error()
		return resp, err
	}

	sqlStr := `insert into tiktok_comment.comment(video_id, user_id, comment_id, content) value(?, ?, ?, ?)`
	_, err = l.svcCtx.CommentDB.Exec(sqlStr, in.VideoId, in.UserId, cid, in.CommentText)
	if err != nil {
		logx.Errorw("mysql insert comment record failed",
			logx.Field("err", err),
		)
		resp.Msg = err.Error()
		return resp, err
	}

	resp.Msg = "success"
	return resp, nil
}
