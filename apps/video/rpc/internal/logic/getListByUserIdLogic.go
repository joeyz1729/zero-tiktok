package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListByUserIdLogic {
	return &GetListByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetListByUserIdLogic) GetListByUserId(in *model.GetListByUserIdRequest) (*model.GetListByUserIdResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetListByUserIdResponse{}, nil
}
