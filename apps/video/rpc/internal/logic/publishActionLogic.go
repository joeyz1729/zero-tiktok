package logic

import (
	"bytes"
	"context"
	"fmt"
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
	// gen url
	// 同步加入redis和异步加入rabbitmq修改mysql，
	// 加入redis后返回
	_, err = l.svcCtx.VideoCache.ZAdd(l.ctx, "tiktok:video:time", &redis.Z{Score: float64(time.Now().Unix() - time.Date(2023, time.September, 1, 1, 46, 40, 0, time.UTC).Unix()), Member: vid}).Result()
	if err != nil {
		logx.Errorw("[redis] add video time zset failed",
			logx.Field("err", err))
		return resp, nil
	}
	logx.Info("[redis] add video time zset success")
	resp.VideoId = int64(vid)
	logx.Info("[mysql] asynchronous add video into database")

	// 读取file信息并修改filename
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
	logx.Info("upload file success", uploadInfo)
	playURL := minio.MinioVideoBucketName + "/" + filename + in.Type
	filepath, err := minio.Client.PresignedGetObject(l.ctx, minio.MinioVideoBucketName, filename+in.Type, time.Minute*1, nil)
	if err != nil {
		logx.Errorw("get object path failed",
			logx.Field("err", err))
		return resp, err
	}
	fmt.Printf("get object path success, %s\n", filepath.String())

	buf, err := ffmpeg.GetSnapshot(filepath.String()) //TODO
	if err != nil || buf.Len() == 0 {
		logx.Errorw("get video snapshot failed",
			logx.Field("err", err))
		return resp, err
	}
	logx.Infof("video cover snapshot size: %d\n", buf.Len())
	upInfo, err := minio.PutToBucketByBuf(l.ctx, minio.MinioImgBucketName, filename+".png", buf)
	if err != nil {
		logx.Errorw("upload cover img failed",
			logx.Field("err", err))
		return resp, err
	}
	logx.Infof("upload video cover success, size: %d\n", upInfo.Size)

	go func() {
		_, err = l.svcCtx.VideoModel.Insert(context.Background(), &model.Video{
			VideoId:     int64(vid),
			AuthorId:    in.GetUserId(),
			Title:       in.GetTitle(),
			PlayUrl:     playURL,
			CoverUrl:    minio.MinioImgBucketName + "/" + filename + ".png",
			PublishTime: timeNow,
		})
		if err != nil {
			logx.Errorw("[mysql] add video failed",
				logx.Field("err", err))
		} else {
			logx.Info("[mysql] add video success")
		}
	}()

	return resp, nil

}
