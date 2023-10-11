package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoIdsLogic {
	return &GetVideoIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoIdsLogic) GetVideoIds(in *model.GetVideoIdsRequest) (*model.GetVideoIdsResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetVideoIdsResponse)
	//var videoIds []int64
	//
	//sqlStr := `select video_id from tiktok_favorite.favorite where user_id = ?`
	//err := l.svcCtx.FavorRepo.FavorDB.Select(&videoIds, sqlStr, in.UserId)
	//if err != nil {
	//	logx.Errorw("select video ids failed",
	//		logx.Field("err", err),
	//	)
	//	return resp, err
	//}
	//resp.VideoNum = int64(len(videoIds))
	//resp.VideoIds = videoIds
	return resp, nil
}
