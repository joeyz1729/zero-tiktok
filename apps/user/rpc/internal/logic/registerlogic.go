package logic

import (
	"context"
	"errors"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/model"
	"github.com/YiZou89/zero-tiktok/pkg/tool"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *model.RegisterRequest) (*model.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.RegisterResponse)

	exist, err := l.svcCtx.UserRepo.CheckUserValid(in.Username)
	if err != nil {
		return resp, err
	}
	if exist {
		// username already exists.
		return resp, errors.New("username already exists")
	}
	// !exist, generate user id
	uid, err := snowflake.GenID()
	if err != nil {
		return resp, errors.New("generate user id failed")
	}
	err = l.svcCtx.UserRepo.Register(int64(uid), in.GetUsername(), tool.Encrypt(in.GetPassword()))
	if err != nil {
		return resp, err
	}
	return &model.RegisterResponse{UserId: int64(uid)}, nil
}
