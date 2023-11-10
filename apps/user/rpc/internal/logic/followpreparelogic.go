package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/model"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowPrepareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowPrepareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowPrepareLogic {
	return &FollowPrepareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowPrepareLogic) FollowPrepare(in *model.FollowPrepareRequest) (*model.FollowPrepareRequest, error) {
	// todo: add your logic here and delete this line

	return &model.FollowPrepareRequest{}, nil
}
