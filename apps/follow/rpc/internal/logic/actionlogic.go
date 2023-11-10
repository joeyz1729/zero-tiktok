package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"

	_ "github.com/go-sql-driver/mysql"
)

const (
	IfPushMq = false
)

var (
	ErrRepeatedOperation = errors.New("repeated operation")
	DtmServer            = "etcd://localhost:2379/dtmservice"
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
	userServer, err := l.svcCtx.Config.UserRpc.BuildTarget()
	if err != nil {
		return nil, err
	}
	gid := dtmgrpc.MustGenGid(DtmServer)
	msg := dtmgrpc.NewMsgGrpc(DtmServer, gid).Add(userServer+"/user.User/UpdateFollowInfo", &user.UpdateFollowInfoRequest{
		UserId:     in.UserId,
		ToUserId:   in.ToUserId,
		ActionType: in.ActionType == int32(1),
	})

	db, err := sql.Open("mysql", "root:root1234@tcp(localhost:13306)/tiktok_follow?parseTime=true&charset=utf8")
	if err != nil {
		return nil, err
	}
	err = msg.DoAndSubmitDB(userServer+"/user.User/FollowPrepare", db, func(tx *sql.Tx) error {
		//userId, toUserId := in.UserId, in.ToUserId
		var e error
		if in.ActionType == int32(1) {
			e = l.svcCtx.FollowRepo.AddRelation(in.UserId, in.ToUserId)
			//sqlStr := `insert into tiktok_follow.followed(user_id, followed_id) value(?, ?)`
			//_, err = tx.Exec(sqlStr, userId, toUserId)
			//// 删除redis数据
			//err = l.svcCtx.FollowCache.RemRelation(l.ctx, userId, toUserId)
			//return err
		} else {
			e = l.svcCtx.FollowRepo.DelRelation(in.UserId, in.ToUserId)
			//sqlStr := `delete from tiktok_follow.followed where user_id = ? and followed_id = ? limit 1`
			//_, err = tx.Exec(sqlStr, userId, toUserId)
			//err = l.svcCtx.FollowCache.RemRelation(l.ctx, userId, toUserId)
			//return err
		}
		return e
	})
	// 修改关系
	//if in.ActionType == int32(1) {
	//	err = l.svcCtx.FollowRepo.AddRelation(in.UserId, in.ToUserId)
	//} else {
	//	err = l.svcCtx.FollowRepo.DelRelation(in.UserId, in.ToUserId)
	//}
	//if err != nil {
	//	return nil, err
	//}

	// 更新计数
	//_, err = l.svcCtx.UserRpc.UpdateFollowInfo(l.ctx, &user.UpdateFollowInfoRequest{
	//	UserId:     in.UserId,
	//	ToUserId:   in.ToUserId,
	//	ActionType: in.ActionType == int32(1),
	//})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
