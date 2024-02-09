package logic

import (
	"bytes"
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/internal/data"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/mw/ffmpeg"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/mw/minio"
	"time"

	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/model"
	"github.com/joeyz1729/zero-tiktok/pkg/snowflake"
	"github.com/joeyz1729/zero-tiktok/pkg/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishActionLogic) PublishAction(in *model.PublishActionRequest) (*model.PublishActionResponse, error) {
	resp := new(model.PublishActionResponse)
	var err error
	// snowflake gen video id
	vid, err := snowflake.GenID()

	// 视频文件提交到minio
	timeNow := time.Now()
	filename := utils.NewFileName(in.GetUserId(), timeNow.Unix())
	buffer := bytes.NewBuffer(in.Data)
	uploadInfo, err := minio.PutToBucketByBuf(
		l.ctx,
		minio.MinioVideoBucketName,
		filename+in.Type,
		buffer,
	)
	if err != nil {
		logx.Errorw("upload file failed",
			logx.Field("err", err))
		return nil, err
	}
	playURL := minio.MinioVideoBucketName + "/" + filename + in.Type
	logx.Info("upload file success", uploadInfo)

	// 获取视频路径，并截取视频帧作为封面
	filepath, err := minio.Client.PresignedGetObject(l.ctx, minio.MinioVideoBucketName, filename+in.Type, time.Minute*1, nil)
	if err != nil {
		logx.Errorw("get object path failed",
			logx.Field("err", err))
		return nil, err
	}

	buf, err := ffmpeg.GetSnapshot(filepath.String())
	if err != nil || buf.Len() == 0 {
		logx.Errorw("get video snapshot failed",
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

	err = l.svcCtx.VideoRepo.AddVideo(l.ctx, &data.Video{
		VideoId:  int64(vid),
		AuthorId: in.GetUserId(),
		Title:    in.GetTitle(),
		PlayUrl:  playURL,
		CoverUrl: coverURL,
	})

	if err != nil {
		logx.Errorw("[mysql] add video failed",
			logx.Field("err", err))
		return nil, err
	}
	resp.VideoId = int64(vid)
	return resp, nil

}
