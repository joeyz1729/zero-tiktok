package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/model"
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
func (l *GetListByAuthorIdLogic) GetListByAuthorId(in *model.GetListByAuthorIdRequest) (*model.GetListByAuthorIdResponse, error) {
	videos, err := l.svcCtx.VideoRepo.GetVideosByAuthorId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errors.New("empty set")
	}
	logx.Infof("get %d videos by tiktok-user %d\n", len(videos), in.UserId)
	resp := new(model.GetListByAuthorIdResponse)
	resp.VideoList = make([]*model.VideoDetail, len(videos))
	for i, v := range videos {
		resp.VideoList[i] = &model.VideoDetail{
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
