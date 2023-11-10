package favorite

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteActionLogic) FavoriteAction(req *types.FavoriteActionRequest) (resp *types.FavoriteActionResponse, err error) {
	resp = new(types.FavoriteActionResponse)
	// action type
	if req.ActionType != int32(1) && req.ActionType != int32(2) {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "invalid action type"
		return
	}
	// user id
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("jwt parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "invalid token: " + err.Error()
		return resp, nil
	}
	uid := claims.UserId
	// video id
	videoRes := new(video.GetVideoByIdResponse)
	videoRes, err = l.svcCtx.VideoRpc.GetVideoById(l.ctx, &video.GetVideoByIdRequest{
		VideoId: req.VideoId,
	})
	if err != nil {
		logx.Errorw("invalid video id",
			logx.Field("err", err))
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "check video id failed: " + err.Error()
		return resp, nil
	}
	authorId := videoRes.VideoInfo.AuthorId
	// start
	_, err = l.svcCtx.FavoriteRpc.Action(l.ctx, &favorite.ActionRequest{
		UserId:     uid,
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
		AuthorId:   authorId,
	})
	if err != nil {
		logx.Errorw("favorite action",
			logx.Field("err", err))
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "update action failed:" + err.Error()
		return resp, nil
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "OK"
	return resp, nil
}
