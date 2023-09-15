package logic

import (
	"context"
	"errors"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *model.UserInfoRequest) (*model.UserInfoResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.UserInfoResponse)
	user, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return resp, errors.New("user id does not exist")
		}
		return resp, errors.New("mysql query error")
	}

	resp.User.Name = user.Username
	resp.User.Id = user.UserId
	resp.User.Avatar = "no avatar"
	resp.User.BackgroundImage = "no background image"
	resp.User.Signature = "no signature"

	// get follow count

	// get favorite count

	// get work count

	return resp, nil

	return &model.UserInfoResponse{}, nil
}
