package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowIdsLogic {
	return &GetFollowIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFollowIds 获取关注列表
func (l *GetFollowIdsLogic) GetFollowIds(in *pb.GetFollowIdsRequest) (*pb.GetFollowIdsResponse, error) {
	resp := new(pb.GetFollowIdsResponse)
	ids, err := l.svcCtx.FollowRepo.GetFollowedIds(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("empty set")
	}
	resp.FollowIds = ids
	return resp, nil
}
