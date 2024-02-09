package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type FollowPrepareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowPrepareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowPrepareLogic {
	return &FollowPrepareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowPrepareLogic) FollowPrepare(in *pb.FollowPrepareRequest) (*pb.FollowPrepareRequest, error) {
	// todo: add your logic here and delete this line

	return &pb.FollowPrepareRequest{}, nil
}
