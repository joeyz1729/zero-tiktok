package service

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/YiZou89/zero-tiktok/apps/comment/rmq/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	chanCount   = 5
	bufferCount = 1024
)

type Service struct {
	c        config.Config
	UserRpc  user.User
	VideoRpc video.Video

	addWaiter sync.WaitGroup
	delWaiter sync.WaitGroup
	addChan   []chan *AddData
	delChan   []chan *DelData
}

type AddData struct {
	UserId      int64  `json:"user_id"`
	VideoId     int64  `json:"video_id"`
	CommentText string `json:"comment_text"`
}

type DelData struct {
	UserId    int64 `json:"user_id"`
	VideoId   int64 `json:"video_id"`
	CommentId int64 `json:"comment_id"`
}

func NewService(c config.Config) *Service {
	s := &Service{
		c:        c,
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc: video.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		addChan:  make([]chan *AddData, chanCount),
		delChan:  make([]chan *DelData, chanCount),
	}

	for i := 0; i < chanCount; i++ {
		ch1 := make(chan *AddData, bufferCount)
		ch2 := make(chan *DelData, bufferCount)
		s.addChan[i] = ch1
		s.delChan[i] = ch2
		s.addWaiter.Add(1)
		s.delWaiter.Add(1)
		go s.consumeAdd(ch1)
		go s.consumeDel(ch2)
	}
	return s
}

func (s *Service) consumeAdd(ch chan *AddData) {
	defer s.addWaiter.Done()
	for {
		m, ok := <-ch
		if !ok {
			log.Fatal("add comment mq exit")
		}
		logx.Infof("add comment consume msg: %v\n", m)

	}
}

func (s *Service) consumeDel(ch chan *DelData) {
	defer s.delWaiter.Done()
	for {
		m, ok := <-ch
		if !ok {
			log.Fatal("del comment mq exit")
		}
		logx.Infof("del comment consume msg: %v\n", m)
		// biz block

	}
}

//func (s *Service) consumeDTM(ch chan *KafkaData) {
//
//}

func (s *Service) ConsumeAdd(_ string, value string) error {
	logx.Infof("Consume value: %s\n", value)
	var data []*AddData
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}
	// gen comment id

	for _, d := range data {
		s.addChan[d.VideoId&chanCount] <- d
	}
	return nil
}

func (s *Service) ConsumeDel(_ string, value string) error {
	logx.Infof("Consume value: %s\n", value)
	var data []*DelData
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}
	for _, d := range data {
		s.delChan[d.CommentId&chanCount] <- d
	}
	return nil
}
