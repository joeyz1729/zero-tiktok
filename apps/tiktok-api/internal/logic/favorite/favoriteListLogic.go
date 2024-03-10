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
	idsRes := new(favorite.GetThumbupVideoIdsResponse)
	idsRes, err = l.svcCtx.FavoriteRpc.GetThumbupVideoIds(l.ctx, &favorite.GetThumbupVideoIdsRequest{
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
	videosRes := new(videoservice.GetVideosResponse)
	videosRes, err = l.svcCtx.VideoRpc.GetVideos(l.ctx, &videoservice.GetVideosRequest{
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
