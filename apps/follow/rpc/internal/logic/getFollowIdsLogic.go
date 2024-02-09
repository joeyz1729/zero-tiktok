package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/model"
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

// GetFollowIds 获取关注列表
func (l *GetFollowIdsLogic) GetFollowIds(in *model.GetFollowIdsRequest) (*model.GetFollowIdsResponse, error) {
	resp := new(model.GetFollowIdsResponse)
	ids, err := l.svcCtx.FollowRepo.GetFollowedIds(in.UserId)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("empty set")
	}
	resp.FollowIds = ids
	return resp, nil
}
