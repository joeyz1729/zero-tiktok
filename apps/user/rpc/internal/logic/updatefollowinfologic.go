package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/model"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFollowInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFollowInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFollowInfoLogic {
	return &UpdateFollowInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFollowInfoLogic) UpdateFollowInfo(in *model.UpdateFollowInfoRequest) (*model.UpdateFollowInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &model.UpdateFollowInfoResponse{}, nil
}
