package job

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/db"
	"github.com/joeyz1729/zero-tiktok/pkg/worker"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func FavoriteCountWorker(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo) *worker.Worker {
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
			err = db.UpdateFavoriteCount(userId, incr, repo.DB)
			if err != nil {
				logx.Errorw("update favorite count", logx.Field("err", err),
					logx.Field("userId", userId))
				return err
			}
			logx.Infow("update favorite count success", logx.Field("incr", incr),
				logx.Field("userId", userId))
		}
		return nil
	}
	c.Topic = TopicFavor
	c.GroupID = GroupUpdateFavorCount
	return &worker.Worker{
		Handler:      handler,
		ReaderConfig: c,
	}

}

func TotalFavoritedWorkerStart(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo) error {
	c.Topic = TopicTotalFavorited
	c.GroupID = GroupUpdateTotalFavorited
	reader := kafka.NewReader(c)
	for {
		m, err := reader.FetchMessage(ctx)
		if errors.Is(err, context.Canceled) {
			return err
		}
		if err != nil {
			logx.Errorw("fetch message", logx.Field("err", err))
			break
		}
		userId, err := strconv.ParseInt(string(m.Key), 10, 64)
		if err != nil {
			return err
		}
		incr, err := strconv.ParseInt(string(m.Value), 10, 64)
		if err != nil {
			return err
		}
		if err := db.UpdateTotalFavorited(userId, incr, repo.DB); err != nil {
			return err
		}
		if err := reader.CommitMessages(ctx, m); err != nil {
			logx.Errorw("commit messages", logx.Field("err", err),
				logx.Field("authorId", userId))
			return err
		}
	}
	return nil
}
