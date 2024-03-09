package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetListByAuthorIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListByAuthorIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListByAuthorIdLogic {
	return &GetListByAuthorIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetListByAuthorId 根据用户id查询发布的所有视频信息，不需要再查询用户信息
func (l *GetListByAuthorIdLogic) GetListByAuthorId(in *pb.GetListByAuthorIdRequest) (*pb.GetListByAuthorIdResponse, error) {
	videos, err := l.svcCtx.Repo.GetVideosByAuthor(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	resp := new(pb.GetListByAuthorIdResponse)
	resp.VideoList = make([]*pb.Video, len(videos))
	for i, v := range videos {
		ifThumbup, err := l.svcCtx.FavorRpc.CheckThumbup(l.ctx, &favorite.CheckThumbupRequest{
			UserId:  in.UserId,
			VideoId: v.ID,
		})
		if err != nil {
			logx.Errorw("favor rpc check thumbup", logx.Field("err", err))
			return nil, err
		}
		resp.VideoList[i] = &pb.Video{
			Id:            v.ID,
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
