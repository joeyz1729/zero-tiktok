package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/dao/cache"
	"github.com/YiZou89/zero-tiktok/apps/follow/model"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strconv"
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
	resp := new(model.DelResponse)
	var err error
	uid, tid := in.UserId, in.ToUserId
	uidStr, tidStr := strconv.FormatInt(in.UserId, 10), strconv.FormatInt(in.UserId, 10)

	tx, err := l.svcCtx.FollowDB.Beginx()
	if err != nil {
		logx.Errorw("[mysql] begin transaction failed",
			logx.Field("err", err))
		return resp, err
	}

	sqlStr := `delete from tiktok_follow.followed where user_id = ? and followed_id = ?`
	_, err = tx.Exec(sqlStr, in.UserId, in.ToUserId)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logx.Errorw("del followed failed, and rollback failed")
			return resp, err
		}
		logx.Errorw("del followed failed, rollback success")
		return resp, err
	}
	sqlStr = `delete from tiktok_follow.follower where follower_id = ? and user_id = ?`
	_, err = tx.Exec(sqlStr, in.ToUserId, in.UserId)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logx.Errorw("del follower failed, and rollback failed")
			return resp, err
		}
		logx.Errorw("del follower failed, rollback success")
		return resp, err
	}
	if err = tx.Commit(); err != nil {
		logx.Errorw("transaction commit failed")
		return resp, err
	}

	pipeline := l.svcCtx.FollowCache.TxPipeline()
	_, err = pipeline.SRem(l.ctx, cache.FollowedPrefix+uidStr, tid).Result()
	_, err = pipeline.SRem(l.ctx, cache.FollowerPrefix+tidStr, uid).Result()
	_, err = pipeline.HDel(l.ctx, cache.CountPrefix+uidStr).Result()
	_, err = pipeline.HDel(l.ctx, cache.CountPrefix+tidStr).Result()
	_, err = pipeline.Exec(l.ctx)
	if err != nil {
		logx.Errorw("[redis] transaction pipeline failed",
			logx.Field("err", err),
		)
		resp.Code = http.StatusInternalServerError
		resp.Msg = "[redis] del relation failed"
		return resp, err
	}
	logx.Info("[redis] del relation success")
	resp.Code = http.StatusOK
	resp.Msg = "success"

	return resp, nil

}
