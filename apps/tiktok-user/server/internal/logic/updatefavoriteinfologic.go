package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFavoriteInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteInfoLogic {
	return &UpdateFavoriteInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 点赞后更新用户计数
func (l *UpdateFavoriteInfoLogic) UpdateFavoriteInfo(in *pb.UpdateFavoriteInfoRequest) (*pb.UpdateFavoriteInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateFavoriteInfoResponse{}, nil
}
