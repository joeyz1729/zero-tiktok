package kmq

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	chanCount   = 5
	bufferCount = 1024
)

type Kmq struct {
	c      kq.KqConf
	db     *sqlx.DB
	waiter sync.WaitGroup
	//addWaiter sync.WaitGroup
	//delWaiter sync.WaitGroup
	dataChan []chan *KafkaData
	//addChan   []chan *AddData
	//delChan   []chan *DelData
}

type KafkaData struct {
	ActionType  bool   `json:"action_type"`
	UserId      int64  `json:"user_id"`
	VideoId     int64  `json:"video_id"`
	CommentText string `json:"comment_text"`
	CommentId   int64  `json:"comment_id"`
}

//type AddData struct {
//	UserId      int64  `json:"user_id"`
//	VideoId     int64  `json:"video_id"`
//	CommentText string `json:"comment_text"`
//}
//
//type DelData struct {
//	UserId    int64 `json:"user_id"`
//	VideoId   int64 `json:"video_id"`
//	CommentId int64 `json:"comment_id"`
//}

func NewMq(c kq.KqConf, db *sqlx.DB) *Kmq {
	s := &Kmq{
		c:        c,
		dataChan: make([]chan *KafkaData, chanCount*2),
		db:       db,
		//addChan: make([]chan *AddData, chanCount),
		//delChan: make([]chan *DelData, chanCount),
	}

	for i := 0; i < chanCount*2; i++ {
		ch := make(chan *KafkaData, bufferCount)
		s.dataChan[i] = ch
		s.waiter.Add(1)
		go s.consume(ch)
		//ch1 := make(chan *AddData, bufferCount)
		//ch2 := make(chan *DelData, bufferCount)
		//s.addChan[i] = ch1
		//s.delChan[i] = ch2
		//s.addWaiter.Add(1)
		//s.delWaiter.Add(1)
		//go s.consumeAdd(ch1)
		//go s.consumeDel(ch2)
	}
	logx.Info("init kafka mq success")
	return s
}

func (s *Kmq) consume(ch chan *KafkaData) {
	defer s.waiter.Done()
	for {
		in, ok := <-ch
		if !ok {
			log.Fatal("add comment mq exit")
		}
		logx.Infof("add comment consume msg: %v\n", in)
		if in.ActionType {
			if in.CommentId == 0 {
				return
			}

			logx.Info("start insert into comment database")
			sqlStr := `insert into tiktok_comment.comment(video_id, user_id, comment_id, content) value(?, ?, ?, ?)`
			_, err := s.db.Exec(sqlStr, in.VideoId, in.UserId, in.CommentId, in.CommentText)
			if err != nil {
				logx.Errorw("mysql insert comment record failed",
					logx.Field("err", err),
				)
				//resp.Msg = err.Error()
			}
		} else {
			// del action
			delStr := `delete from tiktok_comment.comment where user_id = ? and video_id = ? and comment_id = ?`
			_, err := s.db.Exec(delStr, in.UserId, in.VideoId, in.CommentId)
			if err != nil {
				logx.Errorw("del comment record failed",
					logx.Field("err", err),
				)
			}
		}

	}
}

func (s *Kmq) Consume(_ string, value string) error {
	logx.Infof("Consume value: %s\n", value)
	var data KafkaData
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		logx.Errorw("unmarshal kafka data failed",
			logx.Field("err", err))
		return err
	}
	logx.Info("json unmarshal success")
	//for _, d := range data {

	select {
	case s.dataChan[data.CommentId&chanCount] <- &data:
		logx.Info("add kafka data to comment chan success")
		return nil
	case <-time.After(time.Second):
		logx.Errorf("add kafka data timeout, channel len: %d\n", len(s.dataChan[data.CommentId&chanCount]))
		return errors.New("add mq timeout")
	}

}
