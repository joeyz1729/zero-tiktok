package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/follow"
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
	user, err := l.svcCtx.UserRepo.GetUserDetail(in.ToUserId)
	if err != nil {
		logx.Errorw("get tiktok-user info failed",
			logx.Field("err", err))
		return resp, err
	}
	resp.User = &pb.User{
		Id:              user.Id,
		Name:            user.Name,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		FollowerCount:   user.FollowerCount,
		FollowCount:     user.FollowCount,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
	relation, err := l.svcCtx.FollowRpc.GetRelation(l.ctx, &follow.GetRelationRequest{UserId: in.UserId, ToUserId: in.ToUserId})
	if err != nil {
		return nil, err
	}
	resp.User.IsFollow = relation.IsFollowing
	return resp, nil
}
