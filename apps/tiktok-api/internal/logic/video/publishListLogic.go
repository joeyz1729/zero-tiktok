package video

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/videoservice"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListRequest) (resp *types.PublishListResponse, err error) {
	resp = new(types.PublishListResponse)
	// 1. jwt
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	// 2. logic
	videosRes, err := l.svcCtx.VideoRpc.PublishList(l.ctx, &videoservice.PublishListRequest{
		UserId:   claims.UserId,
		AuthorId: req.UserId,
	})
	if err != nil {
		logx.Errorw("[PublishList]",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	// 3. return
	resp.VideoList = make([]types.Video, len(videosRes.VideoList))
	for i, v := range videosRes.VideoList {
		resp.VideoList[i] = types.Video{
			Id: v.Id,
			Author: types.Author{
				Id:              v.Author.Id,
				Name:            v.Author.Name,
				FollowCount:     v.Author.FollowCount,
				FollowerCount:   v.Author.FollowerCount,
				IsFollow:        v.Author.IsFollow,
				Avatar:          v.Author.Avatar,
				BackgroundImage: v.Author.BackgroundImage,
				Signature:       v.Author.Signature,
				TotalFavorited:  v.Author.TotalFavorited,
				WorkCount:       v.Author.WorkCount,
				FavoriteCount:   v.Author.FavoriteCount,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			Title:         v.Title,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
		}
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
