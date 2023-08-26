package favorite

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/favorite"
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
	uid := claims.UserId
	actionRes, err := l.svcCtx.FavoriteRpc.Action(l.ctx, &favorite.ActionRequest{
		UserId:     uid,
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
	})
	if err != nil {
		return
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = actionRes.Msg
	return
}
