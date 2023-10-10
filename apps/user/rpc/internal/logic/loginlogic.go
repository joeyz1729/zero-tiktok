package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/model"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
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

func (l *LoginLogic) Login(in *model.LoginRequest) (*model.LoginResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.LoginResponse)
	userId, err := l.svcCtx.UserRepo.CheckLogin(in.Username, in.Password)
	if err != nil {
		return resp, err
	}

	resp.UserId = userId
	return resp, nil
}
