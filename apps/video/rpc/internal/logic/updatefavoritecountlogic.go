package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/video/data"
	"strconv"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFavoriteCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteCountLogic {
	return &UpdateFavoriteCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFavoriteCountLogic) UpdateFavoriteCount(in *model.UpdateFavoriteCountRequest) (*model.UpdateFavoriteCountResponse, error) {
	// todo: add your logic here and delete this line

	var sqlStr1 string
	var err error
	if in.ActionType {
		sqlStr1 = `update tiktok_video.video_count set favorite_count = favorite_count + 1 where video_id = ? limit 1`
	} else {
		sqlStr1 = `update tiktok_video.video_count set favorite_count = favorite_count - 1 where video_id = ? limit 1`
	}
	if _, err = l.svcCtx.VideoDB.Exec(sqlStr1, in.VideoId); err != nil {
		logx.Error("db update favorite count ", err)
		return nil, err
	}

	// delete cache
	vidStr := strconv.FormatInt(in.VideoId, 10)
	if _, err = l.svcCtx.VideoCache.Del(context.Background(), data.VideoInfoPrefix+vidStr).Result(); err != nil {
		return nil, err
	}
	return &model.UpdateFavoriteCountResponse{}, nil
}
