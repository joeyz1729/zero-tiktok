package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/es"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// UserCountStart user_count表更新的时候，同步到es中
func (w *Worker) UserCountStart(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicUserCount
	w.ReaderConfig.GroupID = TopicUserCount
	reader := kafka.NewReader(w.ReaderConfig)
	for {
		m, err := reader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled) {
			return err
		}
		if err != nil {
			break
		}
		msg := new(Msg)
		if err := json.Unmarshal(m.Value, msg); err != nil {
			continue
		}
		if msg.Type == "UPDATE" {
			for idx := range msg.Data {
				w.updateCountES(ctx, msg.Data[idx])
			}
		} else {
		}
	}
	return nil
}

func (w *Worker) updateCountES(ctx context.Context, data map[string]interface{}) error {
	userId := data["id"].(string)
	resp, err := w.Repo.ES.Update(es.UserIndex, userId).Doc(data).Do(context.TODO())
	if err != nil {
		logx.Error(err)
		return err
	}
	logx.Info(resp)
	return nil
}
