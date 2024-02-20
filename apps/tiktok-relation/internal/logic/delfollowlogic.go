package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelFollowLogic {
	return &DelFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelFollowLogic) DelFollow(in *pb.DelFollowRequest) (*pb.DelFollowResponse, error) {
	exist, err := l.svcCtx.FollowRepo.CheckRelation(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrRepeatedOperation
	}
	err = l.svcCtx.FollowRepo.DelRelation(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		return nil, err
	}
	return &pb.DelFollowResponse{}, nil
}
