package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/es"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// VideoPublishStart 上传视频时事务更新video和video_count，通过mq同步canal。
func (w *Worker) VideoPublishStart(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicVideo
	w.ReaderConfig.GroupID = GroupCreateVideo
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
				w.EsCreateVideo(ctx, data)
			}
		}
	}
	return nil
}

func (w *Worker) EsCreateVideo(ctx context.Context, data map[string]interface{}) {
	err := es.CreateVideo(ctx, data, w.Repo.ES)
	if err != nil {
		logx.Errorw("es create video", logx.Field("err", err))
	}
	logx.Infow("es create video success", logx.Field("data", data))
	return
}
