package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/joeyz1729/zero-tiktok/pkg/tool"
	"github.com/zeromicro/go-zero/core/logx"
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
	userId, err := l.svcCtx.UserRepo.Login(in.Username, tool.Encrypt(in.Password))
	if err != nil {
		return resp, err
	}

	resp.UserId = userId
	return resp, nil
}
