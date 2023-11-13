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
	resp := new(model.AddCommentResponse)
	var err error
	// 生成commentId
	cid, err := snowflake.GenID()
	if err != nil {
		logx.Error("snowflake gen comment id ", err)
		return nil, err
	}

	// 包装消息并序列化
	kafkaData := kmq.KafkaData{
		ActionType:  true,
		UserId:      in.UserId,
		VideoId:     in.VideoId,
		CommentId:   int64(cid),
		CommentText: in.CommentText,
	}
	logx.Info("prepare to send msg ", kafkaData)
	kafkaBytes, err := json.Marshal(kafkaData)
	if err != nil {
		logx.Error("encode message ", err)
		return nil, err
	}

	// 发送消息，go-zero包装的kafka-go好像只能设置时间作为key
	if err = l.svcCtx.KafkaPusher.Push(string(kafkaBytes)); err != nil {
		logx.Error("push AddComment ", err)
		return resp, err
	}

	logx.Info("push message success")
	return resp, nil
}
