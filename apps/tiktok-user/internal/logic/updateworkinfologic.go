package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWorkInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkInfoLogic {
	return &UpdateWorkInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateWorkInfoLogic) UpdateWorkInfo(in *pb.UpdateWorkInfoRequest) (*pb.UpdateWorkInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateWorkInfoResponse{}, nil
}
