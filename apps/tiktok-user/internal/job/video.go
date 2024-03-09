package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"time"
)

// VideoPublishStart 用户发布视频之后，更新work_count计数。
func (w *Worker) VideoPublishStart(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicVideo
	w.ReaderConfig.GroupID = TopicVideo
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
				err = w.updateWorkCount(ctx, data)
				if err != nil {
					logx.Errorf("update work count", logx.Field("err", err))
				}
			}
		} else {
		}
	}
	return nil
}

func (w *Worker) updateWorkCount(ctx context.Context, data map[string]interface{}) error {
	userId, err := strconv.ParseInt(data["id"].(string), 10, 64)
	if err != nil {
		return err
	}
	layout := "2006-01-02 15:04:05"
	createdTime, err := time.Parse(layout, data["create_time"].(string))
	if err != nil {
		return err
	}
	// todo
	return w.Repo.DBCreateCount(userId, createdTime)
}
