package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddActionLogic {
	return &AddActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddActionLogic) AddAction(in *model.ActionRequest) (*model.ActionResponse, error) {
	//// todo: add your logic here and delete this line
	resp := new(model.ActionResponse)
	//logx.Info(in.UserId, in.VideoId)
	//
	//// 1. query record
	//vidStr := strconv.Itoa(int(in.VideoId))
	//uidStr := strconv.Itoa(int(in.UserId))
	//exist, err := l.svcCtx.FavorRepo.FavorCache.SIsMember(l.ctx, data.FavoriteSetPrefix+vidStr, uidStr).Result()
	//if err != nil {
	//	logx.Errorw("redis check failed",
	//		logx.Field("err", err))
	//	return resp, err
	//}
	//if exist {
	//	resp.Code = int32(0)
	//	resp.Msg = "repeat operation"
	//	return resp, nil
	//}
	//// 不存在，添加
	//_, err = l.svcCtx.FavorRepo.FavorCache.SAdd(l.ctx, data.FavoriteSetPrefix+vidStr, uidStr).Result()
	//if err != nil {
	//	resp.Code = int32(0)
	//	resp.Msg = "[redis] add favorite failed"
	//	return resp, err
	//}
	//l.svcCtx.FavorRepo.AddFavorite(uidStr, vidStr)

	resp.Code = int32(01)
	resp.Msg = "add follow relation success"
	return resp, nil
}

func (l *AddActionLogic) SyncAddAction(in *model.ActionRequest) (*model.ActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.ActionResponse)
	//logx.Info(in.UserId, in.VideoId)
	//
	//// 1. query record
	//vidStr := strconv.Itoa(int(in.VideoId))
	//uidStr := strconv.Itoa(int(in.UserId))
	//exist, err := l.svcCtx.FavorRepo.FavorCache.SIsMember(l.ctx, data.FavoriteSetPrefix+vidStr, "0"+uidStr).Result()
	//if err != nil {
	//	logx.Errorw("redis check failed",
	//		logx.Field("err", err))
	//	return resp, err
	//}
	//if exist {
	//	resp.Code = int32(0)
	//	resp.Msg = "repeat operation"
	//	return resp, nil
	//}
	//
	//var ifCancel bool
	//sqlStr := fmt.Sprintf(`select cancel from tiktok_favorite.favorite where user_id = ? and video_id = ? limit 1`)
	//err = l.svcCtx.FavorRepo.FavorDB.Get(&ifCancel, sqlStr, in.UserId, in.VideoId)
	//if err != nil && err != sqlc.ErrNotFound {
	//	logx.Errorw("find relation failed",
	//		logx.Field("err", err))
	//	resp.Code = int32(0)
	//	resp.Msg = "mysql query record failed"
	//	return resp, err
	//}
	//if err == nil {
	//	// query success, already favorite
	//	resp.Code = int32(0)
	//	resp.Msg = "repeat operation"
	//	return resp, nil
	//}
	//// err == ErrNotFound
	//_, err = l.svcCtx.FavorRepo.FavoriteModel.Insert(l.ctx, &data.Favorite{
	//	UserId:  in.UserId,
	//	VideoId: in.VideoId,
	//})
	//if err != nil {
	//	logx.Errorw("add following relation failed",
	//		logx.Field("err", err),
	//	)
	//	resp.Code = int32(0)
	//	resp.Msg = "mysql insert failed"
	//	return resp, err
	//}
	resp.Code = int32(01)
	resp.Msg = "add follow relation success"
	return resp, nil
}
