package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/user"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/model"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/video"

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
	logx.Infof("get videos %v\n", in.VideoIds)
	// 根据video ids获取列表
	length := len(in.VideoIds)
	if length == 0 {
		return nil, errors.New("empty set")
	}
	videoList := make([]*video.VideoDetail, length)
	for i, vid := range in.VideoIds {
		// 获取视频详细信息
		v, err := l.svcCtx.VideoRepo.GetVideoById(l.ctx, vid)
		if err != nil {
			return nil, err
		}
		// 根据author id获取用户详细信息
		authorRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{UserId: v.AuthorId})
		if err != nil {
			return nil, err
		}
		vd := &video.VideoDetail{
			Id: v.VideoId,
			Author: &video.UserInfo{
				Id:              authorRes.User.Id,
				Name:            authorRes.User.Name,
				Avatar:          authorRes.User.Avatar,
				BackgroundImage: authorRes.User.BackgroundImage,
				Signature:       authorRes.User.Signature,
				FavoriteCount:   authorRes.User.FavoriteCount,
				WorkCount:       authorRes.User.WorkCount,
				TotalFavorited:  authorRes.User.TotalFavorited,
				FollowerCount:   authorRes.User.FollowerCount,
				FollowCount:     authorRes.User.FollowCount,
				IsFollow:        authorRes.User.IsFollow,
			},
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    true, // 点赞
		}
		videoList[i] = vd
	}
	resp.VideoList = videoList
	return resp, nil
}
