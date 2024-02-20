package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetThumbupVideoIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetThumbupVideoIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetThumbupVideoIdsLogic {
	return &GetThumbupVideoIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetThumbupVideoIdsLogic) GetThumbupVideoIds(in *pb.GetThumbupVideoIdsRequest) (*pb.GetThumbupVideoIdsResponse, error) {
	resp := new(pb.GetThumbupVideoIdsResponse)
	ids, err := l.svcCtx.ThumbupDao.GerUserThumbupList(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	resp.VideoIds = ids
	return resp, nil
}
