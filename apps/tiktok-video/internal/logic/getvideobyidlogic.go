package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoByIdLogic {
	return &GetVideoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoByIdLogic) GetVideoById(in *pb.GetVideoByIdRequest) (*pb.GetVideoByIdResponse, error) {
	resp := new(pb.GetVideoByIdResponse)
	videoDetail, err := l.svcCtx.Repo.GetVideoById(l.ctx, in.VideoId)
	if err != nil {
		logx.Errorw("repo get video", logx.Field("err", err))
		return nil, err
	}
	users, err := l.svcCtx.UserRpc.GetUsers(l.ctx, &userservice.GetUsersRequest{
		UserIds: []int64{videoDetail.AuthorID},
	})
	if err != nil || users == nil || len(users.UserList) != 0 {
		logx.Errorw("rpc get users", logx.Field("err", err))
		return nil, err
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

	resp.VideoInfo = &pb.Video{
		Id:            videoDetail.ID,
		Author:        author,
		PlayUrl:       videoDetail.PlayURL,
		CoverUrl:      videoDetail.CoverURL,
		Title:         videoDetail.Title,
		FavoriteCount: videoDetail.ThumbupCount,
		CommentCount:  videoDetail.CommentCount,
	}
	return resp, nil
}
