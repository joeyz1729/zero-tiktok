package logic

import (
	"context"
	"database/sql"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/data"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"

	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationLogic {
	return &GetRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationLogic) GetRelation(in *model.GetRelationRequest) (*model.GetRelationResponse, error) {
	// todo: add your logic here and delete this line
	var resp = new(model.GetRelationResponse)
	var err error
	mid := strconv.Itoa(int(in.UserId))
	fid := strconv.Itoa(int(in.ToUserId))
	logx.Info("begin get follow relation")
	// 1. bloom filter
	bloomKey := data.BloomPrefix + mid + ":" + fid
	ok, err := l.svcCtx.Filter.ExistsCtx(l.ctx, []byte(bloomKey))
	if err != nil {
		logx.Errorw("bloom filter failed",
			logx.Field("err", err),
		)
		return resp, err
	}
	if !ok {
		logx.Info("bloom filter not exist, return")
		resp.IfFollowing = 0
		resp.IfFollower = 11111
		return resp, nil
	}
	// 2. check redis
	logx.Info("bloom filter check success, check redis cache")
	redisKey := data.FollowingPrefix + mid
	_, err = l.svcCtx.FollowCache.ZScore(l.ctx, redisKey, fid).Result()
	// err
	if err != nil && err != redis.Nil {
		logx.Errorw("redis get follow relation failed",
			logx.Field("err", err),
		)
		return resp, err
	}
	// redis success
	if err == nil {
		logx.Info("redis get follow relation success")
		resp.IfFollowing = 1
		return resp, nil
	}
	// not exist
	resp.IfFollowing = 0

	// miss
	sqlStr := `select id from tiktok_follow.follow where user_id = ? and follower_id = ?`
	_, err = l.svcCtx.FollowDB.Query(sqlStr, in.UserId, in.ToUserId)
	// failed
	if err != nil && err != sql.ErrNoRows {
		logx.Errorw("mysql query failed",
			logx.Field("err", err),
		)
		return resp, err
	}

	if err == sql.ErrNoRows {
		logx.Info("mysql does not exist")
		resp.IfFollowing = 0
		return resp, nil
	}
	// mysql success get
	resp.IfFollowing = 1
	go func() {
		// asynchronous add redis
		// TODO 异步处理， 消息队列，错误处理，或者刷新时间
		pipeline := l.svcCtx.FollowCache.TxPipeline()
		pipeline.ZAdd(l.ctx, data.FollowingPrefix+mid, &redis.Z{Score: float64(time.Now().Unix()), Member: fid}).Result()
		pipeline.ZAdd(l.ctx, data.FollowerPrefix+fid, &redis.Z{Score: float64(time.Now().Unix()), Member: mid}).Result()
		_, err := pipeline.Exec(l.ctx)
		if err != nil {
			logx.Errorw("[redis] asynchronous add failed",
				logx.Field("err", err),
			)
		}
		logx.Info("[redis] asynchronous add success")
	}()

	return resp, nil

	//return &model.GetRelationResponse{}, nil
}
