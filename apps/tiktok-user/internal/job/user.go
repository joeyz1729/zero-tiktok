package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/es"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"time"
)

func (w *Worker) UserStart(ctx context.Context) error {
	for {
		w.ReaderConfig.Topic = TopicUser
		w.ReaderConfig.GroupID = TopicUser
		reader := kafka.NewReader(w.ReaderConfig)
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
			for idx := range msg.Data {
				err = w.insertCount(ctx, msg.Data[idx])
				if err != nil {
					logx.Error(err)
				}
				err = w.insertEs(ctx, msg.Data[idx])
				if err != nil {
					logx.Error(err)
				}
			}
		} else {

		}
	}
	return nil
}

func (w *Worker) CreateUser(ctx context.Context, data map[string]interface{}) error {
	err := w.insertCount(ctx, data)
	if err != nil {
		return err
	}
	err = w.insertEs(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) insertCount(ctx context.Context, data map[string]interface{}) error {
	userId, err := strconv.ParseInt(data["id"].(string), 10, 64)
	if err != nil {
		return err
	}
	layout := "2006-01-02 15:04:05"
	createdTime, err := time.Parse(layout, data["create_time"].(string))
	if err != nil {
		return err
	}
	return w.Repo.DBCreateCount(userId, createdTime)
}

func (w *Worker) insertEs(ctx context.Context, data map[string]interface{}) error {
	userId := data["id"].(string)
	resp, err := w.Repo.ES.Index(es.UserIndex).Id(userId).Document(data).Do(ctx)
	if err != nil {
		logx.Error(err)
		return err
	}
	logx.Info(resp)
	return nil
}
