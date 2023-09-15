package logic

import (
	"context"
	"encoding/json"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/kmq"
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

	// push kafka mq
	kafkaData := kmq.KafkaData{
		ActionType:  true,
		UserId:      in.UserId,
		VideoId:     in.VideoId,
		CommentId:   int64(cid),
		CommentText: in.CommentText,
	}
	logx.Info(kafkaData)
	kafkaBytes, err := json.Marshal(kafkaData)
	if err != nil {
		logx.Errorw("encode kafka data failed",
			logx.Field("err", err))
		return resp, err
	}

	if err = l.svcCtx.KafkaPusher.Push(string(kafkaBytes)); err != nil {
		logx.Errorw("push add comment to kafka failed",
			logx.Field("err", err))
		return resp, err
	}
	logx.Info("push add comment success")

	resp.Msg = "push kafka mq success"
	return resp, nil
}
