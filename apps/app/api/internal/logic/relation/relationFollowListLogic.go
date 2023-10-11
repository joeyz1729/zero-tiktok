package relation

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

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
	// todo: add your logic here and delete this line
	resp = new(types.FollowListResponse)
	followRes, err := l.svcCtx.FollowRpc.GetFollowIds(l.ctx, &follow.GetFollowIdsRequest{
		UserId: req.UserId,
	})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = err.Error()
		return resp, nil
	}
	logx.Infof("follow ids: %v\n", followRes.FollowIds)

	usersRes := new(user.GetUsersResponse)
	usersRes, err = l.svcCtx.UserRpc.GetUsers(l.ctx, &user.GetUsersRequest{
		UserIds: followRes.FollowIds,
	})
	if err != nil {
		return nil, err
	}
	userList := make([]types.UserInfo, len(followRes.FollowIds))
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
		}
	}
	resp.UserList = userList
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
