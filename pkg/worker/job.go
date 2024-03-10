package worker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type Worker struct {
	Handler      func(*Msg) error
	ReaderConfig kafka.ReaderConfig
}

type Msg struct {
	Type     string                   `json:"type"`
	Database string                   `json:"database"`
	Table    string                   `json:"table"`
	IsDdl    bool                     `json:"isDdl"`
	Data     []map[string]interface{} `json:"data"`
}

func (w *Worker) Start(ctx context.Context) error {
	reader := kafka.NewReader(w.ReaderConfig)
	for {
		m, err := reader.FetchMessage(ctx)
		if errors.Is(err, context.Canceled) {
			return err
		}
		if err != nil {
			logx.Errorw("fetch message", logx.Field("err", err))
			break
		}
		msg := new(Msg)
		if err := json.Unmarshal(m.Value, msg); err != nil {
			logx.Errorw("json unmarshal", logx.Field("m.Value", string(m.Value)))
			continue
		}
		if err = w.Handler(msg); err != nil {
			logx.Errorw("consume messages", logx.Field("err", err), logx.Field("msg", msg))
		}
		if err := reader.CommitMessages(ctx, m); err != nil {
			logx.Errorw("commit messages", logx.Field("err", err), logx.Field("msg", msg))
		}
	}
	return nil
}
