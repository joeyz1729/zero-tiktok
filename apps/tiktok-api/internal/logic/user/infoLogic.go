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
	var userRes = &userservice.GetUsersResponse{}
	userRes, err = l.svcCtx.UserRpc.GetUsers(l.ctx, &userservice.GetUsersRequest{
		UserIds: []int64{toUserId}})
	if err != nil {
		return nil, err
	}
	resp.UserInfo = types.User{
		Id:              userRes.UserList[0].Id,
		Name:            userRes.UserList[0].Name,
		Avatar:          userRes.UserList[0].Avatar,
		BackgroundImage: userRes.UserList[0].BackgroundImage,
		Signature:       userRes.UserList[0].Signature,
		FollowCount:     userRes.UserList[0].FollowCount,
		FollowerCount:   userRes.UserList[0].FollowerCount,
		TotalFavorited:  userRes.UserList[0].TotalFavorited,
		FavoriteCount:   userRes.UserList[0].FavoriteCount,
	}

	return resp, nil
}
