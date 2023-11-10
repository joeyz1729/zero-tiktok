package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowListLogic) GetFollowList(in *model.GetFollowListRequest) (*model.GetFollowListResponse, error) {
	resp := new(model.GetFollowListResponse)
	ids, err := l.svcCtx.FollowRepo.GetFollowedIds(in.UserId)
	if err != nil {
		return nil, err
	}
	resp.FollowedIds = ids
	resp.Relations = make([]bool, len(ids))
	for i, id := range ids {
		ok, err := l.svcCtx.FollowRepo.CheckRelation(in.UserId, id)
		if err != nil {
			return nil, err
		}
		resp.Relations[i] = ok
	}
	return resp, nil
}
