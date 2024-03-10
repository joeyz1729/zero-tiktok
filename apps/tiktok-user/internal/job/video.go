package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"strconv"
)

func VideoPublishWorker(ctx context.Context, c config.KafkaConfig, repo *repository.Repo) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		if msg.Type == "INSERT" {
			for _, data := range msg.Data {
				updateWorkCount(data, repo.DB)
			}
		}
		return nil
	}
	return &worker.Worker{
		Handler: handler,
		ReaderConfig: kafka.ReaderConfig{
			Brokers: c.Brokers,
			Topic:   TopicVideo,
			GroupID: GroupUpdateWorkCount,
			// todo
		},
	}

}

func updateWorkCount(data map[string]interface{}, DB *gorm.DB) {
	userId, err := strconv.ParseInt(data["author_id"].(string), 10, 64)
	if err != nil {
		logx.Errorw("get author_id from msg data", logx.Field("err", err))
		return
	}

	err = db.UpdateWorkCount(userId, 1, DB)
	if err != nil {
		logx.Errorw("db update work count", logx.Field("err", err))
		return
	}
	logx.Infow("update work count success", logx.Field("data", data))
}
