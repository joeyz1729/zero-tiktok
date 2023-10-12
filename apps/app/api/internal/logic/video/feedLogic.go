package video

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
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
	// todo: add your logic here and delete this line
	resp = new(types.FeedResponse)
	resp.VideoList = []types.Video{}

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
	uid := claims.UserId
	var feedRes *video.FeedResponse
	feedRes, err = l.svcCtx.VideoRpc.Feed(l.ctx, &video.FeedRequest{
		UserId:     uid,
		LatestTime: req.LatestTime,
	})
	if err != nil {
		logx.Errorw("video rpc failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "rpc feed failed"
		resp.NextTime = time.Now().Unix()
		return resp, nil
	}
	if feedRes.VideoLen == 0 {
		resp.StatusMsg = "empty list"
		resp.NextTime = time.Now().Unix()
		resp.StatusCode = http.StatusOK
		return resp, nil
	}

	resp.VideoList = make([]types.Video, feedRes.VideoLen)
	for i, v := range feedRes.VideoList {
		vi := types.Video{
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
		resp.VideoList[i] = vi
	}
	resp.NextTime = feedRes.NextTime
	resp.StatusCode = http.StatusOK
	return resp, nil
}
