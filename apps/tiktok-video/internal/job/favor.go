package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
)

// UpdateFavorCount 当点赞/取消点赞之后，更新对应视频的点赞信息。todo：需要同时更新video和user的计数，最好在favor中添加一个consumer分别投递到user和video
func (w *Worker) UpdateFavorCount(ctx context.Context) error {
	w.ReaderConfig.Topic = TopicFavor
	w.ReaderConfig.GroupID = TopicFavor
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
			for idx := range msg.Data {
				//w.updateCountES(ctx, msg.Data[idx])
			}
		} else {
		}
	}
	return nil
}

//func (w *Worker) updateCountES(ctx context.Context, data map[string]interface{}) error {
//	userId := data["id"].(string)
//	resp, err := w.Repo.ES.Update(es.UserIndex, userId).Doc(data).Do(context.TODO())
//	if err != nil {
//		logx.Error(err)
//		return err
//	}
//	logx.Info(resp)
//	return nil
//}
