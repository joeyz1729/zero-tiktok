package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoIdsLogic {
	return &GetVideoIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoIdsLogic) GetVideoIds(in *model.GetVideoIdsRequest) (*model.GetVideoIdsResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetVideoIdsResponse)
	//var videoIds []int64
	//
	ids, err := l.svcCtx.FavorRepo.GetFavorIds(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	resp.VideoNum = int64(len(ids))
	resp.VideoIds = ids
	return resp, nil
}
