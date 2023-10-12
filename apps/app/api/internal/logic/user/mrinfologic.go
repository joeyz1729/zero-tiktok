package user

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type MrInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMrInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MrInfoLogic {
	return &MrInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MrInfoLogic) MrInfo(req *types.UserMrInfoRequest) (resp *types.UserMrInfoResponse, err error) {
	// token 校验
	resp = new(types.UserMrInfoResponse)
	//claims, err := jwtx.ParseToken(req.Token)
	//if err != nil {
	//	logx.Errorw("parse token failed, or token is invalid",
	//		logx.Field("err", err))
	//	resp.StatusCode = http.StatusUnauthorized
	//	resp.StatusMsg = "unauthorized"
	//	return resp, nil
	//}
	//
	//userId := claims.UserId // 该用户id
	//toUserId := req.UserId  // 查询用户id
	//resp.UserInfo = types.User{}
	//errChan := make(chan error, 7)
	//defer close(errChan)
	//var wg sync.WaitGroup
	//wg.Add(4)
	//
	//// user 模块
	//go func() {
	//	defer wg.Done()
	//	var userRes *user.GetUserByIdResponse
	//	userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{UserId: toUserId})
	//	if err != nil {
	//		errChan <- err
	//	} else {
	//		resp.UserInfo.Id = toUserId
	//		resp.UserInfo.Name = userRes.Name
	//		resp.UserInfo.Avatar = userRes.Avatar
	//		resp.UserInfo.BackgroundImage = userRes.BackgroundImage
	//		resp.UserInfo.Signature = userRes.Signature
	//	}
	//}()
	//
	//// 调用follow模块
	//go func() {
	//	defer wg.Done()
	//	var followRes *follow.GetFollowCountResponse
	//	followRes, err = l.svcCtx.FollowRpc.GetFollowCount(l.ctx, &follow.GetFollowCountRequest{UserId: userId, ToUserId: toUserId})
	//	if err != nil {
	//		errChan <- err
	//	} else {
	//		resp.UserInfo.FollowCount = followRes.FollowCount
	//		resp.UserInfo.FollowerCount = followRes.FollowerCount
	//		resp.UserInfo.IsFollow = followRes.IsFollowing
	//	}
	//}()
	//
	//// 调用video模块
	//go func() {
	//	defer wg.Done()
	//	var videoRes *video.GetWorkCountResponse
	//	videoRes, err = l.svcCtx.VideoRpc.GetWorkCount(l.ctx, &video.GetWorkCountRequest{UserId: toUserId})
	//	if err != nil {
	//		errChan <- err
	//	} else {
	//		resp.UserInfo.WorkCount = videoRes.WorkCount
	//	}
	//}()
	//
	////like模块
	//go func() {
	//	defer wg.Done()
	//	var likeRes *favorite.GetFavoriteCountResponse
	//	likeRes, err = l.svcCtx.FavoriteRpc.GetFavoriteCount(l.ctx, &favorite.GetFavoriteCountRequest{UserId: toUserId})
	//	if err != nil {
	//		errChan <- err
	//	} else {
	//		resp.UserInfo.TotalFavorited = likeRes.TotalFavorited
	//		resp.UserInfo.FavoriteCount = likeRes.FavoriteCount
	//	}
	//}()
	//wg.Wait()
	//select {
	//case result := <-errChan:
	//	resp.StatusCode = http.StatusInternalServerError
	//	resp.StatusMsg = "internal server error"
	//	logx.Error(result.Error())
	//	return resp, nil
	//default:
	//}
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
