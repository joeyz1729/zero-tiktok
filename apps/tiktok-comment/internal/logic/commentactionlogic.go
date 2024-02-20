package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentActionLogic) CommentAction(in *pb.CommentActionRequest) (*pb.CommentActionResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.CommentActionResponse{}, nil
}
