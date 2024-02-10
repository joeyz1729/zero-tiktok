package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
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
	// todo: add your logic here and delete this line
	resp := new(pb.GetVideoByIdResponse)
	vi, err := repository.GetVideoById(l.ctx, in.VideoId)
	//vi, err := l.svcCtx.VideoModel.FindOneByVideoId(l.ctx, in.VideoId)
	if err != nil {
		logx.Errorw("mysql query failed",
			logx.Field("err", err),
		)
		return resp, err
	}

	resp.VideoInfo = &pb.Video{
		Id:       in.VideoId,
		Author:   nil,
		PlayUrl:  vi.PlayURL,
		CoverUrl: vi.CoverURL,
		Title:    vi.Title,
	}
	return resp, nil
}
