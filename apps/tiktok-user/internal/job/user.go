package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/es"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// CreateUserStart 创建新用户时更新user_count表和es。
func (w *Worker) CreateUserStart(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicUser
	w.ReaderConfig.GroupID = GroupCreateUser
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
		if msg.Type == "INSERT" {
			for _, data := range msg.Data {
				w.EsCreateUser(ctx, data)
			}
		} else {

		}
	}
	return nil
}

func (w *Worker) EsCreateUser(ctx context.Context, data map[string]interface{}) {
	err := es.CreateUser(ctx, data, w.Repo.ES)
	if err != nil {
		logx.Errorw("es create user", logx.Field("err", err))
		return
	}
	logx.Infow("es create user success", logx.Field("data", data))
	return
}
