package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationLogic {
	return &GetRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationLogic) GetRelation(in *model.GetRelationRequest) (*model.GetRelationResponse, error) {
	// todo: add your logic here and delete this line
	var resp = new(model.GetRelationResponse)
	var err error
	logx.Info("begin get follow relation")
	// 1. bloom filter
	//bloomKey := cache.BloomPrefix + mid + ":" + fid
	//ok, err := l.svcCtx.Filter.ExistsCtx(l.ctx, []byte(bloomKey))
	//if err != nil {
	//	logx.Errorw("bloom filter failed",
	//		logx.Field("err", err),
	//	)
	//	return resp, err
	//}
	//if !ok {
	//	logx.Info("bloom filter not exist, return")
	//	resp.IfFollowing = 0
	//	resp.IfFollower = 11111
	//	return resp, nil
	//}
	// 2. check redis
	ok, err := l.svcCtx.FollowCache.GetRelation(l.ctx, in.UserId, in.ToUserId)
	if err != nil && ok {
		resp.IfFollowing = int32(1)
		return resp, nil
	}
	// miss
	ok, err = l.svcCtx.FollowDB.CheckRelation(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		return resp, err
	}
	if ok {
		resp.IfFollowing = int32(1)
	}

	go l.svcCtx.FollowCache.AddRelation(l.ctx, in.UserId, in.ToUserId, false, 0, 0)
	return resp, nil
}
