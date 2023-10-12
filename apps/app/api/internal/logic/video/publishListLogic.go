package video

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line
	resp = new(types.PublishListResponse)

	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	//myId := claims.UserId
	//loginUserId := claims.UserId

	var wg sync.WaitGroup
	wg.Add(2)
	errCh := make(chan error, 2)
	// 1. 获取user信息
	authorRes := new(user.GetAuthorResponse)
	go func() {
		defer wg.Done()
		authorRes, err = l.svcCtx.UserRpc.GetAuthor(l.ctx, &user.GetAuthorRequest{
			UserId:   claims.UserId,
			AuthorId: req.UserId,
		})
		if err != nil {
			errCh <- err
			return
		}
	}()
	// 查询videos基础信息，以及点赞信息
	videosRes := new(video.GetListByAuthorIdResponse)
	go func() {
		wg.Done()
		videosRes, err = l.svcCtx.VideoRpc.GetListByAuthorId(l.ctx, &video.GetListByAuthorIdRequest{
			UserId: req.UserId,
		})
		if err != nil {
			errCh <- err
			return
		}
	}()
	wg.Wait()
	select {
	case err := <-errCh:
		logx.Error("concurrency rpc ", err)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = err.Error()
		return resp, nil
	default:
	}

	videoList := make([]types.Video, len(videosRes.VideoList))
	// TODO 拼接结果
	for i, vi := range videosRes.VideoList {
		v := types.Video{
			Id: vi.Id,
			Author: types.Author{
				Id:              authorRes.Id,
				Name:            authorRes.Name,
				Avatar:          authorRes.Avatar,
				BackgroundImage: authorRes.BackgroundImage,
				Signature:       authorRes.Signature,

				FollowCount:   authorRes.FollowCount,
				FollowerCount: authorRes.FollowerCount,
				IsFollow:      authorRes.IsFollow,

				FavoriteCount:  authorRes.FavoriteCount,
				TotalFavorited: authorRes.TotalFavorited,
				WorkCount:      authorRes.WorkCount,
			},
			Title:         vi.Title,
			CoverUrl:      vi.CoverUrl,
			PlayUrl:       vi.PlayUrl,
			FavoriteCount: vi.FavoriteCount,
			CommentCount:  vi.CommentCount,
			IsFavorite:    vi.IsFavorite,
		}
		videoList[i] = v
	}
	resp.VideoList = videoList

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
