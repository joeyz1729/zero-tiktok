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

type GetVideosByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideosByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideosByIdsLogic {
	return &GetVideosByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据视频id查询视频详细信息，不包括作者详细信息
func (l *GetVideosByIdsLogic) GetVideosByIds(in *pb.GetVideosByIdsRequest) (*pb.GetVideosByIdsResponse, error) {
	videos := make([]*dto.Video, len(in.VideoIds))
	for i, vid := range in.VideoIds {
		video, err := l.svcCtx.Repo.GetVideoById(l.ctx, vid)
		if err != nil {
			return nil, err
		}
		videos[i] = video
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

	resp := new(pb.GetVideosByIdsResponse)
	resp.VideoList = make([]*pb.Video, len(videos))
	for i, v := range videos {
		ifThumbup, err := l.svcCtx.FavorRpc.CheckThumbup(l.ctx, &favorite.CheckThumbupRequest{
			UserId:  in.UserId,
			VideoId: in.VideoIds[i],
		})
		if err != nil {
			logx.Errorw("favor rpc check thumbup", logx.Field("err", err))
			return nil, err
		}
		resp.VideoList[i] = &pb.Video{
			Id:            v.ID,
			Author:        userMap[v.ID],
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			Title:         v.Title,
			FavoriteCount: v.ThumbupCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    ifThumbup.IsThumbup == int32(1),
		}
	}
	return resp, nil
}
