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

// GetFavorite 检查用户和视频的点赞关系
func (l *GetFavoriteLogic) GetFavorite(in *model.GetFavoriteRequest) (*model.GetFavoriteResponse, error) {
	exist, err := l.svcCtx.FavorRepo.CheckFavor(l.ctx, in.UserId, in.VideoId)
	if err != nil {
		return nil, err
	}
	resp := new(model.GetFavoriteResponse)
	if exist {
		resp.ActionType = int64(1)
	} else {
		resp.ActionType = int64(2)
	}
	return resp, nil
}
