package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/model"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"

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

func (l *GetAuthorLogic) GetAuthor(in *model.GetAuthorRequest) (*model.GetAuthorResponse, error) {
	// todo: add your logic here and delete this line

	return &model.GetAuthorResponse{}, nil
}
