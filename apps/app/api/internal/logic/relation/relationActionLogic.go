package relation

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationActionLogic) RelationAction(req *types.RelationActionRequest) (resp *types.RelationActionResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.RelationActionResponse)
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("jwt parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusUnauthorized
		resp.StatusMsg = "invalid token"
		return resp, nil
	}

	userId := claims.UserId
	if userId == req.ToUserId {
		resp.StatusCode = http.StatusServiceUnavailable
		resp.StatusMsg = "can not follow oneself"
		return resp, nil
	}
	var msg string
	if req.ActionType == int32(1) {
		res := new(follow.AddResponse)
		res, err = l.svcCtx.FollowRpc.Add(l.ctx, &follow.AddRequest{
			UserId:   userId,
			ToUserId: req.ToUserId,
		})
		msg = res.Msg
	} else {
		res := new(follow.DelResponse)
		res, err = l.svcCtx.FollowRpc.Del(l.ctx, &follow.DelRequest{
			UserId:   userId,
			ToUserId: req.ToUserId,
		})
		msg = res.Msg
	}

	if err != nil {
		logx.Errorw("rpc follow add failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = msg
	return resp, nil

}
