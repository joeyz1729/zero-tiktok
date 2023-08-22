package user

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	// todo: add your logic here and delete this line
	// token 校验
	resp = new(types.UserInfoResponse)

	res, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{UserId: req.UserId})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		logx.Error(err)
		return resp, nil
	}
	resp.UserInfo = types.User{
		Id:              req.UserId,
		Name:            res.Name,
		Avatar:          res.Avatar,
		BackgroundImage: res.BackgroundImage,
		Signature:       res.Signature,
	}

	// 调用follow模块

	// 调用video模块和like模块

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
