package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWorkInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkInfoLogic {
	return &UpdateWorkInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateWorkInfoLogic) UpdateWorkInfo(in *model.UpdateWorkInfoRequest) (*model.UpdateWorkInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &model.UpdateWorkInfoResponse{}, nil
}
