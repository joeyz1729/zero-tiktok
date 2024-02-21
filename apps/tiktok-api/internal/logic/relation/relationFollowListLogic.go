package relation

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/follow"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationFollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationFollowListLogic {
	return &RelationFollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationFollowListLogic) RelationFollowList(req *types.FollowListRequest) (resp *types.FollowListResponse, err error) {
	resp = new(types.FollowListResponse)
	// 获取列表
	followRes, err := l.svcCtx.FollowRpc.GetFollowList(l.ctx, &follow.GetFollowListRequest{
		UserId: req.UserId,
	})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = err.Error()
		return resp, nil
	}
	// 获取用户详细信息
	usersRes := new(userservice.GetUsersResponse)
	usersRes, err = l.svcCtx.UserRpc.GetUsers(l.ctx, &userservice.GetUsersRequest{
		UserIds: followRes.FollowedIds,
	})
	if err != nil {
		return nil, err
	}
	// 拼接结果
	userList := make([]types.UserInfo, len(followRes.FollowedIds))
	for i, userInfo := range usersRes.UserList {
		userList[i] = types.UserInfo{
			Id:              userInfo.Id,
			Name:            userInfo.Name,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			Signature:       userInfo.Signature,

			FollowCount:   userInfo.FollowCount,
			FollowerCount: userInfo.FollowerCount,

			TotalFavorited: userInfo.TotalFavorited,
			FavoriteCount:  userInfo.FavoriteCount,
			WorkCount:      userInfo.WorkCount,

			IsFollow: followRes.Relations[i],
		}
	}
	resp.UserList = userList
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
