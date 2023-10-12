package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/model"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuthorsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAuthorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorsLogic {
	return &GetAuthorsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAuthorsLogic) GetAuthors(in *model.GetAuthorsRequest) (*model.GetAuthorsResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetAuthorsResponse{}, nil
}
