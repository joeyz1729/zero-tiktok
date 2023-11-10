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

// UpdateFavoriteInfo 点赞操作时，修改用户和视频作者的计数信息
func (l *UpdateFavoriteInfoLogic) UpdateFavoriteInfo(in *model.UpdateFavoriteInfoRequest) (*model.UpdateFavoriteInfoResponse, error) {
	var err error
	if in.ActionType {
		err = l.svcCtx.UserRepo.AddFollow(in.UserId, in.AuthorId)
	} else {
		err = l.svcCtx.UserRepo.DelFollow(in.UserId, in.AuthorId)
	}
	if err != nil {
		return nil, err
	}
	return &model.UpdateFavoriteInfoResponse{}, nil
}
