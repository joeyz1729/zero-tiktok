package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkCountLogic {
	return &GetWorkCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkCountLogic) GetWorkCount(in *model.GetWorkCountRequest) (*model.GetWorkCountResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetWorkCountResponse{}, nil
}
