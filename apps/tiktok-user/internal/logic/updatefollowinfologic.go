package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFollowInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFollowInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFollowInfoLogic {
	return &UpdateFollowInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateFollowInfo 更新
func (l *UpdateFollowInfoLogic) UpdateFollowInfo(in *pb.UpdateFollowInfoRequest) (*pb.UpdateFollowInfoResponse, error) {
	var err error
	if in.ActionType {
		err = l.svcCtx.UserRepo.AddFollow(in.UserId, in.ToUserId)
	} else {
		err = l.svcCtx.UserRepo.DelFollow(in.UserId, in.ToUserId)
	}
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	return &pb.UpdateFollowInfoResponse{}, nil
}
