package favorite

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

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
	// todo: add your logic here and delete this line
	resp = new(types.FavoriteActionResponse)
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("jwt parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "invalid token"
		return resp, nil
	}

	_, err = l.svcCtx.VideoRpc.GetVideoById(l.ctx, &video.GetVideoByIdRequest{
		VideoId: req.VideoId,
	})
	if err != nil {
		logx.Errorw("invalid video id",
			logx.Field("err", err))
		return resp, err
	}

	// TODO
	uid := claims.UserId
	// 1 add favorite
	var actionRes = new(favorite.ActionResponse)
	if req.ActionType == 1 {
		actionRes, err = l.svcCtx.FavoriteRpc.AddAction(l.ctx, &favorite.ActionRequest{
			UserId:  uid,
			VideoId: req.VideoId,
		})
	} else if req.ActionType == 2 {
		actionRes, err = l.svcCtx.FavoriteRpc.DelAction(l.ctx, &favorite.ActionRequest{
			UserId:  uid,
			VideoId: req.VideoId,
		})
	} else {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "invalid action type"
		return
	}

	// 2. update user info
	_, err = l.svcCtx.UserRpc.UpdateFavoriteInfo(l.ctx, &user.UpdateFavoriteInfoRequest{
		ActionType: true,
	})
	if err != nil {
		return
	}

	//// 3. update video info
	_, err = l.svcCtx.VideoRpc.UpdateFavoriteCount(l.ctx, &video.UpdateFavoriteCountRequest{
		ActionType: true,
	})
	if err != nil {
		return
	}

	//if err != nil {
	//	resp.StatusCode = http.StatusOK
	//	resp.StatusMsg = "failed"
	//	logx.Errorw("favorite rpc",
	//		logx.Field("err", err),
	//	)
	//	return
	//}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = actionRes.Msg
	return
}
