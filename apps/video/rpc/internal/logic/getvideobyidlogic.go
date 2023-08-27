package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

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

func (l *GetVideoByIdLogic) GetVideoById(in *model.GetVideoByIdRequest) (*model.GetVideoByIdResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetVideoByIdResponse)
	vi, err := l.svcCtx.VideoModel.FindOneByVideoId(l.ctx, in.VideoId)
	if err != nil {
		logx.Errorw("mysql query failed",
			logx.Field("err", err),
		)
		return resp, nil
	}

	resp.VideoInfo = &model.VideoInfo{
		VideoId:  in.VideoId,
		AuthorId: vi.AuthorId,
		PlayUrl:  vi.PlayUrl,
		CoverUrl: vi.CoverUrl,
		Title:    vi.Title,
	}
	return resp, nil
}
