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
	// 根据author id查询发布的视频列表
	// 不需要查询author信息了，只需要查询info和favor

	return &model.GetListByAuthorIdResponse{}, nil
}
