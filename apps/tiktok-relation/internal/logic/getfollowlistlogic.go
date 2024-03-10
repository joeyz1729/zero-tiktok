package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/pb"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFollowList 获取点赞列表，并检查对方是否关注自己
func (l *GetFollowListLogic) GetFollowList(in *pb.GetFollowListRequest) (*pb.GetFollowListResponse, error) {
	ids, err := l.svcCtx.FollowRepo.GetFollowedIds(l.ctx, in.ToUserId)
	if err != nil {
		logx.Errorw("get followed ids", logx.Field("err", err),
			logx.Field("userId", in.ToUserId))
		return nil, err
	}
	logx.Infow("followed ids", logx.Field("ids", ids))
	users, err := l.svcCtx.UserRpc.GetUsers(l.ctx, &userservice.GetUsersRequest{
		UserIds: ids,
		UserId:  in.UserId,
	})
	logx.Infow("get users", logx.Field("id", users.UserList[0].Name))
	if err != nil {
		logx.Errorw("user rpc get users", logx.Field("err", err),
			logx.Field("userIds", ids), logx.Field("userId", in.UserId))
		return nil, err
	}

	resp := new(pb.GetFollowListResponse)
	resp.List = make([]*pb.User, len(users.UserList))
	for i, user := range users.UserList {
		resp.List[i] = &pb.User{
			Id:              user.Id,
			Name:            user.Name,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
			IsFollow:        user.IsFollow,
		}
	}

	return resp, nil
}
