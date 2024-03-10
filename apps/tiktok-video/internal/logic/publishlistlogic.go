package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublishList 查询指定用户发布的视频列表
func (l *PublishListLogic) PublishList(in *pb.PublishListRequest) (*pb.PublishListResponse, error) {
	// 查询发布视频列表
	videos, err := l.svcCtx.Repo.GetVideosByAuthor(l.ctx, in.UserId)
	if err != nil {
		logx.Errorw("get videos by author id",
			logx.Field("err", err),
			logx.Field("authorId", in.AuthorId))
		return nil, err
	}
	// 查询作者详细信息
	users, err := l.svcCtx.UserRpc.GetUsers(l.ctx, &userservice.GetUsersRequest{
		UserIds: []int64{in.AuthorId},
		UserId:  in.UserId,
	})
	if err != nil {
		logx.Errorw("user rpc get author info",
			logx.Field("err", err),
			logx.Field("authorId", in.AuthorId),
			logx.Field("userId", in.UserId))
		return nil, err
	} else if len(users.UserList) != 1 {
		logx.Errorw("user rpc ger users, incorrect response length", logx.Field("len", len(users.UserList)))
		return nil, errors.New("get users failed, incorrect response")
	}
	var author = &pb.User{
		Id:              users.UserList[0].Id,
		Name:            users.UserList[0].Name,
		FollowCount:     users.UserList[0].FollowCount,
		FollowerCount:   users.UserList[0].FollowerCount,
		IsFollow:        users.UserList[0].IsFollow,
		Avatar:          users.UserList[0].Avatar,
		BackgroundImage: users.UserList[0].BackgroundImage,
		Signature:       users.UserList[0].Signature,
		TotalFavorited:  users.UserList[0].TotalFavorited,
		WorkCount:       users.UserList[0].WorkCount,
		FavoriteCount:   users.UserList[0].FavoriteCount,
	}

	videoIds := make([]int64, len(videos))
	for i, v := range videos {
		videoIds[i] = v.ID
	}
	// 查询用户点赞情况
	thumbup, err := l.svcCtx.FavorRpc.MCheckThumbup(l.ctx, &favorite.MCheckThumbupRequest{
		UserId:   in.UserId,
		VideoIds: videoIds,
	})
	if err != nil {
		logx.Errorw("favor rpc MCheckThumbup", logx.Field("err", err),
			logx.Field("userId", in.UserId), logx.Field("videoIds", videoIds))
		return nil, err
	}

	resp := new(pb.PublishListResponse)
	resp.VideoList = make([]*pb.Video, len(videos))
	for i, v := range videos {
		resp.VideoList[i] = &pb.Video{
			Id:            v.ID,
			Author:        author,
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			Title:         v.Title,
			FavoriteCount: v.ThumbupCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    thumbup.IsThumbup[i],
		}
	}
	return resp, nil
}
