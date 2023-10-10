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

	userList := make([]types.UserInfo, len(followRes.FollowIds))
	for i, id := range followRes.FollowIds {
		userRes := new(user.GetUserByIdResponse)
		userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
			UserId: id,
		})
		if err != nil {
			logx.Errorw("user rpc failed",
				logx.Field("err", err),
			)
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = err.Error()
			return resp, nil
		}
		userList[i] = types.UserInfo{
			Id:              userRes.Id,
			Name:            userRes.Name,
			Avatar:          userRes.Avatar,
			BackgroundImage: userRes.BackgroundImage,
			Signature:       userRes.Signature,
		}
		//TODO, comprehensive user info

	}
	//fmt.Println(userList)
	resp.UserList = userList
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}

func (l *RelationFollowListLogic) RelationFollowListByPage(req *types.FollowListRequest) (resp *types.FollowListResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.FollowListResponse)
	followRes, err := l.svcCtx.FollowRpc.GetFollowIds(l.ctx, &follow.GetFollowIdsRequest{
		UserId: req.UserId,
	})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "follow rpc failed"
		return resp, nil
	}
	if len(followRes.FollowIds) == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "empty list"
		return resp, nil
	}
	logx.Infof("follow ids: %v\n", followRes.FollowIds)
	logx.Info("start get user detail")
	userList := make([]types.User, len(followRes.FollowIds))
	for i, id := range followRes.FollowIds {
		userRes := new(user.GetUserByIdResponse)
		userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
			UserId: id,
		})
		if err != nil {
			logx.Errorw("user rpc failed",
				logx.Field("err", err),
			)
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "user rpc failed"
			return resp, nil
		}
		userInfo := types.User{
			Id:              userRes.Id,
			Name:            userRes.Name,
			Avatar:          userRes.Avatar,
			BackgroundImage: userRes.BackgroundImage,
			Signature:       userRes.Signature,
		}
		//TODO, comprehensive user info
		userList[i] = userInfo
	}
	//fmt.Println(userList)
	//resp.UserList = userList
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
