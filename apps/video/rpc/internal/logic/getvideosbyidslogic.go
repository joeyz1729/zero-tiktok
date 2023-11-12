package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/data"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"

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
func (l *GetVideosByIdsLogic) GetVideosByIds(in *data.GetVideosByIdsRequest) (*data.GetVideosByIdsResponse, error) {
	// todo: add your logic here and delete this line

	return &data.GetVideosByIdsResponse{}, nil
}
