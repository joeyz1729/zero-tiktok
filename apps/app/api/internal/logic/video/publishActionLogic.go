package video

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishActionLogic) PublishAction(req *types.PublishActionRequest) (resp *types.PublishActionResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.PublishActionResponse)
	// get user_id from jwt token
	claim, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusUnauthorized
		resp.StatusMsg = "invalid auth token"
		return resp, nil
	}

	_, err = l.svcCtx.VideoRpc.PublishAction(l.ctx, &video.PublishActionRequest{
		UserId: claim.UserId,
		Title:  req.Title,
		Data:   req.Data,
	})
	if err != nil {
		logx.Errorw("publish action failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
