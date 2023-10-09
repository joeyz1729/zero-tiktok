package logic

import (
	"context"
	"errors"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/model"
	"github.com/YiZou89/zero-tiktok/pkg/tool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
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
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return resp, errors.New("user not found")
		}
		return resp, errors.New("query failed")
	}

	if tool.Encrypt(in.Password) != user.Password {
		return resp, errors.New("incorrect password")
	}
	resp.UserId = user.UserId

	return resp, nil
}
