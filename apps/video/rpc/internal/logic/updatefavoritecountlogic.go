package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFavoriteCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteCountLogic {
	return &UpdateFavoriteCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFavoriteCountLogic) UpdateFavoriteCount(in *model.UpdateFavoriteCountRequest) (*model.UpdateFavoriteCountResponse, error) {
	var err error
	if in.ActionType {
		err = l.svcCtx.VideoRepo.AddFavoriteCount(l.ctx, in.VideoId)
	} else {
		err = l.svcCtx.VideoRepo.DelFavoriteCount(l.ctx, in.VideoId)
	}
	if err != nil {
		return nil, err
	}
	return &model.UpdateFavoriteCountResponse{}, nil
}
