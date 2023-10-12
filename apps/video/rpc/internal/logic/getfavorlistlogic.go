package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"sync"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavorListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavorListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavorListLogic {
	return &GetFavorListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFavorList 获取用户点赞列表
func (l *GetFavorListLogic) GetFavorList(in *model.GetFavorListRequest) (*model.GetFavorListResponse, error) {
	var resp = new(model.GetFavorListResponse)
	// 根据video ids获取列表，
	length := len(in.VideoIds)
	var wg sync.WaitGroup
	wg.Add(length)
	errCh := make(chan error, length)
	videoList := make([]*video.VideoDetail, len(in.VideoIds))
	for i, vid := range in.VideoIds {
		go func(i int, vid int64) {
			defer wg.Done()
			v, err := l.svcCtx.VideoRepo.GetVideoById(l.ctx, vid)
			if err != nil {
				errCh <- err
				return
			}
			authorRes, err := l.svcCtx.UserRpc.GetAuthor(l.ctx, &user.GetAuthorRequest{UserId: in.UserId, AuthorId: v.AuthorId})
			if err != nil {
				errCh <- err
				return
			}
			//TODO
			vd := &video.VideoDetail{
				Id: v.VideoId,
				Author: &video.UserInfo{
					Id:              authorRes.Id,
					Name:            authorRes.Name,
					Avatar:          authorRes.Avatar,
					BackgroundImage: authorRes.BackgroundImage,
					Signature:       authorRes.Signature,
					FavoriteCount:   authorRes.FavoriteCount,
					WorkCount:       authorRes.WorkCount,
					TotalFavorited:  authorRes.TotalFavorited,
					FollowerCount:   authorRes.FollowerCount,
					FollowCount:     authorRes.FollowCount,
					IsFollow:        authorRes.IsFollow,
				},
				Title:         v.Title,
				PlayUrl:       v.PlayUrl,
				CoverUrl:      v.CoverUrl,
				FavoriteCount: v.FavoriteCount,
				CommentCount:  v.CommentCount,
				IsFavorite:    true,
			}
			videoList[i] = vd
		}(i, vid)
	}
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, err
	default:
	}
	resp.VideoList = videoList
	return resp, nil
}
