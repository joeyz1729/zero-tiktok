package video

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/follow"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/video"
	"github.com/joeyz1729/zero-tiktok/pkg/jwtx"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
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
	// token校验
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	logx.Infof("PublishList: myUserId: %v, toUserId: %v\n", claims.UserId, req.UserId)
	// 验证用户id是否正确，这里没有关注关系
	authorRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &userservice.UserInfoRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logx.Errorw("get author info failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	// 用户合法，查询关注关系和视频列表
	var wg sync.WaitGroup
	var errChan = make(chan error, 2)
	wg.Add(2)
	// 查询视频详细信息列表
	//videosRes := new(video.GetListByAuthorIdResponse)
	videosRes := new(video.GetPublishListResponse)
	go func() {
		wg.Done()
		videosRes, err = l.svcCtx.VideoRpc.GetPublishList(l.ctx, &video.GetPublishListRequest{
			//videosRes, err = l.svcCtx.VideoRpc.GetListByAuthorId(l.ctx, &video.GetListByAuthorIdRequest{
			UserId:   claims.UserId,
			AuthorId: req.UserId,
		})
		if err != nil {
			errChan <- err
			return
		}
	}()
	// 查询关注关系
	var followRes *follow.GetRelationResponse
	go func() {
		defer wg.Done()
		var err error
		followRes, err = l.svcCtx.FollowRpc.GetRelation(l.ctx, &follow.GetRelationRequest{
			UserId:   claims.UserId,
			ToUserId: req.UserId,
		})
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	select {
	case err := <-errChan:
		logx.Error("get follow relation, or get video list failed ", err)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = err.Error()
		return resp, nil
	default:
	}
	// 整理作者信息
	author := types.Author{
		Id:              authorRes.User.Id,
		Name:            authorRes.User.Name,
		Avatar:          authorRes.User.Avatar,
		BackgroundImage: authorRes.User.BackgroundImage,
		Signature:       authorRes.User.Signature,
		FollowCount:     authorRes.User.FollowCount,
		FollowerCount:   authorRes.User.FollowerCount,
		FavoriteCount:   authorRes.User.FavoriteCount,
		TotalFavorited:  authorRes.User.TotalFavorited,
		WorkCount:       authorRes.User.WorkCount,
		IsFollow:        followRes.IsFollowing,
	}
	// 整理返回结果

	resp.VideoList = make([]types.Video, len(videosRes.VideoList))
	for i, vi := range videosRes.VideoList {
		v := types.Video{
			Id:            vi.Id,
			Author:        author,
			Title:         vi.Title,
			CoverUrl:      vi.CoverUrl,
			PlayUrl:       vi.PlayUrl,
			FavoriteCount: vi.FavoriteCount,
			CommentCount:  vi.CommentCount,
			IsFavorite:    vi.IsFavorite,
		}
		resp.VideoList[i] = v
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
