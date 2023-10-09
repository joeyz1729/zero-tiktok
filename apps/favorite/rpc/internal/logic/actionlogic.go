package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"

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

func (l *ActionLogic) Action(in *model.ActionRequest) (*model.ActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.ActionResponse)
	logx.Info(in.UserId, in.VideoId, in.ActionType)
	var ifCancel bool

	// 1. query record
	sqlStr := fmt.Sprintf(`select cancel from tiktok_favorite.favorite where user_id = ? and video_id = ? limit 1`)
	err := l.svcCtx.FavorRepo.FavoriteDB.Get(&ifCancel, sqlStr, in.UserId, in.VideoId)
	if err != nil {
		if err != sqlc.ErrNotFound {
			logx.Errorw("find relation failed",
				logx.Field("err", err))
		}
		resp.Code = int32(0)
		resp.Msg = "mysql query failed"
		return resp, err
	}
	if err == sqlc.ErrNotFound {

	}
	if err != nil {
		fmt.Println(err)
		if err != sqlc.ErrNotFound {
			// query failed
			logx.Errorw("find relation failed",
				logx.Field("err", err),
			)
			resp.Code = int32(0)
			resp.Msg = "mysql query failed"
			return resp, err
		}
		// relation does not exist
		if in.ActionType == int32(2) {
			resp.Code = int32(02)
			resp.Msg = "unfollowing failed, relation does not exist"
			return resp, nil
		}
		// add follow relation
		_, err = l.svcCtx.FavorRepo.FavoriteModel.Insert(l.ctx, &model.Favorite{
			UserId:  in.UserId,
			VideoId: in.VideoId,
			Cancel:  false,
		})
		if err != nil {
			logx.Errorw("add following relation failed",
				logx.Field("err", err),
			)
			resp.Code = int32(0)
			resp.Msg = "mysql insert failed"
			return resp, err
		}
		resp.Code = int32(01)
		resp.Msg = "add follow relation success"
		return resp, nil
	}

	// query success, update
	if in.ActionType == int32(1) && ifCancel == false || in.ActionType == int32(2) && ifCancel == true {
		resp.Code = int32(00)
		resp.Msg = "repeat operation"
		return resp, nil
	}
	var cancel bool
	if in.ActionType == int32(2) {
		cancel = true
	}
	sqlStr = `update tiktok_favorite.favorite set cancel = ? where user_id = ? and video_id = ?`
	_, err = l.svcCtx.FavorRepo.FavoriteDB.Exec(sqlStr, cancel, in.UserId, in.VideoId)
	if err != nil {
		logx.Errorw("update relation failed",
			logx.Field("err", err),
		)
		resp.Code = int32(0)
		resp.Msg = "update relation failed"
		return resp, err
	}

	resp.Code = http.StatusOK
	resp.Msg = "update record success"
	return resp, nil
}
