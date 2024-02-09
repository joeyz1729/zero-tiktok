package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/server/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIdByNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIdByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIdByNameLogic {
	return &GetIdByNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIdByNameLogic) GetIdByName(in *pb.GetIdByNameRequest) (*pb.GetIdByNameResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetIdByNameResponse{}, nil
}
