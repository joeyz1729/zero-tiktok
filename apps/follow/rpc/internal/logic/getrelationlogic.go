package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/model"

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

// GetRelation 检查指定用户是否关注自己，以及自己是否关注
func (l *GetRelationLogic) GetRelation(in *model.GetRelationRequest) (*model.GetRelationResponse, error) {
	var resp = new(model.GetRelationResponse)
	var err error
	resp.IsFollowing, err = l.svcCtx.FollowRepo.CheckRelation(in.UserId, in.ToUserId)
	if err != nil {
		return nil, err
	}
	resp.IsFollowed, err = l.svcCtx.FollowRepo.CheckRelation(in.ToUserId, in.UserId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
