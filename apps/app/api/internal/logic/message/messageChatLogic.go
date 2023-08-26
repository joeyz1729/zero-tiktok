package message

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/message/message"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"
	"time"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageChatLogic {
	return &MessageChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageChatLogic) MessageChat(req *types.MessageChatRequest) (resp *types.MessageChatResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.MessageChatResponse)
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("jwt parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusUnauthorized
		resp.StatusMsg = "invalid token"
		resp.MessageList = []types.Message{}
		return resp, nil
	}
	uid := claims.UserId

	listRes := new(message.ListResponse)
	listRes, err = l.svcCtx.MessageRpc.List(l.ctx, &message.ListRequest{
		UserId:     uid,
		ToUserId:   req.ToUserId,
		PreMsgTime: time.Now().Add(time.Duration(req.PreMsgTime) * time.Second).String(),
	})
	if err != nil {
		logx.Errorw("message list rpc failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "message rpc failed"
		resp.MessageList = []types.Message{}
		return resp, err
	}

	if listRes.MessageLen == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "empty list"
		resp.MessageList = []types.Message{}
		return resp, err
	}

	resp.MessageList = make([]types.Message, listRes.MessageLen)
	for i, m := range listRes.MessageList {
		msg := types.Message{
			Id:         m.Id,
			FromUserId: m.FromUserId,
			ToUserId:   m.ToUserId,
			Content:    m.Content,
			CreateTime: m.CreateTime,
		}
		resp.MessageList[i] = msg
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return
}
