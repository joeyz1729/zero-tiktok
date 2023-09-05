package logic

import (
	"context"
	"strconv"

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
	var followerIds []string
	var err error
	followerIds, err = l.svcCtx.FollowCache.SMembers(l.ctx, "tiktok:follower:"+strconv.Itoa(int(in.UserId))).Result()
	if err != nil {
		logx.Errorw("[redis] get follower list failed",
			logx.Field("err", err),
		)
		resp.FollowerIds = []int64{}
		return resp, err
	}
	length := min(20, len(followerIds))
	resp.FollowerIds = make([]int64, length)
	for i := 0; i < length; i++ {
		id, err := strconv.Atoi(followerIds[i])
		if err != nil {
			break
		}
		resp.FollowerIds[i] = int64(id)
	}
	//sqlStr := `select user_id from tiktok_follow.follow where follower_id = ?`
	//err := l.svcCtx.FollowDB.Select(&followerIds, sqlStr, in.GetUserId())
	//if err != nil {
	//	logx.Errorw("mysql query failed",
	//		logx.Field("err", err),
	//	)
	//	resp.FollowerIds = []int64{}
	//	return resp, err
	//}
	return resp, nil
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
