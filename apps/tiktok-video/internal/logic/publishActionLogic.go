package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/pb"
	"github.com/joeyz1729/zero-tiktok/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishActionLogic) PublishAction(in *pb.PublishActionRequest) (*pb.PublishActionResponse, error) {
	resp := new(pb.PublishActionResponse)
	var err error
	// snowflake gen videoservice id
	vid, err := snowflake.GenID()

	err = l.svcCtx.VideoRepo.AddVideo(l.ctx, &repository.Video{
		VideoId:  int64(vid),
		AuthorId: in.GetUserId(),
		Title:    in.GetTitle(),
		PlayUrl:  in.GetPlayUrl(),
		CoverUrl: in.GetCoverUrl(),
	})

	if err != nil {
		logx.Errorw("[mysql] add videoservice failed",
			logx.Field("err", err))
		return nil, err
	}
	resp.VideoId = int64(vid)
	return resp, nil
}
