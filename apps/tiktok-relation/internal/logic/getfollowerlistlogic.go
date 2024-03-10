package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerListLogic {
	return &GetFollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowerListLogic) GetFollowerList(in *pb.GetFollowerListRequest) (*pb.GetFollowerListResponse, error) {
	// 获取id列表
	ids, err := l.svcCtx.FollowRepo.GetFollowerIds(l.ctx, in.ToUserId)
	if err != nil {
		logx.Errorw("get follower ids", logx.Field("err", err),
			logx.Field("userId", in.ToUserId))
		return nil, err
	}

	// 获取详细用户信息
	users, err := l.svcCtx.UserRpc.GetUsers(l.ctx, &userservice.GetUsersRequest{
		UserIds: ids,
		UserId:  in.UserId,
	})
	if err != nil {
		logx.Errorw("user rpc get users", logx.Field("err", err),
			logx.Field("userIds", ids), logx.Field("userId", in.UserId))
		return nil, err
	}

	// 组装结果
	resp := new(pb.GetFollowerListResponse)
	resp.List = make([]*pb.User, len(users.UserList))
	for i, user := range users.UserList {
		resp.List[i] = &pb.User{
			Id:              user.Id,
			Name:            user.Name,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			WorkCount:       user.WorkCount,
			TotalFavorited:  user.TotalFavorited,
			FavoriteCount:   user.FavoriteCount,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			IsFollow:        user.IsFollow,
		}
	}
	return resp, nil
}
