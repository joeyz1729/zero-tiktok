package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishActionLogic) PublishAction(in *video.PublishActionRequest) (*video.PublishActionResponse, error) {
	// todo: add your logic here and delete this line

	return &video.PublishActionResponse{}, nil
}
