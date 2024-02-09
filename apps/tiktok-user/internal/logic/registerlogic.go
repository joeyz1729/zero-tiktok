package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/joeyz1729/zero-tiktok/pkg/tool"

	"github.com/joeyz1729/zero-tiktok/pkg/snowflake"
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

func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(pb.RegisterResponse)

	exist, err := l.svcCtx.UserRepo.CheckUserValid(in.Username)
	if err != nil {
		return resp, err
	}
	if exist {
		// username already exists.
		return resp, errors.New("username already exists")
	}
	// !exist, generate tiktok-user id
	uid, err := snowflake.GenID()
	if err != nil {
		return resp, errors.New("generate tiktok-user id failed")
	}
	err = l.svcCtx.UserRepo.Register(int64(uid), in.GetUsername(), tool.Encrypt(in.GetPassword()))
	if err != nil {
		return resp, err
	}
	return &pb.RegisterResponse{UserId: int64(uid)}, nil
}
