package video

import (
	"bytes"
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/videoservice"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"github.com/joeyz1729/zero-tiktok/pkg/mw/ffmpeg"
	"github.com/joeyz1729/zero-tiktok/pkg/mw/minio"
	"github.com/joeyz1729/zero-tiktok/pkg/utils"
	"io"
	"net/http"
	"path"
	"time"

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
	defer func() {
		if err = file.Close(); err != nil {
			logx.Errorw("close file", logx.Field("err", err))
		}
	}()
	fileData, err := io.ReadAll(file)
	if err != nil {
		logx.Errorw("read file repository failed",
			logx.Field("err", err))
		return resp, nil
	}

	// 视频文件提交到minio
	timeNow := time.Now()
	filename := utils.NewFileName(claim.UserId, timeNow.Unix())
	uploadInfo, err := minio.PutToBucketByBuf(
		l.ctx,
		minio.MinioVideoBucketName,
		filename+typ,
		bytes.NewBuffer(fileData),
	)
	if err != nil {
		logx.Errorw("upload file failed",
			logx.Field("err", err))
		return nil, err
	}
	playURL := minio.MinioVideoBucketName + "/" + filename + typ
	logx.Info("upload file success", uploadInfo)

	// 获取视频路径，并截取视频帧作为封面
	filepath, err := minio.Client.PresignedGetObject(l.ctx, minio.MinioVideoBucketName, filename+typ, time.Minute*1, nil)
	if err != nil {
		logx.Errorw("get object path failed",
			logx.Field("err", err))
		return nil, err
	}

	buf, err := ffmpeg.GetSnapshot(filepath.String())
	if err != nil || buf.Len() == 0 {
		logx.Errorw("get videoservice snapshot failed",
			logx.Field("err", err))
		return nil, err
	}

	// 将封面文件上传至minio
	coverURL := minio.MinioImgBucketName + "/" + filename + ".png"
	_, err = minio.PutToBucketByBuf(l.ctx, minio.MinioImgBucketName, filename+".png", buf)
	if err != nil {
		logx.Errorw("upload cover img failed",
			logx.Field("err", err))
		return nil, err
	}

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

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
