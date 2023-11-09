package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFavoriteInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteInfoLogic {
	return &UpdateFavoriteInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFavoriteInfoLogic) UpdateFavoriteInfo(in *model.UpdateFavoriteInfoRequest) (*model.UpdateFavoriteInfoResponse, error) {
	// todo: add your logic here and delete this line
	// 需要更新的内容：uid的点赞数量，aid的被点赞数量，
	//err := l.svcCtx.UserRepo.UpdateFavoriteCount(in.UserId, in.ActionType)
	//if err != nil {
	//	return nil, err
	//}
	//err = l.svcCtx.UserRepo.UpdateTotalFavorited(in.AuthorId, in.ActionType)
	//if err != nil {
	//	return nil, err
	//}
	err := l.svcCtx.UserRepo.UpdateFavorTx(in.UserId, in.VideoId, in.ActionType)
	if err != nil {
		logx.Error("update favor tx ", err)
		return nil, err
	}
	return &model.UpdateFavoriteInfoResponse{}, nil
}
