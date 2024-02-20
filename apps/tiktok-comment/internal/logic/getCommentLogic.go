package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentLogic {
	return &GetCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentLogic) GetComment(in *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetCommentResponse{}, nil
}
