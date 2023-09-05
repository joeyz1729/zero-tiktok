package relation

import (
	"context"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationFollowCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationFollowCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationFollowCheckLogic {
	return &RelationFollowCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationFollowCheckLogic) RelationFollowCheck(req *types.FollowCheckRequest) (resp *types.FollowCheckResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.FollowCheckResponse)
	res, err := l.svcCtx.FollowRpc.GetRelation(l.ctx, &follow.GetRelationRequest{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
	})
	if err != nil {
		resp.StatusMsg = fmt.Sprintf("rpc failed: %s", err.Error())
		resp.StatusCode = http.StatusInternalServerError
		return
	}
	resp.IfFollowing = res.IfFollowing
	resp.IfFollower = res.IfFollower
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
