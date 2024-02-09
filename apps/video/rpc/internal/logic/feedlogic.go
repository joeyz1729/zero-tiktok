package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/user"
	"sync"
	"time"

	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FeedLogic) Feed(in *model.FeedRequest) (*model.FeedResponse, error) {
	resp := new(model.FeedResponse)
	// 1 根据last_time查询video id列表，注意要限制查询条数
	videoIds, nextTime, err := l.svcCtx.VideoRepo.FeedIds(l.ctx, in.LatestTime)
	if err != nil {
		logx.Error("get feed ids ", err)
		return nil, err
	}
	if len(videoIds) == 0 {
		return &model.FeedResponse{VideoLen: 0, NextTime: time.Now().Unix(), VideoList: []*model.VideoDetail{}}, nil
	}
	// 可以并发请求，多条video不相关
	feedList := make([]*model.VideoDetail, len(videoIds))
	var wg sync.WaitGroup
	errCh := make(chan error, len(videoIds))
	wg.Add(len(videoIds))
	for i, vid := range videoIds {
		go func(i int, vid int64) {
			defer wg.Done()
			// 2 根据video ids查询video详细信息，
			vi, err := l.svcCtx.VideoRepo.GetVideoById(l.ctx, vid)
			if err != nil {
				errCh <- err
				return
			}
			// 2->3 根据author ids查询author详细信息，然后根据user id和author id查询关注信息
			author, err := l.svcCtx.UserRpc.GetAuthor(l.ctx, &user.GetAuthorRequest{UserId: in.UserId, AuthorId: vi.AuthorId})
			if err != nil {
				errCh <- err
				return
			}
			// 2 根据vid uid查询点赞信息
			favor, err := l.svcCtx.FavorRpc.GetFavorite(l.ctx, &favorite.GetFavoriteRequest{UserId: in.UserId, VideoId: vid})
			if err != nil {
				errCh <- err
				return
			}
			v := &model.VideoDetail{
				Id: vid,
				Author: &model.UserInfo{
					Id:              author.Id,
					Name:            author.Name,
					Avatar:          author.Avatar,
					BackgroundImage: author.BackgroundImage,
					Signature:       author.Signature,

					FollowCount:   author.FollowCount,
					FollowerCount: author.FollowerCount,
					IsFollow:      author.IsFollow,

					FavoriteCount:  author.FavoriteCount,
					WorkCount:      author.WorkCount,
					TotalFavorited: author.TotalFavorited,
				},
				Title:         vi.Title,
				PlayUrl:       vi.PlayUrl,
				CoverUrl:      vi.CoverUrl,
				FavoriteCount: vi.FavoriteCount,
				CommentCount:  vi.CommentCount,
				IsFavorite:    favor.ActionType == int64(1),
			}
			feedList[i] = v
		}(i, vid)
	}
	wg.Wait()
	select {
	case err := <-errCh:
		logx.Error("concurrency query video ", err)
		return nil, err
	default:
	}

	resp.VideoLen = int64(len(feedList))
	resp.VideoList = feedList
	resp.NextTime = nextTime - 1
	return resp, nil
}
