package logic

import (
	"context"
	"errors"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
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
	videos, err := repository.GetVideosByAuthorId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errors.New("empty set")
	}
	logx.Infof("get %d videos by tiktok-user %d\n", len(videos), in.UserId)
	resp := new(pb.GetListByAuthorIdResponse)
	resp.VideoList = make([]*pb.Video, len(videos))
	for i, v := range videos {
		resp.VideoList[i] = &pb.Video{
			Id:       v.ID,
			PlayUrl:  v.PlayURL,
			CoverUrl: v.CoverURL,
			Title:    v.Title,
		}
	}
	return resp, nil
}
