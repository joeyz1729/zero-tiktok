package user

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MrInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMrInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MrInfoLogic {
	return &MrInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MrInfoLogic) MrInfo(req *types.UserMrInfoRequest) (resp *types.UserMrInfoResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
