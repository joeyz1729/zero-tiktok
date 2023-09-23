package logic

import (
	"context"
	"encoding/json"
	"github.com/YiZou89/zero-tiktok/apps/follow/dao"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/data/cache"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type ActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ActionLogic) Action(in *model.ActionRequest) (*model.ActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.ActionResponse)
	var err error
	// 1. 添加 mq 异步修改

	actionData, err := json.Marshal(dao.Action{
		UserId:     in.UserId,
		ToUserId:   in.ToUserId,
		ActionType: in.ActionType,
	})
	if err != nil {
		resp.Msg = "json marshal failed"
		logx.Errorw("json marshal failed",
			logx.Field("err", err))
		return resp, err
	}
	err = l.svcCtx.KqPusher.Push(string(actionData))
	if err == nil {
		logx.Info("push kafka mq success")
		resp.Msg = "push kafka mq success"
		return resp, nil
	}
	// 2. 如果 mq 异步失败，则同步修改数据库 + redis
	logx.Info("push kafka mq failed, start update mysql and redis")
	uidStr, tidStr := strconv.FormatInt(in.UserId, 10), strconv.FormatInt(in.UserId, 10)

	// 1. 添加到 bloom filter
	//bloomKey := data.BloomPrefix + uidStr + ":" + tidStr
	//err = l.svcCtx.Filter.AddCtx(l.ctx, []byte(bloomKey))
	//if err != nil {
	//	logx.Errorw("[bloom filter] add failed",
	//		logx.Field("err", err),
	//	)
	//	resp.Code = http.StatusInternalServerError
	//	resp.Msg = "add bloom filter failed"
	//	return resp, err
	//}
	//logx.Info("add bloom filter success")

	// 添加数据库两张表，并修改计数
	sqlStr := `insert into tiktok_follow.followed(user_id, followed_id) value(?, ?)`
	_, err = l.svcCtx.FollowDB.Exec(sqlStr, in.UserId, in.ToUserId)
	if err != nil {
		logx.Errorw("[mysql] add followed failed",
			logx.Field("err", err),
		)
	}
	sqlStr = `insert into tiktok_follow.follower(user_id, follower_id) value(?, ?)`
	_, err = l.svcCtx.FollowDB.Exec(sqlStr, in.ToUserId, in.UserId)
	if err != nil {
		logx.Errorw("[mysql] add follower failed",
			logx.Field("err", err))
		// rollback followed
	}

	// 删除redis数据
	pipeline := l.svcCtx.FollowCache.TxPipeline()
	pipeline.SRem(l.ctx, cache.FollowedPrefix+uidStr, tidStr)
	pipeline.SRem(l.ctx, cache.FollowerPrefix+tidStr, uidStr)
	pipeline.HDel(l.ctx, cache.CountPrefix+uidStr)
	pipeline.HDel(l.ctx, cache.CountPrefix+tidStr)
	_, err = pipeline.Exec(l.ctx)
	if err != nil {
		// 注意要添加缓存时间保证最终一致性
		logx.Errorw("delete redis failed",
			logx.Field("err", err))
		return resp, err
	}
	logx.Info("delete redis success")
	return resp, nil
}
