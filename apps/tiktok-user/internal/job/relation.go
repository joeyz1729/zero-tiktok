package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

// RelationStart 关注动作之后，更新双方用户的计数信息。
func (w *Worker) RelationStart(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicRelation
	w.ReaderConfig.GroupID = TopicRelation
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
		if msg.Type == "INSERT" {
			for idx := range msg.Data {
				err = w.incrCount(ctx, msg.Data[idx])
				if err != nil {
					logx.Error(err)
				}
			}
		} else if msg.Type == "DELETE" {
			for idx := range msg.Data {
				err = w.decrCount(ctx, msg.Data[idx])
				if err != nil {
					logx.Error(err)
				}
			}
		}
	}
	return nil
}

func (w *Worker) incrCount(ctx context.Context, data map[string]interface{}) error {
	userId, err := strconv.ParseInt(data["user_id"].(string), 10, 64)
	if err != nil {
		return err
	}
	toUserId, err := strconv.ParseInt(data["followed_id"].(string), 10, 64)
	if err != nil {
		return err
	}
	logx.Infow("add follow", logx.Field("user_id", userId), logx.Field("toUserId", toUserId))
	return w.Repo.DBUpdateRelationCount(userId, toUserId, 1)
}

func (w *Worker) decrCount(ctx context.Context, data map[string]interface{}) error {
	userId, err := strconv.ParseInt(data["user_id"].(string), 10, 64)
	if err != nil {
		return err
	}
	toUserId, err := strconv.ParseInt(data["followed_id"].(string), 10, 64)
	if err != nil {
		return err
	}
	logx.Infow("del follow", logx.Field("user_id", userId), logx.Field("toUserId", toUserId))
	return w.Repo.DBUpdateRelationCount(userId, toUserId, -1)
}
