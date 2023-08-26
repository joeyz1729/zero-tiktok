package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/message/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/message/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *model.ListRequest) (*model.ListResponse, error) {
	// todo: add your logic here and delete this line

	return &model.ListResponse{}, nil
}
