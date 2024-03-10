package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

// VideoPublishStart 上传视频时事务更新video和video_count，通过mq同步canal。
func (w *Worker) VideoPublishStart(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicVideo
	w.ReaderConfig.GroupID = GroupUpdateWorkCount
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
				w.updateWorkCount(data)
			}
		}
	}
	return nil
}

func (w *Worker) updateWorkCount(data map[string]interface{}) {
	userId, err := strconv.ParseInt(data["author_id"].(string), 10, 64)
	if err != nil {
		logx.Errorw("get author_id from msg data", logx.Field("err", err))
		return
	}

	err = db.UpdateWorkCount(userId, 1, w.Repo.DB)
	if err != nil {
		logx.Errorw("db update work count", logx.Field("err", err))
		return
	}
	logx.Infow("update work count success", logx.Field("data", data))
}
