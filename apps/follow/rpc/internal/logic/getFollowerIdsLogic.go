package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerIdsLogic {
	return &GetFollowerIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowerIdsLogic) GetFollowerIds(in *model.GetFollowerIdsRequest) (*model.GetFollowerIdsResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetFollowerIdsResponse)
	var followerIds []int64
	sqlStr := `select user_id from tiktok_follow.follow where follower_id = ?`
	err := l.svcCtx.FollowDB.Select(&followerIds, sqlStr, in.GetUserId())
	if err != nil {
		logx.Errorw("mysql query failed",
			logx.Field("err", err),
		)
		resp.FollowerIds = []int64{}
		return resp, err
	}
	resp.FollowerIds = followerIds
	return resp, nil
}
