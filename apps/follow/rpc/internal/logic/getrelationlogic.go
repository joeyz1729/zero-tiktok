package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
)
type GetRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationLogic {
	return &GetRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationLogic) GetRelation(in *model.GetRelationRequest) (*model.GetRelationResponse, error) {
	// todo: add your logic here and delete this line
	var resp = new(model.GetRelationResponse)
	var err error
	logx.Info("begin get follow relation")
	resp.IfFollowing = int32(2)
	// cache
	ok, err := l.svcCtx.FollowCache.GetRelation(l.ctx, in.UserId, in.ToUserId)
	if err == nil && ok {
		// cache hit
		resp.IfFollowing = int32(1)
		return resp, nil
	}
	// miss or failed
	ok, err = l.svcCtx.FollowDB.CheckRelation(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		return resp, err
	}
	if !ok {
		return resp, nil
	}

	resp.IfFollowing = int32(1)
	err = l.svcCtx.FollowCache.AddFollow(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
