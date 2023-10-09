package logic

import (
	"context"
	"errors"
	model2 "github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/model"
	"github.com/YiZou89/zero-tiktok/pkg/tool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/pkg/snowflake"
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

func (l *RegisterLogic) Register(in *model2.RegisterRequest) (*model2.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model2.RegisterResponse)
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)

	if err == nil {
		// username already exists.
		resp.UserId = user.UserId
		return resp, errors.New("username already exists")
	}

	if err != nil && err != sqlc.ErrNotFound {
		// query err
		return resp, err
	}

	// generate user id
	uid, err := snowflake.GenID()
	if err != nil {
		return resp, errors.New("generate user id failed")
	}
	// encrypt password
	// insert into mysql
	//_, err = l.svcCtx.UserModel.Insert(l.ctx, &model2.User{
	//	UserId:   int64(uid),
	//	Username: in.GetUsername(),
	//	Password: tool.Encrypt(in.GetPassword()),
	//})
	err = l.svcCtx.UserRepo.Register(int64(uid), in.GetUsername(), tool.Encrypt(in.GetPassword()))

	if err != nil {
		return resp, err
	}

	return &model2.RegisterResponse{UserId: int64(uid)}, nil
}
