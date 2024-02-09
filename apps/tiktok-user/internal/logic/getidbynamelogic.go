package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
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
	// 更新mysql计数表和redis count值，userId的关注数，toUserId的粉丝数

	return &pb.GetIdByNameResponse{}, nil
}
