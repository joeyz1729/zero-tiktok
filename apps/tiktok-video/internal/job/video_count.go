package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/es"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

func VideoCountWorker(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		switch msg.Type {
		case "UPDATE":
		default:
			return nil
		}
		for _, data := range msg.Data {
			err := es.UpdateVideoCount(ctx, data, repo.ES)
			if err != nil {
				logx.Errorw("es update video count", logx.Field("err", err))
				return err
			}
			logx.Infow("es update video count success", logx.Field("data", data))
		}
		return nil
	}
	c.Topic = TopicVideoCount
	c.GroupID = GroupUpdateVideoCount
	return &worker.Worker{
		Handler:      handler,
		ReaderConfig: c,
	}

}
