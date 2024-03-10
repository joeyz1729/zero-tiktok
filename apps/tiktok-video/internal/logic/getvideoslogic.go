package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/dto"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideosLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideosLogic {
	return &GetVideosLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideosLogic) GetVideos(in *pb.GetVideosRequest) (*pb.GetVideosResponse, error) {
	// 批量查询视频信息
	videos := make([]*dto.Video, len(in.VideoIds))
	for i, vid := range in.VideoIds {
		video, err := l.svcCtx.Repo.GetVideoById(l.ctx, vid)
		if err != nil {
			logx.Errorw("get video by id", logx.Field("err", err),
				logx.Field("videoId", vid))
			return nil, err
		}
		videos[i] = video
	}

	// 批量查询用户信息
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
		UserId:  in.UserId,
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

	// 批量查询点赞信息
	favor, err := l.svcCtx.FavorRpc.MCheckThumbup(l.ctx, &favorite.MCheckThumbupRequest{
		UserId:   in.UserId,
		VideoIds: in.VideoIds,
	})
	if err != nil {
		logx.Errorw("favor rpc check thumbup", logx.Field("err", err),
			logx.Field("userId", in.UserId), logx.Field("videoIds", in.VideoIds))
		return nil, err
	}

	// 组装返回结果
	resp := new(pb.GetVideosResponse)
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
			IsFavorite:    favor.IsThumbup[i],
		}
	}
	return resp, nil
}
