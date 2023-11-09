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
	resp := new(model.UserInfoResponse)
	resp.User = new(model.UserInfo)
	user, err := l.svcCtx.UserRepo.GetUserInfo(in.UserId)
	if err != nil {
		logx.Errorw("get user info failed",
			logx.Field("err", err))
		return resp, err
	}
	resp.User.Id = in.UserId
	resp.User.Name = user.Username
	resp.User.Avatar = user.Avatar
	resp.User.BackgroundImage = user.BackgroundImage
	resp.User.Signature = user.Signature

	resp.User.FollowCount = user.FollowedCount
	resp.User.FollowerCount = user.FollowerCount

	resp.User.TotalFavorited = user.TotalFavorited
	resp.User.WorkCount = user.WorkCount
	resp.User.FavoriteCount = user.FavoriteCount
	return resp, nil
}
