package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// CreateVideoStart 上传video并添加到video表之后，同步到video_count和es中。
func (w *Worker) CreateVideoStart(ctx context.Context) error {
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
		logx.Info(reader.Offset(), msg.Data)
		// todo 成功消费才更新offset
		if msg.Type == "INSERT" {
			for _, data := range msg.Data {
				err = w.Repo.CreateVideoCount(ctx, data)
				if err != nil {
					logx.Errorf("create video count", logx.Field("err", err))
				}
				err = w.Repo.EsCreateVideo(ctx, data)
				if err != nil {
					logx.Errorf("es create video", logx.Field("err", err))
				}
			}
		}
	}
	return nil
}
