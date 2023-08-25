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

type RelationFollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationFollowerListLogic {
	return &RelationFollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationFollowerListLogic) RelationFollowerList(req *types.FollowerListRequest) (resp *types.FollowerListResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.FollowerListResponse)
	followerRes, err := l.svcCtx.FollowRpc.GetFollowerIds(l.ctx, &follow.GetFollowerIdsRequest{
		UserId: req.UserId,
	})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "follow rpc failed"
		return resp, nil
	}
	if len(followerRes.FollowerIds) == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "empty list"
		return resp, nil
	}
	logx.Infof("follow ids: %v\n", followerRes.FollowerIds)
	logx.Info("start get user detail")
	userList := make([]types.User, len(followerRes.FollowerIds))
	for i, id := range followerRes.FollowerIds {
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
	resp.UserList = userList
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
