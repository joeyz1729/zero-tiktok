package video

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListRequest) (resp *types.PublishListResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.PublishListResponse)

	_, err = jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	//myId := claims.UserId
	//loginUserId := claims.UserId

	// 1. 获取user信息
	authorRes := new(user.GetAuthorResponse)
	authorRes, err = l.svcCtx.UserRpc.GetAuthor(l.ctx, &user.GetUserByIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logx.Error("get author ", err)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = err.Error()
		return resp, nil
	}

	// 查询videos基础信息，以及点赞信息
	videosRes := new(video.GetListByAuthorIdResponse)
	videosRes, err = l.svcCtx.VideoRpc.GetListByAuthorId(l.ctx, &video.GetListByAuthorIdRequest{
		UserId: req.UserId,
	})

	// 拼接结果

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}

//
//// 并行请求
//var wg sync.WaitGroup
//wg.Add(5)
//errChan := make(chan error, 5)
//defer close(errChan)
//
//// 1. 根据userid 查询 user 信息
//var author types.User
//go func() {
//	defer wg.Done()
//	userInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{UserId: req.UserId})
//	if err != nil {
//		logx.Errorw("get user by id failed",
//			logx.Field("err", err),
//		)
//		errChan <- err
//	} else {
//		author.Id = userInfo.Id
//		author.Name = userInfo.Name
//		author.Avatar = userInfo.Avatar
//		author.BackgroundImage = userInfo.BackgroundImage
//	}
//}()
//
//// 2. follow人数
//go func() {
//	defer wg.Done()
//	var followRes *follow.GetFollowCountResponse
//	followRes, err = l.svcCtx.FollowRpc.GetFollowCount(l.ctx, &follow.GetFollowCountRequest{UserId: myId, ToUserId: req.UserId})
//	if err != nil {
//		errChan <- err
//	} else {
//		author.FollowCount = followRes.FollowCount
//		author.FollowerCount = followRes.FollowerCount
//		author.IsFollow = followRes.IsFollowing
//	}
//}()
//
//// 3. 发布视频数
//go func() {
//	defer wg.Done()
//	var videoRes *video.GetWorkCountResponse
//	videoRes, err = l.svcCtx.VideoRpc.GetWorkCount(l.ctx, &video.GetWorkCountRequest{UserId: req.UserId})
//	if err != nil {
//		errChan <- err
//	} else {
//		author.WorkCount = videoRes.WorkCount
//	}
//}()
//
//// 4.点赞总数
//go func() {
//	defer wg.Done()
//	var likeRes *favorite.GetFavoriteCountResponse
//	likeRes, err = l.svcCtx.FavoriteRpc.GetFavoriteCount(l.ctx, &favorite.GetFavoriteCountRequest{UserId: req.UserId})
//	if err != nil {
//		errChan <- err
//	} else {
//		author.TotalFavorited = likeRes.TotalFavorited
//		author.FavoriteCount = likeRes.FavoriteCount
//	}
//}()
//
//// 5. 发布视频列表
//var listRes = new(video.GetListByUserIdResponse)
//go func() {
//	defer wg.Done()
//	listRes, err = l.svcCtx.VideoRpc.GetListByUserId(l.ctx, &video.GetListByUserIdRequest{
//		UserId: req.UserId,
//	})
//	if err != nil {
//		logx.Errorw("get list by user id failed or empty list",
//			logx.Field("err", err),
//		)
//		errChan <- err
//		return
//	} else if len(listRes.GetVideoList()) == 0 {
//		errChan <- errors.New("empty video list")
//		return
//	}
//
//}()
//wg.Wait()
//select {
//case result := <-errChan:
//resp.StatusCode = http.StatusInternalServerError
//resp.StatusMsg = "internal server error"
//logx.Error(result.Error())
//default:
//}
//
//resp.VideoList = make([]types.Video, len(listRes.GetVideoList()))
//for i, v := range listRes.GetVideoList() {
//var videoInfo types.Video
//videoInfo = types.Video{
//Id:       v.GetVideoId(),
//Author:   types.Author(author),
//PlayUrl:  v.GetPlayUrl(),
//CoverUrl: v.GetCoverUrl(),
//Title:    v.GetTitle(),
//}
//resp.VideoList[i] = videoInfo
//}
