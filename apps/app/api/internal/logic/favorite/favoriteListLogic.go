package favorite

import (
	"context"
	"encoding/json"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

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
	// 2. favorite rpc 根据 user_id 查询点赞 video ids 列表
	idsRes := new(favorite.GetVideoIdsResponse)
	idsRes, err = l.svcCtx.FavoriteRpc.GetVideoIds(l.ctx, &favorite.GetVideoIdsRequest{
		UserId: req.UserId,
	})
	length := len(idsRes.VideoIds)
	logx.Info("get favorite video ids: ", idsRes.VideoIds)
	// 出错或者没有数据
	if err != nil {
		logx.Errorw("favorite rpc failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "favorite rpc get video ids: " + err.Error()
		return resp, nil
	}
	if length == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "empty list"
		return resp, nil
	}

	// 3. video rpc 根据 video ids 列表获取详细的视频信息
	videosRes := new(video.GetFavorListResponse)
	videosRes, err = l.svcCtx.VideoRpc.GetFavorList(l.ctx, &video.GetFavorListRequest{
		UserId:   req.UserId,
		VideoIds: idsRes.VideoIds,
	})

	if err != nil || length != len(videosRes.VideoList) {
		logx.Error("get videos by ids: ", err)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "video rpc get videos by ids: " + err.Error()
		return resp, nil
	}
	resp.VideoList = make([]types.Video, length)
	for i, vi := range videosRes.VideoList {
		b, err := json.Marshal(vi)
		if err != nil {
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "marshal video info: " + err.Error()
			return resp, nil
		}
		videoInfo := types.Video{}
		if err = json.Unmarshal(b, &videoInfo); err != nil {
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "unmarshal video info: " + err.Error()
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
//var videoRes = new(video.GetVideoByIdResponse)
//videoRes, err = l.svcCtx.VideoRpc.GetVideoById(l.ctx, &video.GetVideoByIdRequest{
//VideoId: vid,
//})
//if err != nil {
//logx.Errorw("video rpc failed",
//logx.Field("err", err),
//)
//resp.VideoList = []types.Video{}
//resp.StatusCode = http.StatusInternalServerError
//resp.StatusMsg = "video rpc failed"
//return resp, nil
//}
//if videoRes.VideoInfo == nil {
//continue
//}
//var userRes = new(user.GetUserByIdResponse)
//userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
//UserId: videoRes.VideoInfo.AuthorId,
//})
//if err != nil {
//logx.Errorw("user rpc failed",
//logx.Field("err", err),
//)
//resp.StatusCode = http.StatusInternalServerError
//resp.StatusMsg = "user rpc failed"
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
