package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"

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

// GetFollowerIds 获取粉丝列表
func (l *GetFollowerIdsLogic) GetFollowerIds(in *model.GetFollowerIdsRequest) (*model.GetFollowerIdsResponse, error) {
	resp := new(model.GetFollowerIdsResponse)
	ids, err := l.svcCtx.FollowRepo.GetFollowerIds(in.UserId)
	if err != nil {
		return nil, err
	}
	resp.FollowerIds = ids
	return resp, nil
}
