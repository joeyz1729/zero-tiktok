package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func VideoPublishWorker(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		switch msg.Type {
		case "INSERT":
		default:
		}
		for _, data := range msg.Data {
			userId, err := strconv.ParseInt(data["author_id"].(string), 10, 64)
			if err != nil {
				return err
			}
			err = db.UpdateWorkCount(userId, 1, repo.DB)
			if err != nil {
				return err
			}
			logx.Infow("update work count success", logx.Field("data", data))
		}
		return nil
	}
	c.Topic = TopicVideo
	c.GroupID = GroupUpdateWorkCount
	return &worker.Worker{
		Handler:      handler,
		ReaderConfig: c,
	}

}
