package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func RelationWorker(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo) *worker.Worker {
	handler := func(msg *worker.Msg) error {
		var incr int64
		switch msg.Type {
		case "INSERT":
			incr = 1
		case "DELETE":
			incr = -1
		default:
			return nil
		}
		for _, data := range msg.Data {
			userId, err := strconv.ParseInt(data["user_id"].(string), 10, 64)
			if err != nil {
				return err
			}
			toUserId, err := strconv.ParseInt(data["followed_id"].(string), 10, 64)
			if err != nil {
				return err
			}
			err = db.UpdateRelationCount(userId, toUserId, incr, repo.DB)
			if err != nil {
				logx.Errorw("update relation count", logx.Field("err", err),
					logx.Field("userId", userId), logx.Field("toUserId", toUserId))
				return err
			}
			logx.Infow("update relation count success", logx.Field("incr", incr),
				logx.Field("userId", userId), logx.Field("toUserId", toUserId))
		}
		return nil
	}
	c.Topic = TopicRelation
	c.GroupID = GroupUpdateRelationCount
	return &worker.Worker{
		Handler:      handler,
		ReaderConfig: c,
	}

}
