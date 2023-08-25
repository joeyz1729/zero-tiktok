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
	// todo: add your logic here and delete this line

	return &model.GetFollowIdsResponse{}, nil
}
