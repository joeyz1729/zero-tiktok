package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/dao/cache"
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
	page, size := int(in.Page), int(in.Size)
	page, size = 1, 10
	start, end := (page-1)*size, page*size-1
	// redis cache
	res, err := l.svcCtx.FollowCache.ZRange(l.ctx, cache.FollowedPrefix+mid, int64(start), int64(end)).Result()
	if err != nil {
		logx.Errorw("[redis] get following list failed",
			logx.Field("err", err),
		)
		return resp, err
	}

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

	// use mysql database
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
