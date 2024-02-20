package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckThumbupLogic {
	return &CheckThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckThumbupLogic) CheckThumbup(in *pb.CheckThumbupRequest) (*pb.CheckThumbupResponse, error) {
	exist, err := l.svcCtx.ThumbupDao.IsThumbup(l.ctx, in.UserId, in.VideoId)
	if err != nil {
		logx.Error("is favorite videoservice failed: " + err.Error())
		return nil, err
	}
	if exist {
		return &pb.CheckThumbupResponse{
			IsThumbup: 1,
		}, nil
	} else {
		return &pb.CheckThumbupResponse{
			IsThumbup: 0,
		}, nil

	}
}
