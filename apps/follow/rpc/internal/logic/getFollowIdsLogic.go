package logic

import (
	"context"
	"strconv"

	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowIdsLogic {
	return &GetFollowIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowIdsLogic) GetFollowIds(in *model.GetFollowIdsRequest) (*model.GetFollowIdsResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetFollowIdsResponse)
	mid := strconv.Itoa(int(in.UserId))
	res, nextCursor, err := l.svcCtx.FollowCache.SScan(l.ctx, "tiktok:following:"+mid, 0, "", 1).Result()
	if err != nil {
		logx.Errorw("[redis] get following list failed",
			logx.Field("err", err),
		)
		return resp, err
	}
	logx.Infof("[redis] get following list success, next cursor: %d", nextCursor)
	resp.FollowIds = make([]int64, len(res))
	for i := range res {
		id, err := strconv.Atoi(res[i])
		if err != nil {
			resp.FollowIds = []int64{}
			return resp, err
		}
		resp.FollowIds[i] = int64(id)
	}
	return resp, nil
	//var followIds []int64
	//sqlStr := `select follower_id from tiktok_follow.follow where user_id = ?`
	//err := l.svcCtx.FollowDB.Select(&followIds, sqlStr, in.GetUserId())
	//if err != nil {
	//	logx.Errorw("mysql query failed",
	//		logx.Field("err", err),
	//	)
	//	resp.FollowIds = []int64{}
	//	return resp, err
	//}
	//resp.FollowIds = followIds
	//return resp, nil
}
