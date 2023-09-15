package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountLogic {
	return &GetFavoriteCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteCountLogic) GetFavoriteCount(in *model.GetFavoriteCountRequest) (*model.GetFavoriteCountResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetFavoriteCountResponse{}, nil
}
