package favorite

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/videoservice"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
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
	// tiktok-user id
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
	// videoservice id
	videoRes := new(videoservice.GetVideoByIdResponse)
	videoRes, err = l.svcCtx.VideoRpc.GetVideoById(l.ctx, &videoservice.GetVideoByIdRequest{
		VideoId: req.VideoId,
	})
	if err != nil {
		logx.Errorw("invalid videoservice id",
			logx.Field("err", err))
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "check videoservice id failed: " + err.Error()
		return resp, nil
	}
	authorId := videoRes.VideoInfo.AuthorId
	// 添加更新，需要修改关系，user的点赞数，author的被点赞数，video的被点赞数
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
