package logic

import (
	"context"
	"encoding/json"

	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/kmq"
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

	kafkaData := kmq.KafkaData{
		ActionType: false,
		UserId:     in.UserId,
		VideoId:    in.VideoId,
		CommentId:  in.CommentId,
	}
	kafkaBytes, err := json.Marshal(kafkaData)
	if err != nil {
		logx.Errorw("encode kafka data failed",
			logx.Field("err", err))
		return resp, err
	}
	//kq.NewPusher()
	if err = l.svcCtx.KafkaPusher.Push(string(kafkaBytes)); err != nil {
		logx.Errorw("push add comment to kafka failed",
			logx.Field("err", err))
		return resp, err
	}

	// 3. query success

	resp.Msg = "push kafka mq success"
	return resp, nil
}
