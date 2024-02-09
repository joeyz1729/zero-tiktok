package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/message/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/message/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *model.ListRequest) (*model.ListResponse, error) {

	resp := new(model.ListResponse)
	var err error
	// todo: create_time > pre_msg_time
	var ml []*model.Message
	sqlStr := `select * from tiktok_message.message where user_id = ? and to_user_id = ? and create_time > ?`
	err = l.svcCtx.MessageDB.Select(&ml, sqlStr, in.UserId, in.ToUserId, "2023-08-27 00:35:33")
	if err != nil {
		logx.Errorw("mysql get message list failed",
			logx.Field("err", err),
		)
		resp.MessageLen = 0
		resp.MessageList = []*model.MessageInfo{}
		return resp, err
	}
	if len(ml) == 0 {
		resp.MessageLen = 0
		resp.MessageList = []*model.MessageInfo{}
		return resp, nil
	}

	resp.MessageLen = int64(len(ml))
	resp.MessageList = make([]*model.MessageInfo, len(ml))
	for i, m := range ml {
		mi := &model.MessageInfo{
			Id:         m.Id,
			FromUserId: m.UserId,
			ToUserId:   m.ToUserId,
			Content:    m.Content,
			CreateTime: m.CreateTime.String(),
		}
		resp.MessageList[i] = mi
	}
	return resp, nil
}
