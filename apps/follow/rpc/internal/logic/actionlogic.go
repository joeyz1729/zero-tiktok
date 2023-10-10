package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/YiZou89/zero-tiktok/apps/follow/dao"
	"github.com/YiZou89/zero-tiktok/apps/follow/model"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	IfPushMq = false
)

var (
	ErrRepeatedOperation = errors.New("repeated operation")
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
	logx.Infof("user_id: %d, to_user_id: %d, action_type: %d", in.UserId, in.ToUserId, in.ActionType)
	var resp = new(model.ActionResponse)
	// 0. 先检查关系是否重复
	ok, err := l.svcCtx.FollowCache.GetRelation(l.ctx, in.UserId, in.ToUserId)
	if err == nil && ok {
		// 查询成功，存在关系
		if in.ActionType == int32(1) {
			return nil, ErrRepeatedOperation
		}
	}
	// redis没查到或出错，走数据库
	ok, err = l.svcCtx.FollowDB.CheckRelation(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if ok && in.ActionType == int32(1) || !ok && in.ActionType == int32(2) {
		return nil, ErrRepeatedOperation
	}

	// 1. 添加 mq 异步修改
	if IfPushMq {
		actionData, err := json.Marshal(dao.Action{
			UserId:     in.UserId,
			ToUserId:   in.ToUserId,
			ActionType: in.ActionType,
		})
		if err != nil {
			logx.Errorw("json marshal failed",
				logx.Field("err", err))
			return nil, err
		}
		err = l.svcCtx.KqPusher.Push(string(actionData))
		//err = l.svcCtx.KqWriter.WriteMessages(l.ctx,
		//	kafka.Message{
		//		Value: actionData,
		//	})
		if err == nil {
			logx.Info("push kafka mq success")
			resp.Msg = "push kafka mq success"
			_, err = l.svcCtx.UserRpc.UpdateFollowInfo(l.ctx, &user.UpdateFollowInfoRequest{
				UserId:     in.UserId,
				ToUserId:   in.ToUserId,
				ActionType: in.ActionType == int32(1),
			})
			if err != nil {
				return resp, err
			}
			return resp, nil
		}
		logx.Info("push kafka mq failed, start update mysql and redis")
		// 2. 如果 mq 异步失败，则同步修改数据库 + redis
	}

	// 添加数据库两张表，并修改计数
	if in.ActionType == int32(1) {
		sqlStr := `insert into tiktok_follow.followed(user_id, followed_id) value(?, ?)`
		_, err = l.svcCtx.FollowDB.Exec(sqlStr, in.UserId, in.ToUserId)
		if err != nil {
			logx.Errorw("[mysql] add followed failed",
				logx.Field("err", err),
			)
		}
	} else {
		sqlStr := `delete from tiktok_follow.followed where user_id = ? and followed_id = ? limit 1`
		_, err = l.svcCtx.FollowDB.Exec(sqlStr, in.UserId, in.ToUserId)
		if err != nil {
			logx.Errorw("[mysql] del followed failed",
				logx.Field("err", err),
			)
		}
	}

	// 更新user模块计数
	_, err = l.svcCtx.UserRpc.UpdateFollowInfo(l.ctx, &user.UpdateFollowInfoRequest{
		UserId:     in.UserId,
		ToUserId:   in.ToUserId,
		ActionType: in.ActionType == int32(1),
	})
	if err != nil {
		return resp, err
	}
	// 删除redis数据
	if err = l.svcCtx.FollowCache.RemRelation(l.ctx, in.UserId, in.ToUserId); err != nil {
		// TODO, 异步添加消息队列进行删除
		logx.Errorw("delete redis failed",
			logx.Field("err", err))
		return resp, err
	}

	logx.Info("delete redis success")
	return resp, nil
}

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
