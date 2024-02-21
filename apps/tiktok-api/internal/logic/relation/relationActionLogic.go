package relation

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/follow"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"net/http"

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
	resp = new(types.RelationActionResponse)
	if req.ActionType != int32(1) && req.ActionType != int32(2) {
		resp.StatusCode = http.StatusBadRequest
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

	userId := claims.UserId
	if userId == req.ToUserId {
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMsg = "can not follow oneself"
		return resp, nil
	}
	_, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &userservice.GetUserByIdRequest{UserId: req.ToUserId})
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMsg = err.Error()
		return resp, nil
	}

	if req.ActionType == int32(1) {
		_, err = l.svcCtx.FollowRpc.AddFollow(l.ctx, &follow.AddFollowRequest{
			UserId:   userId,
			ToUserId: req.ToUserId,
		})
	} else {
		_, err = l.svcCtx.FollowRpc.DelFollow(l.ctx, &follow.DelFollowRequest{
			UserId:   userId,
			ToUserId: req.ToUserId,
		})
	}
	if err != nil {
		logx.Errorw("tiktok-user follow action failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = err.Error()
		return resp, nil
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "OK"
	return resp, nil

}
