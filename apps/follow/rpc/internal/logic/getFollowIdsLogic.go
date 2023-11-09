package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
)
type GetFollowIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowIdsLogic {
	return &GetFollowIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowIdsLogic) GetFollowIds(in *model.GetFollowIdsRequest) (*model.GetFollowIdsResponse, error) {
	resp := new(model.GetFollowIdsResponse)
	var err error
	var ids []int64
	// redis hit
	ids, err = l.svcCtx.FollowCache.GetFollowedIds(l.ctx, in.UserId)
	if err == nil {
		logx.Info("get follower ids from redis success")
		resp.FollowIds = ids
		return resp, nil
	}
	// mysql
	ids, err = l.svcCtx.FollowDB.GetFollowedIds(l.ctx, in.UserId)
	if err != nil {
		logx.Error("get follower ids from mysql failed", err)
		return nil, err
	}

	// add ids to redis
	err = l.svcCtx.FollowCache.AddFollowed(l.ctx, in.UserId, ids)
	if err != nil {
		return nil, err
	}

	resp.FollowIds = ids
	return resp, nil
}
