package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/data"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelLogic {
	return &DelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelLogic) Del(in *model.DelRequest) (*model.DelResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.DelResponse)
	var err error

	mid, fid := strconv.Itoa(int(in.UserId)), strconv.Itoa(int(in.ToUserId))

	pipeline := l.svcCtx.FollowCache.TxPipeline()
	_, err = pipeline.ZRem(l.ctx, data.FollowingPrefix+mid, fid).Result()
	_, err = pipeline.ZRem(l.ctx, data.FollowingPrefix+fid, mid).Result()
	_, err = pipeline.Decr(l.ctx, data.FollowingCountPerfix+mid).Result()
	_, err = pipeline.Decr(l.ctx, data.FollowerCountPrefix+fid).Result()
	_, err = pipeline.Exec(l.ctx)
	if err != nil {
		logx.Errorw("[redis] del failed",
			logx.Field("err", err),
		)
		resp.Code = http.StatusInternalServerError
		resp.Msg = "del redis failed"
		return resp, err
	}
	logx.Info("[redis] del success")

	go func() {
		// TODO add rabbitmq
		logx.Info("[mysql] asynchronous del from database")
		sqlStr := `delete from  tiktok_follow.follow where user_id = ? and follower_id = ?`
		_, err = l.svcCtx.FollowDB.Exec(sqlStr, in.UserId, in.ToUserId)
		if err != nil {
			logx.Errorw("[mysql] del failed",
				logx.Field("err", err),
			)
		}
	}()

	resp.Code = http.StatusOK
	resp.Msg = "success"
	return resp, nil

}
