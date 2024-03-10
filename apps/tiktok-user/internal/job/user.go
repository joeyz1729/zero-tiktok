package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/es"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
)

func CreateUserWorker(ctx context.Context, c config.KafkaConfig, repo *repository.Repo) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		if msg.Type == "INSERT" {
			for _, data := range msg.Data {
				err := es.CreateUser(ctx, data, repo.ES)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return &worker.Worker{
		Handler: handler,
		ReaderConfig: kafka.ReaderConfig{
			Brokers: c.Brokers,
			Topic:   TopicUser,
			GroupID: GroupCreateUser,
			// todo
		},
	}

}
