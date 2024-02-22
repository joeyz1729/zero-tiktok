package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp := new(pb.LoginResponse)
	user, err := l.svcCtx.UserRepo.DBGetUserByName(in.Username)
	if err != nil {
		return resp, err
	}
	if user.Password != in.Password {
		return nil, ErrInvalidPassword
	}
	resp.UserId = user.ID
	return resp, nil
}
