package logic

import (
	"context"
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
	videos, err := l.svcCtx.Repo.GetVideosByAuthor(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	var (
		userIds   = make([]int64, 0)
		userIdMap = make(map[int64]struct{})
		userMap   = make(map[int64]*pb.User)
	)
	for _, video := range videos {
		userIdMap[video.AuthorID] = struct{}{}
	}
	for uid := range userIdMap {
		userIds = append(userIds, uid)
	}
	users, err := l.svcCtx.UserRpc.GetUsers(l.ctx, &userservice.GetUsersRequest{
		UserIds: userIds,
	})
	if err != nil {
		logx.Errorw("user rpc get users", logx.Field("err", err))
		return nil, err
	}
	for _, user := range users.UserList {
		userMap[user.Id] = &pb.User{
			Id:              user.Id,
			Name:            user.Name,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			IsFollow:        user.IsFollow,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
		}
	}

	resp := new(pb.PublishListResponse)
	resp.VideoList = make([]*pb.Video, len(videos))
	for i, v := range videos {
		resp.VideoList[i] = &pb.Video{
			Id:            v.ID,
			Author:        userMap[v.ID],
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			Title:         v.Title,
			FavoriteCount: v.ThumbupCount,
			CommentCount:  v.CommentCount,
		}
	}
	return resp, nil
}
