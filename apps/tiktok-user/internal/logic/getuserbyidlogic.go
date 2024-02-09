package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserRepo.GetUserInfo(in.UserId)
	if err != nil {
		return nil, err
	}
	resp := new(pb.GetUserByIdResponse)
	resp.Id = in.UserId
	resp.Name = user.Username
	resp.Avatar = "no avatar"
	resp.BackgroundImage = "no background image"
	resp.Signature = "no signature"
	return resp, nil
}
