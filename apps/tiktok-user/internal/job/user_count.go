package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/es"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
)

func UserCountWorker(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		switch msg.Type {
		case "UPDATE":
		default:
		}
		for _, data := range msg.Data {
			err := es.UpdateUserCount(ctx, data, repo.ES)
			if err != nil {
				return err
			}
		}

		return nil
	}
	c.Topic = TopicUserCount
	c.GroupID = GroupUpdateUserCount
	return &worker.Worker{
		Handler:      handler,
		ReaderConfig: c,
	}

}
