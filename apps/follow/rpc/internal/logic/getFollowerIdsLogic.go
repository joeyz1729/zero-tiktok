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

func (l *GetFollowerIdsLogic) GetFollowerIds(in *model.GetFollowerIdsRequest) (*model.GetFollowerIdsResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetFollowerIdsResponse{}, nil
}
