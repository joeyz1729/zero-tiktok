package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/joeyz1729/zero-tiktok/pkg/tool"
	"gorm.io/gorm"

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
	resp := new(pb.RegisterResponse)

	_, err := l.svcCtx.UserRepo.DBGetUserByName(in.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == nil {
		return resp, errors.New("user already exists")
	}
	uid, err := l.svcCtx.IDGenerator.Snowflake.NextID()
	if err != nil {
		return resp, errors.New("generate tiktok-user id failed")
	}
	err = l.svcCtx.UserRepo.DBCreateUser(int64(uid), in.GetUsername(), tool.Encrypt(in.GetPassword()))
	if err != nil {
		return resp, err
	}
	return &pb.RegisterResponse{UserId: int64(uid)}, nil
}
