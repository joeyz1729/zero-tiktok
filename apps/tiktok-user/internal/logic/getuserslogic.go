package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersLogic {
	return &GetUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsersLogic) GetUsers(in *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	//// todo: add your logic here and delete this line
	resp := new(pb.GetUsersResponse)
	//users := make([]*pb.UserInfo, len(in.UserIds))
	//for i, uid := range in.UserIds {
	//	user, err := l.svcCtx.UserRepo.GetUserInfo(uid)
	//	if err != nil {
	//		logx.Error("get tiktok-user info failed")
	//		return nil, err
	//	}
	//	logx.Info(user)
	//	users[i] = &pb.UserInfo{
	//		Id:              uid,
	//		Name:            user.Username,
	//		Avatar:          "no avatar",
	//		BackgroundImage: "no background image",
	//		Signature:       "no signature",
	//
	//		FollowCount:   user.FollowedCount,
	//		FollowerCount: user.FollowerCount,
	//
	//		TotalFavorited: user.TotalFavorited,
	//		FavoriteCount:  user.FavoriteCount,
	//		WorkCount:      user.WorkCount,
	//	}
	//}
	//resp.UserList = users
	return resp, nil
}
