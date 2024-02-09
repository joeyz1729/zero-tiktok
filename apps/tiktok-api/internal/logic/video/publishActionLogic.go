package video

import (
	"context"
	"errors"
	"fmt"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/video"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"io"
	"net/http"
	"path"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	maxSize = 10 << 20
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

func (l *PublishActionLogic) PublishAction() (resp *types.PublishActionResponse, err error) {
	resp = new(types.PublishActionResponse)
	err = l.r.ParseMultipartForm(maxSize)
	if err != nil {
		logx.Errorw("parse multipart form failed",
			logx.Field("err", err))
		return resp, err
	}

	token := l.r.FormValue("token")
	title := l.r.FormValue("title")
	logx.Infof("token: %s, title: %s", token, title)
	if len(token) == 0 || len(title) == 0 {
		return resp, errors.New("invalid params")
	}

	claim, err := jwtx.ParseToken(token)
	if err != nil {
		logx.Errorw("parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusUnauthorized
		resp.StatusMsg = "invalid auth token"
		return resp, nil
	}

	file, header, err := l.r.FormFile("repository")
	typ := path.Ext(header.Filename)
	fmt.Println(typ)
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		logx.Errorw("read file repository failed",
			logx.Field("err", err))
		return resp, nil
	}
	fmt.Printf("read video file, size: %d", len(fileData))

	_, err = l.svcCtx.VideoRpc.PublishAction(l.ctx, &video.PublishActionRequest{
		UserId: claim.UserId,
		Title:  title,
		Data:   fileData,
		Type:   typ,
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
