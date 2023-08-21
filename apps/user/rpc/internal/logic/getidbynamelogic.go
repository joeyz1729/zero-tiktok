package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIdByNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIdByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIdByNameLogic {
	return &GetIdByNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIdByNameLogic) GetIdByName(in *model.GetIdByNameRequest) (*model.GetIdByNameResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetIdByNameResponse{}, nil
}
