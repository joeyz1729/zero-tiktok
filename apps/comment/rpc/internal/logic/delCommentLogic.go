package logic

import (
	"context"
	"encoding/json"

	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/kmq"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
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
	resp := new(model.DelCommentResponse)
	var err error
	// 1. 检查要删除的comment是否存在
	var count = 0
	getStr := `select count(*) from tiktok_comment.comment where video_id = ?  and comment_id = ? and user_id = ? limit 1`
	err = l.svcCtx.CommentDB.Get(&count, getStr, in.VideoId, in.CommentId, in.UserId)
	if err != nil {
		logx.Error("get comment record ", err)
		return nil, err
	}

	// 2. no comment record
	if count == 0 {
		logx.Error("no record")
		return nil, err
	}
	// 3. 添加到消息队列
	kafkaData := kmq.KafkaData{
		ActionType: false,
		UserId:     in.UserId,
		VideoId:    in.VideoId,
		CommentId:  in.CommentId,
	}
	logx.Info("kafka message ", kafkaData)
	kafkaBytes, err := json.Marshal(kafkaData)
	if err != nil {
		logx.Error("encode message ", err)
		return resp, err
	}

	// 暂时使用默认的pusher，kq.NewPusher()
	if err = l.svcCtx.KafkaPusher.Push(string(kafkaBytes)); err != nil {
		logx.Error("push DelComment ", err)
		return resp, err
	}

	// 3. query success
	return resp, nil
}
