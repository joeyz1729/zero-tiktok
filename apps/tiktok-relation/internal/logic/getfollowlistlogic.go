package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/pb"

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

// GetFollowList 获取点赞列表，并检查对方是否关注自己
func (l *GetFollowListLogic) GetFollowList(in *pb.GetFollowListRequest) (*pb.GetFollowListResponse, error) {
	resp := new(pb.GetFollowListResponse)
	ids, err := l.svcCtx.FollowRepo.GetFollowedIds(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	resp.FollowedIds = ids
	resp.Relations = make([]bool, len(ids))
	for i, id := range ids {
		ok, err := l.svcCtx.FollowRepo.CheckRelation(l.ctx, in.UserId, id)
		if err != nil {
			return nil, err
		}
		resp.Relations[i] = ok
	}
	return resp, nil
}
