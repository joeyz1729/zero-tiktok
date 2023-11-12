package relation

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"net/http"

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
	resp = new(types.FollowerListResponse)

	followerRes, err := l.svcCtx.FollowRpc.GetFollowerIds(l.ctx, &follow.GetFollowerIdsRequest{
		UserId: req.UserId,
	})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = err.Error()
		return resp, nil
	}

	logx.Infof("follow ids: %v\n", followerRes.FollowerIds)

	usersRes := new(user.GetUsersResponse)
	usersRes, err = l.svcCtx.UserRpc.GetUsers(l.ctx, &user.GetUsersRequest{
		UserIds: followerRes.FollowerIds,
	})
	if err != nil {
		return nil, err
	}
	userList := make([]types.UserInfo, len(followerRes.FollowerIds))
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

//func (l *RelationFollowerListLogic) RelationFollowerListByPage(req *types.FollowerListRequest) (resp *types.FollowerListResponse, err error) {
//	resp = new(types.FollowerListResponse)
//
//	var page = pagination.PageToken(req.PageToken).Decode()
//	var cursor uint64 = 0
//	var size int64 = 2
//	if time.Now().Unix()-page.NextTimeAtUTC < int64(time.Hour)*24 {
//		// not expire
//		if page.NextId > 0 {
//			cursor = page.NextId
//		}
//		if page.PageSize > 0 && page.PageSize <= 10 {
//			size = page.PageSize
//		}
//	}
//
//	followerRes, err := l.svcCtx.FollowRpc.GetFollowerIds(l.ctx, &follow.GetFollowerIdsRequest{
//		UserId:   req.UserId,
//		PageSize: size,
//		Cursor:   cursor,
//	})
//
//	if err != nil {
//		resp.StatusCode = http.StatusInternalServerError
//		resp.StatusMsg = "follow rpc failed"
//		return resp, nil
//	}
//	if len(followerRes.FollowerIds) == 0 {
//		resp.StatusCode = http.StatusOK
//		resp.StatusMsg = "empty list"
//		return resp, nil
//	}
//	logx.Infof("follow ids: %v\n", followerRes.FollowerIds)
//
//	logx.Info("start get user detail")
//	userList := make([]types.User, len(followerRes.FollowerIds))
//	for i, id := range followerRes.FollowerIds {
//		userRes := new(user.GetUserByIdResponse)
//		userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
//			UserId: id,
//		})
//		if err != nil {
//			logx.Errorw("user rpc failed",
//				logx.Field("err", err),
//			)
//			resp.StatusCode = http.StatusInternalServerError
//			resp.StatusMsg = "user rpc failed"
//			return resp, nil
//		}
//		userInfo := types.User{
//			Id:              userRes.Id,
//			Name:            userRes.Name,
//			Avatar:          userRes.Avatar,
//			BackgroundImage: userRes.BackgroundImage,
//			Signature:       userRes.Signature,
//		}
//		userList[i] = userInfo
//	}
//	//fmt.Println(userList)
//
//	logx.Info("generate next page token")
//	nextPage := pagination.Page{
//		NextId:        uint64(followerRes.NextCursor),
//		NextTimeAtUTC: time.Now().Unix(),
//		PageSize:      size,
//	}
//	resp.NextToken = string(nextPage.Encode())
//
//	//resp.UserList = userList
//	resp.StatusCode = http.StatusOK
//	resp.StatusMsg = "success"
//	return resp, nil
//}

//for i, id := range followerRes.FollowerIds {
//	userRes := new(user.GetUserByIdResponse)
//	userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
//		UserId: id,
//	})
//	if err != nil {
//		logx.Errorw("user rpc failed",
//			logx.Field("err", err),
//		)
//		resp.StatusCode = http.StatusInternalServerError
//		resp.StatusMsg = "user rpc failed"
//		return resp, nil
//	}
//	userList[i] = types.UserInfo{
//		Id:              userRes.Id,
//		Name:            userRes.Name,
//		Avatar:          userRes.Avatar,
//		BackgroundImage: userRes.BackgroundImage,
//		Signature:       userRes.Signature,
//	}
//}

//fmt.Println(userList)

//logx.Info("generate next page token")
//nextPage := pagination.Page{
//	NextId: uint64(followerRes.NextCursor),
//}
//resp.NextToken = string(nextPage.Encode())
