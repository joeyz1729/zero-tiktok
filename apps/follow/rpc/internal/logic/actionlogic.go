package logic

import (
	"context"
	"errors"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	IfPushMq = false
)

var (
	ErrRepeatedOperation = errors.New("repeated operation")
)

type ActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ActionLogic) Action(in *model.ActionRequest) (*model.ActionResponse, error) {
	// todo: add your logic here and delete this line
	logx.Infof("user_id: %d, to_user_id: %d, action_type: %d", in.UserId, in.ToUserId, in.ActionType)
	var resp = new(model.ActionResponse)
	// 检查关系
	ok, err := l.svcCtx.FollowRepo.CheckRelation(in.UserId, in.ToUserId)
	if err != nil {
		logx.Error("check relation failed", err)
		return nil, err
	}
	if ok && in.ActionType == int32(1) || !ok && in.ActionType == int32(2) {
		return nil, ErrRepeatedOperation
	}

	// 修改关系
	if in.ActionType == int32(1) {
		err = l.svcCtx.FollowRepo.AddRelation(in.UserId, in.ToUserId)
	} else {
		err = l.svcCtx.FollowRepo.DelRelation(in.UserId, in.ToUserId)
	}
	if err != nil {
		return nil, err
	}

	// 更新计数
	_, err = l.svcCtx.UserRpc.UpdateFollowInfo(l.ctx, &user.UpdateFollowInfoRequest{
		UserId:     in.UserId,
		ToUserId:   in.ToUserId,
		ActionType: in.ActionType == int32(1),
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
