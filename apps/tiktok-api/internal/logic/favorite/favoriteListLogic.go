package favorite

import (
	"context"
	"encoding/json"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/videoservice"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListRequest) (resp *types.FavoriteListResponse, err error) {
	// 1. 用户信息
	resp = new(types.FavoriteListResponse)
	resp.VideoList = []types.Video{}
	// 2. favorite tiktok-user 根据 user_id 查询点赞 videoservice ids 列表
	idsRes := new(favorite.GetVideoIdsResponse)
	idsRes, err = l.svcCtx.FavoriteRpc.GetVideoIds(l.ctx, &favorite.GetVideoIdsRequest{
		UserId: req.UserId,
	})
	length := len(idsRes.VideoIds)
	logx.Info("get favorite videoservice ids: ", idsRes.VideoIds)
	// 出错或者没有数据
	if err != nil {
		logx.Errorw("favorite tiktok-user failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "favorite tiktok-user get videoservice ids: " + err.Error()
		return resp, nil
	}
	if length == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "empty list"
		return resp, nil
	}

	// 3. videoservice tiktok-user 根据 videoservice ids 列表获取详细的视频信息
	videosRes := new(videoservice.GetVideosByIdsResponse)
	videosRes, err = l.svcCtx.VideoRpc.GetVideosByIds(l.ctx, &videoservice.GetVideosByIdsRequest{
		UserId:   req.UserId,
		VideoIds: idsRes.VideoIds,
	})

	if err != nil || length != len(videosRes.VideoList) {
		logx.Error("get videos by ids: ", err)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "videoservice tiktok-user get videos by ids: " + err.Error()
		return resp, nil
	}
	resp.VideoList = make([]types.Video, length)
	for i, vi := range videosRes.VideoList {
		b, err := json.Marshal(vi)
		if err != nil {
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "marshal videoservice info: " + err.Error()
			return resp, nil
		}
		videoInfo := types.Video{}
		if err = json.Unmarshal(b, &videoInfo); err != nil {
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "unmarshal videoservice info: " + err.Error()
			return resp, nil
		}
		resp.VideoList[i] = videoInfo
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return
}

//resp.VideoList = make([]types.Video, idsRes.VideoNum)
//for i, vid := range idsRes.VideoIds {
//
//var videoRes = new(videoservice.GetVideoByIdResponse)
//videoRes, err = l.svcCtx.VideoRpc.GetVideoById(l.ctx, &videoservice.GetVideoByIdRequest{
//VideoId: vid,
//})
//if err != nil {
//logx.Errorw("videoservice tiktok-user failed",
//logx.Field("err", err),
//)
//resp.VideoList = []types.Video{}
//resp.StatusCode = http.StatusInternalServerError
//resp.StatusMsg = "videoservice tiktok-user failed"
//return resp, nil
//}
//if videoRes.VideoInfo == nil {
//continue
//}
//var userRes = new(tiktok-user.GetUserByIdResponse)
//userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &tiktok-user.GetUserByIdRequest{
//UserId: videoRes.VideoInfo.AuthorId,
//})
//if err != nil {
//logx.Errorw("tiktok-user tiktok-user failed",
//logx.Field("err", err),
//)
//resp.StatusCode = http.StatusInternalServerError
//resp.StatusMsg = "tiktok-user tiktok-user failed"
//return resp, nil
//}
//
//resp.VideoList[i] = types.Video{
//Id: vid,
//Author: types.Author(types.User{
//Id:              videoRes.VideoInfo.AuthorId,
//Name:            userRes.Name,
//Avatar:          userRes.Avatar,
//BackgroundImage: userRes.BackgroundImage,
//Signature:       userRes.Signature,
//}),
//PlayUrl:  videoRes.VideoInfo.PlayUrl,
//CoverUrl: videoRes.VideoInfo.CoverUrl,
//Title:    videoRes.VideoInfo.Title,
//}
//}
