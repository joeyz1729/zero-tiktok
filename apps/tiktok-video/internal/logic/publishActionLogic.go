package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/db"
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
	vid, err := snowflake.GenID()

	err = repository.AddVideo(l.ctx, &db.Video{
		ID:       int64(vid),
		AuthorID: in.GetUserId(),
		Title:    in.GetTitle(),
		PlayURL:  in.GetPlayUrl(),
		CoverURL: in.GetCoverUrl(),
	})

	if err != nil {
		logx.Errorw("[PublishAction] repo add video failed",
			logx.Field("err", err))
		return nil, err
	}
	resp.VideoId = int64(vid)
	return resp, nil
}
