package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/model"

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

func (l *UserInfoLogic) UserInfo(in *model.UserInfoRequest) (*model.UserInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &model.UserInfoResponse{}, nil
}
