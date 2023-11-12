package logic

import (
	"context"
	"errors"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/data"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"

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

// GetListByAuthorId 根据用户id查询发布的所有视频信息，不需要再查询用户信息
func (l *GetListByAuthorIdLogic) GetListByAuthorId(in *data.GetListByAuthorIdRequest) (*data.GetListByAuthorIdResponse, error) {
	videos, err := l.svcCtx.VideoRepo.GetVideosByAuthorId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errors.New("empty set")
	}
	if err != nil {
		return nil, err
	}

	resp := new(data.GetListByAuthorIdResponse)
	resp.VideoList = make([]*data.VideoDetail, len(videos))
	for i, v := range videos {
		resp.VideoList[i] = &data.VideoDetail{
			Id:            v.VideoId,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			Title:         v.Title,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			//IsFavorite:    false,
		}
	}
	return resp, nil
}
