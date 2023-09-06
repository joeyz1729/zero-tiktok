package logic

import (
	"context"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/zeromicro/go-zero/core/logx"
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

// Deprecated: Use Add or Del instead.
func (l *ActionLogic) Action(in *model.ActionRequest) (*model.ActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.ActionResponse)
	var err error

	var ifCancel bool
	sqlStr := fmt.Sprintf(`select cancel from tiktok_follow.follow where user_id = ? and follower_id = ? limit 1`)
	err = l.svcCtx.FollowDB.Get(&ifCancel, sqlStr, in.UserId, in.ToUserId)
	if err != nil {
		fmt.Println(err)
		if err != sqlc.ErrNotFound {
			// query failed
			logx.Errorw("find relation failed",
				logx.Field("err", err),
			)
			resp.Msg = "mysql query failed"
			return resp, err
		}
		// relation does not exist
		if in.ActionType == int32(2) {
			resp.Msg = "unfollowing failed, relation does not exist"
			return resp, nil
		}
		// add follow relation
		_, err = l.svcCtx.FollowModel.Insert(l.ctx, &model.Follow{
			UserId:     in.UserId,
			FollowerId: in.ToUserId,
			Cancel:     false,
		})
		if err != nil {
			logx.Errorw("add following relation failed",
				logx.Field("err", err),
			)
			resp.Msg = "internal server error"
			return resp, err
		}
		resp.Msg = "add follow relation success"
		return resp, nil
	}
	fmt.Println("success")
	// query success
	if in.ActionType == int32(1) && ifCancel == false || in.ActionType == int32(2) && ifCancel == true {
		resp.Msg = "repeat operation"
		return resp, nil
	}
	var cancel bool
	if in.ActionType == int32(2) {
		cancel = true
	}
	sqlStr = `update tiktok_follow.follow set cancel = ? where user_id = ? and follower_id = ?`
	_, err = l.svcCtx.FollowDB.Exec(sqlStr, cancel, in.UserId, in.ToUserId)
	if err != nil {
		logx.Errorw("update relation failed",
			logx.Field("err", err),
		)
		resp.Msg = "update relation failed"
		return resp, err
	}

	resp.Msg = "update relation success"
	return resp, nil
}
