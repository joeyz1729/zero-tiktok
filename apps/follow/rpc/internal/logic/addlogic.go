package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/data"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddLogic) Add(in *model.AddRequest) (*model.AddResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.AddResponse)
	var err error

	// add into bloom filter
	mid, fid := strconv.Itoa(int(in.UserId)), strconv.Itoa(int(in.ToUserId))
	bloomKey := data.BloomPrefix + mid + ":" + fid
	err = l.svcCtx.Filter.AddCtx(l.ctx, []byte(bloomKey))
	if err != nil {
		logx.Errorw("[bloom filter] add failed",
			logx.Field("err", err),
		)
		resp.Code = http.StatusInternalServerError
		resp.Msg = "add bloom filter failed"
		return resp, err
	}
	logx.Info("add bloom filter success")

	pipeline := l.svcCtx.FollowCache.TxPipeline()
	_, err = pipeline.ZAdd(l.ctx, data.FollowingPrefix+mid, &redis.Z{Score: float64(time.Now().Unix()), Member: fid}).Result()
	_, err = pipeline.ZAdd(l.ctx, data.FollowerPrefix+fid, &redis.Z{Score: float64(time.Now().Unix()), Member: mid}).Result()
	_, err = pipeline.Incr(l.ctx, data.FollowingCountPerfix+mid).Result()
	_, err = pipeline.Incr(l.ctx, data.FollowerCountPrefix+fid).Result()
	_, err = pipeline.Exec(l.ctx)
	if err != nil {
		logx.Errorw("[redis] pipeline failed",
			logx.Field("err", err),
		)
		resp.Code = http.StatusInternalServerError
		resp.Msg = "add redis failed"
		return resp, err
	}
	logx.Info("[redis] add success")
	resp.Code = http.StatusOK
	resp.Msg = "success"

	go func() {
		// TODO add rabbitmq
		logx.Info("[mysql] asynchronous add database")
		sqlStr := `insert into tiktok_follow.follow(user_id, follower_id) value(?, ?)`
		_, err = l.svcCtx.FollowDB.Exec(sqlStr, in.UserId, in.ToUserId)
		if err != nil {
			logx.Errorw("[mysql] add failed",
				logx.Field("err", err),
			)
		}
	}()

	return resp, nil

}
