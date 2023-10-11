package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListByAuthorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListByAuthorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListByAuthorLogic {
	return &GetListByAuthorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetListByAuthorLogic) GetListByAuthor(in *model.GetListByAuthorIdRequest) (*model.GetListByAuthorIdResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetListByAuthorIdResponse{}, nil
}
