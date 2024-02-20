package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerIdsLogic {
	return &GetFollowerIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFollowerIds 获取粉丝列表
func (l *GetFollowerIdsLogic) GetFollowerIds(in *pb.GetFollowerIdsRequest) (*pb.GetFollowerIdsResponse, error) {
	resp := new(pb.GetFollowerIdsResponse)
	ids, err := l.svcCtx.FollowRepo.GetFollowerIds(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	resp.FollowerIds = ids
	return resp, nil
}
