package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideosByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideosByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideosByIdsLogic {
	return &GetVideosByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据视频id查询视频详细信息，不包括作者详细信息
func (l *GetVideosByIdsLogic) GetVideosByIds(in *model.GetVideosByIdsRequest) (*model.GetVideosByIdsResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetVideosByIdsResponse{}, nil
}
