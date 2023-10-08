package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteLogic {
	return &GetFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteLogic) GetFavorite(in *model.GetFavoriteRequest) (*model.GetFavoriteResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetFavoriteResponse{}, nil
}
