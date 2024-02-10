package video

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/videoservice"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedRequest) (resp *types.FeedResponse, err error) {
	resp = new(types.FeedResponse)
	// 1. jwt
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("jwt parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusUnauthorized
		resp.StatusMsg = "invalid token"
		resp.NextTime = time.Now().Unix()
		return resp, nil
	}
	// 2. logic
	feedRes, err := l.svcCtx.VideoRpc.Feed(l.ctx, &videoservice.FeedRequest{
		UserId:     claims.UserId,
		LatestTime: req.LatestTime,
	})
	if err != nil {
		logx.Errorw("videoservice tiktok-user failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "tiktok-user feed failed"
		resp.NextTime = time.Now().Unix()
		return resp, nil
	}
	// 3. return
	resp.VideoList = make([]types.Video, len(feedRes.VideoList))
	for i, v := range feedRes.VideoList {
		resp.VideoList[i] = types.Video{
			Id: v.Id,
			Author: types.Author{
				Id:              v.Author.Id,
				Name:            v.Author.Name,
				Avatar:          v.Author.Avatar,
				BackgroundImage: v.Author.BackgroundImage,
				Signature:       v.Author.Signature,
				FollowCount:     v.Author.FollowCount,
				FollowerCount:   v.Author.FollowerCount,
				IsFollow:        v.Author.IsFollow,
				FavoriteCount:   v.Author.FavoriteCount,
				TotalFavorited:  v.Author.TotalFavorited,
				WorkCount:       v.Author.WorkCount,
			},
			Title:         v.Title,
			CoverUrl:      v.CoverUrl,
			PlayUrl:       v.PlayUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
		}
	}
	resp.NextTime = feedRes.NextTime
	resp.StatusCode = http.StatusOK
	return resp, nil
}
