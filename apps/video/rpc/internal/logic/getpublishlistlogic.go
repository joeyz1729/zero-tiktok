package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublishListLogic {
	return &GetPublishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetPublishList 根据用户id查询发布的所有视频信息，已经对应用户的点赞信息，不需要查询作者信息
func (l *GetPublishListLogic) GetPublishList(in *model.GetPublishListRequest) (*model.GetPublishListResponse, error) {
	videos, err := l.svcCtx.VideoRepo.GetVideosByAuthorId(l.ctx, in.AuthorId)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errors.New("empty set")
	}
	logx.Infof("get %d videos by tiktok-user %d\n", len(videos), in.AuthorId)
	resp := new(model.GetPublishListResponse)
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
		favorRes, err := l.svcCtx.FavorRpc.GetFavorite(l.ctx, &favorite.GetFavoriteRequest{UserId: in.UserId, VideoId: v.VideoId})
		if err != nil {
			logx.Error("get favorite relation ", err)
			continue
		} else if favorRes != nil {
			resp.VideoList[i].IsFavorite = favorRes.ActionType == int64(1)
		}
	}

	return resp, err

}
