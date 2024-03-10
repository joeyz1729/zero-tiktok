package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"strconv"
)

func incrCount(ctx context.Context, data map[string]interface{}, DB *gorm.DB) error {
	userId, err := strconv.ParseInt(data["user_id"].(string), 10, 64)
	if err != nil {
		return err
	}
	toUserId, err := strconv.ParseInt(data["followed_id"].(string), 10, 64)
	if err != nil {
		return err
	}
	return db.UpdateRelationCount(userId, toUserId, 1, DB)
}

func decrCount(ctx context.Context, data map[string]interface{}, DB *gorm.DB) error {
	userId, err := strconv.ParseInt(data["user_id"].(string), 10, 64)
	if err != nil {
		return err
	}
	toUserId, err := strconv.ParseInt(data["followed_id"].(string), 10, 64)
	if err != nil {
		return err
	}
	return db.UpdateRelationCount(userId, toUserId, -1, DB)
}

func RelationWorker(ctx context.Context, c config.KafkaConfig, repo *repository.Repo) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		if msg.Type == "INSERT" {
			for idx := range msg.Data {
				if err := incrCount(ctx, msg.Data[idx], repo.DB); err != nil {
					return err
				}
			}
		} else if msg.Type == "DELETE" {
			for idx := range msg.Data {
				if err := decrCount(ctx, msg.Data[idx], repo.DB); err != nil {
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
			Topic:   TopicRelation,
			GroupID: GroupUpdateRelationCount,
			// todo
		},
	}

}
