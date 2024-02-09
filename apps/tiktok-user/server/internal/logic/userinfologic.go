package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户信息
func (l *UserInfoLogic) UserInfo(in *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UserInfoResponse{}, nil
}
