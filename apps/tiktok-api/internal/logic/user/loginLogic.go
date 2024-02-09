package user

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {

	resp = new(types.UserLoginResponse)
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &userservice.LoginRequest{Username: req.Username, Password: req.Password})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal service error"
		return resp, nil
	}
	aToken, _, err := jwtx.GenToken(res.GetUserId(), req.Username)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal service error"
		return resp, nil
	}

	resp.UserId = res.GetUserId()
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	resp.Token = aToken

	return resp, nil
}
