package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/model"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerIdsLogic {
	return &GetFollowerIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowerIdsLogic) GetFollowerIds(in *model.GetFollowerIdsRequest) (*model.GetFollowerIdsResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetFollowerIdsResponse)
	var err error
	var ids []int64
	// redis hit
	ids, err = l.svcCtx.FollowCache.GetFollowerIds(l.ctx, in.UserId)
	if err == nil {
		logx.Info("get follower ids from redis success")
		resp.FollowerIds = ids
		return resp, nil
	}
	// mysql
	ids, err = l.svcCtx.FollowDB.GetFollowerIds(l.ctx, in.UserId)
	if err != nil {
		logx.Error("get follower ids from mysql failed", err)
		return nil, err
	}

	// add ids to redis
	err = l.svcCtx.FollowCache.AddFollower(l.ctx, in.UserId, ids)
	if err != nil {
		return nil, err
	}

	resp.FollowerIds = ids
	return resp, nil
}
