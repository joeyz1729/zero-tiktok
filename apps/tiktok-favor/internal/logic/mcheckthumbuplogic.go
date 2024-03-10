package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MCheckThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMCheckThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MCheckThumbupLogic {
	return &MCheckThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MCheckThumbupLogic) MCheckThumbup(in *pb.MCheckThumbupRequest) (*pb.MCheckThumbupResponse, error) {
	var resp = new(pb.MCheckThumbupResponse)
	resp.IsThumbup = make([]bool, len(in.VideoIds))
	if in.UserId == 0 {
		return resp, nil
	}
	for i, vid := range in.VideoIds {
		exist, err := l.svcCtx.ThumbupDao.IsThumbup(l.ctx, in.UserId, vid)
		if err != nil {
			logx.Error("is favorite videoservice failed: " + err.Error())
			return nil, err
		}
		resp.IsThumbup[i] = exist
	}
	return resp, nil

}
