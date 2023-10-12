package logic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/dao"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/mw/ffmpeg"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/mw/minio"
	"time"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/YiZou89/zero-tiktok/pkg/snowflake"
	"github.com/YiZou89/zero-tiktok/pkg/utils"
	"github.com/zeromicro/go-zero/core/stores/redis"

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
	// todo: add your logic here and delete this line
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
		return resp, err
	}
	playURL := minio.MinioVideoBucketName + "/" + filename + in.Type
	logx.Info("upload file success", uploadInfo)

	// 获取视频路径，并截取视频帧作为封面
	filepath, err := minio.Client.PresignedGetObject(l.ctx, minio.MinioVideoBucketName, filename+in.Type, time.Minute*1, nil)
	if err != nil {
		logx.Errorw("get object path failed",
			logx.Field("err", err))
		return resp, err
	}

	buf, err := ffmpeg.GetSnapshot(filepath.String()) //TODO
	if err != nil || buf.Len() == 0 {
		logx.Errorw("get video snapshot failed",
			logx.Field("err", err))
		return resp, err
	}

	// 将封面文件上传至minio
	coverURL := minio.MinioImgBucketName + "/" + filename + ".png"
	_, err = minio.PutToBucketByBuf(l.ctx, minio.MinioImgBucketName, filename+".png", buf)
	if err != nil {
		logx.Errorw("upload cover img failed",
			logx.Field("err", err))
		return resp, err
	}

	// add video
	// 添加到redis，用于发布列表，视频流，用户work count统计
	pipeline := l.svcCtx.VideoCache.TxPipeline()
	pipeline.ZAdd(l.ctx, "tiktok:video:time", &redis.Z{Score: float64(time.Now().Unix() - time.Date(2023, time.September, 1, 1, 46, 40, 0, time.UTC).Unix()), Member: vid})
	pipeline.SAdd(l.ctx, "tiktok:video:user:"+fmt.Sprintf("%d", in.UserId), vid)
	_, err = pipeline.Exec(l.ctx)
	if err != nil {
		logx.Errorw("[redis] add redis failed",
			logx.Field("err", err))
		return resp, err
	}
	logx.Info("[redis] add redis success")

	// todo, 异步入库
	go func() {

		err = l.svcCtx.VideoRepo.AddVideo(l.ctx, &dao.Video{
			VideoId:  int64(vid),
			AuthorId: in.GetUserId(),
			Title:    in.GetTitle(),
			PlayUrl:  playURL,
			CoverUrl: coverURL,
		})
		//_, err = l.svcCtx.VideoModel.Insert(context.Background(), &model.Video{
		//	VideoId:     int64(vid),
		//	AuthorId:    in.GetUserId(),
		//	Title:       in.GetTitle(),
		//	PlayUrl:     playURL,
		//	CoverUrl:    coverURL,
		//	PublishTime: timeNow,
		//})
		if err != nil {
			logx.Errorw("[mysql] add video failed",
				logx.Field("err", err))
		} else {
			logx.Info("[mysql] add video success")
		}
	}()

	resp.VideoId = int64(vid)
	return resp, nil

}
