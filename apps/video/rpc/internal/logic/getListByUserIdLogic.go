package logic

import (
	"context"
	"errors"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListByUserIdLogic {
	return &GetListByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetListByUserIdLogic) GetListByUserId(in *model.GetListByUserIdRequest) (*model.GetListByUserIdResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetListByUserIdResponse)
	res, err := l.svcCtx.VideoModel.FindVideosByUserId(l.ctx, in.UserId)
	if err != nil {
		return resp, err
	}
	if len(res) == 0 {
		return resp, errors.New("no video data")
	}
	resp.VideoList = make([]*video.VideoInfo, len(res))
	for i, v := range res {
		vi := &video.VideoInfo{
			VideoId:  v.VideoId,
			AuthorId: v.AuthorId,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		}
		resp.VideoList[i] = vi
	}

	return resp, nil
}
