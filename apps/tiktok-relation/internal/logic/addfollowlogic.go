package logic

import (
	"context"
	"errors"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrRepeatedOperation = errors.New("repeated operation")
)

type AddFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFollowLogic {
	return &AddFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddFollowLogic) AddFollow(in *pb.AddFollowRequest) (*pb.AddFollowResponse, error) {
	exist, err := l.svcCtx.FollowRepo.CheckRelation(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrRepeatedOperation
	}
	err = l.svcCtx.FollowRepo.AddRelation(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		return nil, err
	}
	return &pb.AddFollowResponse{}, nil
}
