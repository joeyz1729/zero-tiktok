package job

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/repository"
	"github.com/segmentio/kafka-go"
)

type Worker struct {
	Repo        *repository.Repo
	KafkaReader *kafka.Reader
}

//func (w *Worker) Start(ctx context.Context) error {
//	for {
//		m, err := w.KafkaReader.ReadMessage(ctx)
//		if errors.Is(err, context.Canceled) {
//			return err
//		}
//		if err != nil {
//			break
//		}
//		msg := new(Msg)
//		if err := json.Unmarshal(m.Value, msg); err != nil {
//			continue
//		}
//		if msg.Type == "INSERT" {
//			// 添加关注
//			for idx := range msg.Data {
//				err = w.insertCount(ctx, msg.Data[idx])
//				if err != nil {
//					logx.Error(err)
//					continue
//				}
//				err = w.insertEs(ctx, msg.Data[idx])
//				if err != nil {
//					logx.Error(err)
//				}
//			}
//		} else if msg.Type == "DELETE" {
//			// 取消关注
//
//		}
//	}
//	return nil
//}
//
//type Msg struct {
//	Type     string                   `json:"type"`
//	Database string                   `json:"database"`
//	Table    string                   `json:"table"`
//	IsDdl    bool                     `json:"isDdl"`
//	Data     []map[string]interface{} `json:"data"`
//}
//
//func (w *Worker) incrCount(ctx context.Context, data map[string]interface{}) error {
//	userId, err := strconv.ParseInt(data["id"].(string), 10, 64)
//	if err != nil {
//		return err
//	}
//	layout := "2006-01-02 15:04:05"
//	createdTime, err := time.Parse(layout, data["create_time"].(string))
//	if err != nil {
//		return err
//	}
//	return w.Repo.DBCreateCount(userId, createdTime)
//}
//
//func (w *Worker) decrCount(ctx context.Context, data map[string]interface{}) error {
//	userId, err := strconv.ParseInt(data["id"].(string), 10, 64)
//	if err != nil {
//		return err
//	}
//	layout := "2006-01-02 15:04:05"
//	createdTime, err := time.Parse(layout, data["create_time"].(string))
//	if err != nil {
//		return err
//	}
//	return w.Repo.DBCreateCount(userId, createdTime)
//}
//
////func (w *Worker) insertEs(ctx context.Context, data map[string]interface{}) error {
////	userId := data["id"].(string)
////	resp, err := w.Repo.ES.Index(es.UserIndex).Id(userId).Document(data).Do(ctx)
////	if err != nil {
////		logx.Error(err)
////		return err
////	}
////	logx.Info(resp)
////	return nil
////}
////
////// DecrCount  取消关注时，减少双方的计数
////func (w *Worker) DecrCount(userId, toUserId int64) error {
////	// 修改数据库数据，删除缓存数据
////
////	return
////}
////
////// IncrCount 新增关注时，增加双方的计数
////func (w *Worker) IncrCount(userId, toUserId int64) error {
////
////}
