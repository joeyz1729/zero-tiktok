package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// UserCountStart user_count表更新的时候，同步到es中
func (w *Worker) UserCountStart(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicUserCount
	w.ReaderConfig.GroupID = GroupUpdateUserCount
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
			for _, data := range msg.Data {
				w.EsUpdateUserCount(ctx, data)
			}
		} else {
		}
	}
	return nil
}

func (w *Worker) EsUpdateUserCount(ctx context.Context, data map[string]interface{}) {
	err := w.Repo.EsUpdateUserCount(ctx, data)
	if err != nil {
		logx.Errorw("es update user count", logx.Field("err", err))
		return
	}
	logx.Infow("es update user count", logx.Field("data", data))
	return
}
