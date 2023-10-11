package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

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

func (l *GetVideosByIdsLogic) GetVideosByIds(in *model.GetVideosByIdsRequest) (*model.GetVideosByIdsResponse, error) {
	// todo: add your logic here and delete this line
	//resp := new(model.GetVideosByIdsResponse)
	// 先查询基础信息
	videos, err := l.svcCtx.VideoRepo.GetFavorLists(l.ctx, in.VideoIds)
	if err != nil {
		return nil, err
	}
	// user rpc 根据userId和authorId添加关注关系信息
	authorIds := make([]int64, len(videos))
	for i, video := range videos {
		authorIds[i] = video.AuthorId
	}
	//l.svcCtx.UserRpc.Get

	// favorite rpc 根据userId和videoId添加点赞关系信息
	// 因为是查询的是点赞列表所以不需要

	return &model.GetVideosByIdsResponse{}, nil
}
