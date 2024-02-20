package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoCommentListLogic {
	return &GetVideoCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoCommentListLogic) GetVideoCommentList(in *pb.GetVideoCommentListRequest) (*pb.GetVideoCommentListResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetVideoCommentListResponse{}, nil
}
