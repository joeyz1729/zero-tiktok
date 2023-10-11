package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListByAuthorIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListByAuthorIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListByAuthorIdLogic {
	return &GetListByAuthorIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetListByAuthorIdLogic) GetListByAuthorId(in *model.GetListByAuthorIdRequest) (*model.GetListByAuthorIdResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetListByAuthorIdResponse{}, nil
}
