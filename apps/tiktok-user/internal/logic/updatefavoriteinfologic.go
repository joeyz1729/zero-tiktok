package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
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
func (l *UpdateFavoriteInfoLogic) UpdateFavoriteInfo(in *pb.UpdateFavoriteInfoRequest) (*pb.UpdateFavoriteInfoResponse, error) {
	var err error
	if in.ActionType {
		err = l.svcCtx.UserRepo.AddFavoriteRelation(in.UserId, in.AuthorId)
	} else {
		err = l.svcCtx.UserRepo.DelFavoriteRelation(in.UserId, in.AuthorId)
	}
	if err != nil {
		return nil, err
	}
	return &pb.UpdateFavoriteInfoResponse{}, nil
}
