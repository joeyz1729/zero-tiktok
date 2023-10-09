package user

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserRegisterRequest) (resp *types.UserRegisterResponse, err error) {
	logx.Debugf("register req: %v\n", req)
	resp = new(types.UserRegisterResponse)
	res, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{Username: req.Username, Password: req.Password})
	if err != nil {
		logx.Errorw("user.Register failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}

	aToken, _, err := jwtx.GenToken(res.GetUserId(), req.Username)
	if err != nil {
		logx.Errorw("jwt.GenToken failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}

	resp.UserId = res.GetUserId()
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	resp.Token = aToken

	return resp, nil
}
