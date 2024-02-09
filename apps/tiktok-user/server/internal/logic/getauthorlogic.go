package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuthorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAuthorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorLogic {
	return &GetAuthorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAuthorLogic) GetAuthor(in *pb.GetAuthorRequest) (*pb.GetAuthorResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetAuthorResponse{}, nil
}
