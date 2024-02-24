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
	ids, err := l.svcCtx.FollowRepo.GetFollowerIds(l.ctx, in.GetToUserId())
	if err != nil {
		return nil, err
	}
	users, err := l.svcCtx.UserRpc.GetUsers(l.ctx, &userservice.GetUsersRequest{UserIds: ids})
	if err != nil {
		return nil, err
	}
	resp := new(pb.GetFollowerListResponse)
	resp.List = make([]*pb.User, len(users.UserList))
	for i, user := range users.UserList {
		relation, err := l.svcCtx.FollowRepo.CheckRelation(l.ctx, in.UserId, in.ToUserId)
		if err != nil {
			return nil, err
		}
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
			IsFollow:        relation,
		}
	}

	return resp, nil
}
