package video

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/videoservice"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"github.com/joeyz1729/zero-tiktok/pkg/mw/minio"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"path"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewPublishActionLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		r:      r,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishActionLogic) PublishAction(token string, filename string, title string, fileData []byte) (resp *types.PublishActionResponse, err error) {
	resp = new(types.PublishActionResponse)
	// 1. jwt
	claim, err := jwtx.ParseToken(token)
	if err != nil {
		logx.Errorw("parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusUnauthorized
		resp.StatusMsg = "invalid auth token"
		return resp, nil
	}
	// 2. upload video and cover
	playURL, coverURL, err := minio.UploadVideo(l.ctx, claim.UserId, path.Ext(filename), fileData)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	// 3. logic
	_, err = l.svcCtx.VideoRpc.PublishAction(l.ctx, &videoservice.PublishActionRequest{
		UserId:   claim.UserId,
		Title:    title,
		CoverUrl: coverURL,
		PlayUrl:  playURL,
	})
	if err != nil {
		logx.Errorw("publish action failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	// 4. return
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
