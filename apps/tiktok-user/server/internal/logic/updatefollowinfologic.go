package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/pb"

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

// rpc
func (l *UpdateFollowInfoLogic) UpdateFollowInfo(in *pb.UpdateFollowInfoRequest) (*pb.UpdateFollowInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateFollowInfoResponse{}, nil
}
