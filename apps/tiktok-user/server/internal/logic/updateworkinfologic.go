package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/pb"

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

// 发布视频后更新用户计数
func (l *UpdateWorkInfoLogic) UpdateWorkInfo(in *pb.UpdateWorkInfoRequest) (*pb.UpdateWorkInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateWorkInfoResponse{}, nil
}
