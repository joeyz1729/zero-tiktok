package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/model"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFavoriteInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteInfoLogic {
	return &UpdateFavoriteInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFavoriteInfoLogic) UpdateFavoriteInfo(in *model.UpdateFavoriteInfoRequest) (*model.UpdateFavoriteInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &model.UpdateFavoriteInfoResponse{}, nil
}
