package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
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

func (l *UserInfoLogic) UserInfo(in *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	resp := new(pb.UserInfoResponse)
	resp.User = new(pb.User)
	user, err := l.svcCtx.UserRepo.GetUserById(in.UserId)
	if err != nil {
		logx.Errorw("get tiktok-user info failed",
			logx.Field("err", err))
		return resp, err
	}
	resp.User.Id = in.UserId
	resp.User.Name = user.Username
	resp.User.Avatar = user.Avatar
	resp.User.BackgroundImage = user.BackgroundImage
	resp.User.Signature = user.Signature
	return resp, nil
}
