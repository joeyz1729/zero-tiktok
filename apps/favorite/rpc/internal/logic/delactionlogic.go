package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type DelActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelActionLogic {
	return &DelActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *DelActionLogic) DelAction(in *model.ActionRequest) (*model.ActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.ActionResponse)
	//var err error
	//logx.Info(in.UserId, in.VideoId)
	//
	//vidStr := strconv.Itoa(int(in.VideoId))
	//uidStr := strconv.Itoa(int(in.UserId))
	//exist, err := l.svcCtx.FavorRepo.FavorCache.SIsMember(l.ctx, data.FavoriteSetPrefix+vidStr, uidStr).Result()
	//if err != nil {
	//	logx.Errorw("redis check failed",
	//		logx.Field("err", err))
	//	return resp, err
	//}
	//if !exist {
	//	resp.Code = int32(0)
	//	resp.Msg = "repeat operation"
	//	return resp, nil
	//}
	//_, err = l.svcCtx.FavorRepo.FavorCache.SRem(l.ctx, data.FavoriteSetPrefix+vidStr, uidStr).Result()
	//if err != nil {
	//	resp.Code = int32(1)
	//	resp.Msg = "[redis] del msg failed"
	//	return resp, err
	//}
	//
	//l.svcCtx.FavorRepo.DelFavorite(uidStr, vidStr)

	resp.Code = http.StatusOK
	resp.Msg = "update record success"
	return resp, nil
}

//
//func (l *DelActionLogic) DelActionSync(in *data.ActionRequest) (*data.ActionResponse, error) {
//	// todo: add your logic here and delete this line
//	resp := new(data.ActionResponse)
//	var err error
//	logx.Info(in.UserId, in.VideoId)
//
//	//vidStr := strconv.Itoa(int(in.VideoId))
//	//uidStr := strconv.Itoa(int(in.UserId))
//	//exist, err := l.svcCtx.FavorRepo.FavorCache.SIsMember(l.ctx, data.FavoriteSetPrefix+vidStr, "1"+uidStr).Result()
//	//if err != nil {
//	//	logx.Errorw("redis check failed",
//	//		logx.Field("err", err))
//	//	return resp, err
//	//}
//
//	var ifCancel bool
//	// 1. query record
//	sqlStr := fmt.Sprintf(`select cancel from tiktok_favorite.favorite where user_id = ? and video_id = ? limit 1`)
//	err = l.svcCtx.FavorRepo.FavorDB.Get(&ifCancel, sqlStr, in.UserId, in.VideoId)
//	if err != nil && err != sqlc.ErrNotFound {
//		// query failed
//		logx.Errorw("find relation failed",
//			logx.Field("err", err),
//		)
//		resp.Code = int32(0)
//		resp.Msg = "mysql query failed"
//		return resp, err
//
//	}
//	if err == sqlc.ErrNotFound {
//		// relation does not exist
//		resp.Code = int32(02)
//		resp.Msg = "unfollowing failed, relation does not exist"
//		return resp, nil
//	}
//	// err == nil, query success
//	sqlStr = `delete from tiktok_favorite.favorite where user_id = ? and video_id = ? limit 1`
//	_, err = l.svcCtx.FavorRepo.FavorDB.Exec(sqlStr, in.UserId, in.VideoId)
//	if err != nil {
//		logx.Errorw("update relation failed",
//			logx.Field("err", err),
//		)
//		resp.Code = int32(0)
//		resp.Msg = "update relation failed"
//		return resp, err
//	}
//
//	resp.Code = http.StatusOK
//	resp.Msg = "update record success"
//	return resp, nil
//}
