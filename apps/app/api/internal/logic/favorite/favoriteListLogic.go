package favorite

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
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
	// todo: add your logic here and delete this line
	// get video list by user id
	resp = new(types.FavoriteListResponse)

	// get vid list by uid
	idsRes, err := l.svcCtx.FavoriteRpc.GetVideoIds(l.ctx, &favorite.GetVideoIdsRequest{
		UserId: req.UserId,
	})
	logx.Info(idsRes.VideoIds)
	if err != nil {
		logx.Errorw("favorite rpc failed",
			logx.Field("err", err),
		)
		resp.VideoList = []types.Video{}
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "favorite rpc failed"
		return resp, nil
	}
	if idsRes.VideoNum == 0 {
		resp.VideoList = []types.Video{}
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "empty list"
		return resp, nil
	}
	resp.VideoList = make([]types.Video, idsRes.VideoNum)
	for i, vid := range idsRes.VideoIds {

		var videoRes = new(video.GetVideoByIdResponse)
		videoRes, err = l.svcCtx.VideoRpc.GetVideoById(l.ctx, &video.GetVideoByIdRequest{
			VideoId: vid,
		})
		if err != nil {
			logx.Errorw("video rpc failed",
				logx.Field("err", err),
			)
			resp.VideoList = []types.Video{}
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "video rpc failed"
			return resp, nil
		}
		if videoRes.VideoInfo == nil {
			continue
		}
		var userRes = new(user.GetUserByIdResponse)
		userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
			UserId: videoRes.VideoInfo.AuthorId,
		})
		if err != nil {
			logx.Errorw("user rpc failed",
				logx.Field("err", err),
			)
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "user rpc failed"
			return resp, nil
		}

		resp.VideoList[i] = types.Video{
			Id: vid,
			Author: types.Author(types.User{
				Id:              videoRes.VideoInfo.AuthorId,
				Name:            userRes.Name,
				Avatar:          userRes.Avatar,
				BackgroundImage: userRes.BackgroundImage,
				Signature:       userRes.Signature,
			}),
			PlayUrl:  videoRes.VideoInfo.PlayUrl,
			CoverUrl: videoRes.VideoInfo.CoverUrl,
			Title:    videoRes.VideoInfo.Title,
		}
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return
}
