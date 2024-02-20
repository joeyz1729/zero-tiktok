package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrRepeatedOperation = errors.New("repeated operation")
)

type AddThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddThumbupLogic {
	return &AddThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddThumbupLogic) AddThumbup(in *pb.AddThumbupRequest) (*pb.AddThumbupResponse, error) {
	exist, err := l.svcCtx.ThumbupDao.IsThumbup(l.ctx, in.UserId, in.VideoId)
	if err != nil {
		logx.Error("is favorite videoservice failed: " + err.Error())
		return nil, err
	}
	if exist {
		return nil, ErrRepeatedOperation
	}

	err = l.svcCtx.ThumbupDao.AddThumbup(l.ctx, in.UserId, in.VideoId)
	if err != nil {
		return nil, err
	}
	return &pb.AddThumbupResponse{}, nil
}
