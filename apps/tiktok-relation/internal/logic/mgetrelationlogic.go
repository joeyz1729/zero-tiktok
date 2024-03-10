package logic

import (
	"context"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MGetRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMGetRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MGetRelationLogic {
	return &MGetRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MGetRelationLogic) MGetRelation(in *pb.MGetRelationRequest) (*pb.MGetRelationResponse, error) {
	var resp = new(pb.MGetRelationResponse)
	resp.IsFollowed = make([]bool, len(in.ToUserIds))
	resp.IsFollowing = make([]bool, len(in.ToUserIds))
	for i, uid := range in.ToUserIds {
		following, err := l.svcCtx.FollowRepo.CheckRelation(l.ctx, in.UserId, uid)
		if err != nil {
			logx.Errorw("check relation", logx.Field("err", err),
				logx.Field("userId", in.UserId), logx.Field("toUserId", uid))
			return nil, err
		}
		resp.IsFollowing[i] = following
		followed, err := l.svcCtx.FollowRepo.CheckRelation(l.ctx, uid, in.UserId)
		if err != nil {
			logx.Errorw("check relation", logx.Field("err", err),
				logx.Field("userId", uid), logx.Field("toUserId", in.UserId))
			return nil, err
		}
		resp.IsFollowed[i] = followed
	}
	return resp, nil
}
