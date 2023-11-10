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
