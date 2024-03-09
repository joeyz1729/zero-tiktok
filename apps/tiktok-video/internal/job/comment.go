package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
)

// UpdateCommentCount 更新视频对应的评价数
func (w *Worker) UpdateCommentCount(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicComment
	w.ReaderConfig.GroupID = TopicComment
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
			//for _, data := range msg.Data {
			//
			//}
		} else {

		}
	}
	return nil
}
