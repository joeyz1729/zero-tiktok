package user

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	// token 校验
	resp = new(types.UserInfoResponse)
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed, or token is invalid",
			logx.Field("err", err))
		return resp, nil
	}
	userId := claims.UserId // 该用户id
	logx.Infof("userid:%d\n", userId)
	toUserId := req.UserId // 查询用户id
	resp.UserInfo = types.User{}
	var userRes = &userservice.UserInfoResponse{}
	userRes, err = l.svcCtx.UserRpc.UserInfo(l.ctx, &userservice.UserInfoRequest{UserId: toUserId})
	if err != nil {
		return nil, err
	}
	resp.UserInfo.Id = toUserId
	resp.UserInfo.Name = userRes.User.Name
	resp.UserInfo.Avatar = userRes.User.Avatar
	resp.UserInfo.BackgroundImage = userRes.User.BackgroundImage
	resp.UserInfo.Signature = userRes.User.Signature

	//resp.UserInfo.FollowCount = userRes.User.FollowCount
	//resp.UserInfo.FollowerCount = userRes.User.FollowerCount
	//resp.UserInfo.TotalFavorited = userRes.User.TotalFavorited
	//resp.UserInfo.FavoriteCount = userRes.User.FavoriteCount
	return resp, nil
}
