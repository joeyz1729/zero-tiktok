package logic

import (
	"context"
	"database/sql"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/YiZou89/zero-tiktok/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishActionLogic) PublishAction(in *model.PublishActionRequest) (*model.PublishActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.PublishActionResponse)
	var err error
	// snowflake gen video id
	vid, err := snowflake.GenID()
	// gen url
	// 同步加入redis和异步加入rabbitmq修改mysql，
	// 加入redis后返回
	_, err = l.svcCtx.VideoCache.ZAdd(l.ctx, "tiktok:video:time", &redis.Z{Score: float64(time.Now().Unix() - time.Date(2023, time.September, 1, 1, 46, 40, 0, time.UTC).Unix()), Member: vid}).Result()
	if err != nil {
		logx.Errorw("[redis] add video time zset failed",
			logx.Field("err", err))
		return resp, nil
	}
	logx.Info("[redis] add video time zset success")
	resp.VideoId = int64(vid)
	logx.Info("[mysql] asynchronous add video into database")
	go func() {
		_, err = l.svcCtx.VideoModel.Insert(context.Background(), &model.Video{
			VideoId:  int64(vid),
			AuthorId: in.GetUserId(),
			Title:    in.GetTitle(),
			Data: sql.NullString{
				String: string(in.GetData()),
				Valid:  true,
			},
			PlayUrl:     "need cover url",
			CoverUrl:    "need cover url",
			PublishTime: time.Now(),
		})
		if err != nil {
			logx.Errorw("[mysql] add video failed",
				logx.Field("err", err))
		} else {
			logx.Info("[mysql] add video success")
		}
	}()

	return resp, nil

}
