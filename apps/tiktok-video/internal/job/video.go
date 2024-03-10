package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/es"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

func VideoPublish(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		switch msg.Type {
		case "INSERT":
		default:
			return nil
		}
		for _, data := range msg.Data {
			err := es.CreateVideo(ctx, data, repo.ES)
			if err != nil {
				logx.Errorw("es create video", logx.Field("err", err))
				return err
			}
			logx.Infow("es create video success", logx.Field("data", data))
		}
		return nil
	}

	c.Topic = TopicVideo
	c.GroupID = GroupCreateVideo
	return &worker.Worker{
		Handler:      handler,
		ReaderConfig: c,
	}

}
