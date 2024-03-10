package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/db"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func UpdateFavorCount(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo, w *kafka.Writer) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		var incr int64
		switch msg.Type {
		case "INSERT":
			incr = 1
		case "DELETE":
			incr = -1
		default:
			return nil
		}
		for _, data := range msg.Data {
			videoId, err := strconv.ParseInt(data["video_id"].(string), 10, 64)
			if err != nil {
				logx.Errorw("parse video id from message data", logx.Field("err", err),
					logx.Field("data", data))
				return err
			}
			err = db.UpdateThumbupCount(videoId, incr, repo.DB)
			if err != nil {
				logx.Errorw("update thumbup count", logx.Field("err", err),
					logx.Field("videoId", videoId))
				return err
			}
			logx.Infow("update thumbup count success", logx.Field("videoId", videoId))

			video, err := repo.GetVideoById(ctx, videoId)
			if err != nil {
				logx.Errorw("get video by id", logx.Field("err", err),
					logx.Field("videoId", videoId))
				return err
			}
			err = w.WriteMessages(ctx,
				kafka.Message{
					Key:   []byte(strconv.FormatInt(video.AuthorID, 10)),
					Value: []byte(strconv.FormatInt(incr, 10)),
				})
			if err != nil {
				logx.Errorw("send message", logx.Field("err", err),
					logx.Field("authorId", video.AuthorID))
				return err
			}
		}
		return nil
	}

	c.Topic = TopicFavor
	c.GroupID = GroupUpdateFavorCount
	return &worker.Worker{
		Handler:      handler,
		ReaderConfig: c,
	}

}
