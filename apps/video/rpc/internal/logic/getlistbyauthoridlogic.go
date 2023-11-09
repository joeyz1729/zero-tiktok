package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListByAuthorIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListByAuthorIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListByAuthorIdLogic {
	return &GetListByAuthorIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetListByAuthorIdLogic) GetListByAuthorId(in *model.GetListByAuthorIdRequest) (*model.GetListByAuthorIdResponse, error) {
	// todo: add your logic here and delete this line
	// 根据author id查询发布的视频列表
	resp := new(model.GetListByAuthorIdResponse)
	videos, err := l.svcCtx.VideoRepo.GetVideosByAuthorId(l.ctx, in.UserId)
	if err != nil {
		logx.Error("get videos by author id: ", err)
		return nil, err
	}
	videoList := make([]*model.VideoDetail, len(videos))
	for i, v := range videos {
		//go func(i int, v *data.Video) {
		videoList[i] = &model.VideoDetail{
			Id: v.VideoId,
			// author 不需要
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			Title:         v.Title,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    false,
		}
		//}(i, v)
	}
	resp.VideoList = videoList
	return resp, nil
}
