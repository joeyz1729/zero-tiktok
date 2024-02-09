package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuthorsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAuthorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorsLogic {
	return &GetAuthorsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAuthorsLogic) GetAuthors(in *pb.GetAuthorsRequest) (*pb.GetAuthorsResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetAuthorsResponse{}, nil
}
