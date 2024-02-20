package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteThumbupLogic {
	return &DeleteThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消点赞
func (l *DeleteThumbupLogic) DeleteThumbup(in *pb.DeleteThumbupRequest) (*pb.DeleteThumbupResponse, error) {
	exist, err := l.svcCtx.ThumbupDao.IsThumbup(l.ctx, in.UserId, in.VideoId)
	if err != nil {
		logx.Error("is favorite videoservice failed: " + err.Error())
		return nil, err
	}
	if !exist {
		return nil, ErrRepeatedOperation
	}

	err = l.svcCtx.ThumbupDao.DeleteThumbup(l.ctx, in.UserId, in.VideoId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteThumbupResponse{}, nil
}
