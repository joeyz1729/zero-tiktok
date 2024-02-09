package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/model"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/user"
	"github.com/zeromicro/go-zero/core/logx"

	_ "github.com/go-sql-driver/mysql"
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

// Action 用户关注操作的rpc，使用分布式事务
func (l *ActionLogic) Action(in *model.ActionRequest) (*model.ActionResponse, error) {
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
	userServer, err := l.svcCtx.Config.UserRpc.BuildTarget()
	if err != nil {
		return nil, err
	}
	gid := dtmgrpc.MustGenGid(l.svcCtx.DtmServer)
	msg := dtmgrpc.NewMsgGrpc(l.svcCtx.DtmServer, gid).Add(userServer+"/tiktok-user.User/UpdateFollowInfo", &user.UpdateFollowInfoRequest{
		UserId:     in.UserId,
		ToUserId:   in.ToUserId,
		ActionType: in.ActionType == int32(1),
	})

	err = msg.DoAndSubmitDB(userServer+"/tiktok-user.User/FollowPrepare", l.svcCtx.BarrierDB, func(tx *sql.Tx) error {
		var e error
		if in.ActionType == int32(1) {
			e = l.svcCtx.FollowRepo.AddRelation(in.UserId, in.ToUserId)
		} else {
			e = l.svcCtx.FollowRepo.DelRelation(in.UserId, in.ToUserId)
		}
		return e
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
