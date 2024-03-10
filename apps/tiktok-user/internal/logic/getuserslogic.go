package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/follow"
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
	resp := new(pb.GetUsersResponse)
	users := make([]*pb.User, len(in.UserIds))
	for i, uid := range in.UserIds {
		// 查询用户信息
		user, err := l.svcCtx.UserRepo.GetUserDetail(uid)
		if err != nil {
			logx.Errorw("get user detail", logx.Field("err", err))
			return nil, err
		}
		users[i] = &pb.User{
			Id:              user.ID,
			Name:            user.Username,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
		}

		// userId字段为0时不查询关注关系。
		if in.UserId != 0 {
			relation, err := l.svcCtx.FollowRpc.GetRelation(l.ctx, &follow.GetRelationRequest{
				UserId:   in.UserId,
				ToUserId: uid,
			})
			if err != nil {
				logx.Errorw("follow rpc get relation", logx.Field("err", err))
				return nil, err
			}
			users[i].IsFollow = relation.IsFollowing
		}
	}
	resp.UserList = users
	return resp, nil
}
