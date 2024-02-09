package message

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/message/message"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageActionLogic {
	return &MessageActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageActionLogic) MessageAction(req *types.MessageActionRequest) (resp *types.MessageActionResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.MessageActionResponse)
	if req.ActionType != int32(1) {
		logx.Errorw("invalid action type",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "invalid action type"
		return resp, nil
	}

	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("jwt parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusUnauthorized
		resp.StatusMsg = "invalid token"
		return resp, nil
	}

	uid := claims.UserId
	actionRes := new(message.ActionResponse)
	actionRes, err = l.svcCtx.MessageRpc.Action(l.ctx, &message.ActionRequest{
		UserId:   uid,
		ToUserId: req.ToUserId,
		Content:  req.Content,
	})
	if err != nil {
		logx.Errorw("message tiktok-user failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = actionRes.Msg
		return resp, nil
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
