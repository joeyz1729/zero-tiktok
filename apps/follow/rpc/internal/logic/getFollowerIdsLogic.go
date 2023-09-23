package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"strconv"

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
	var res []string
	var err error
	sqlStr := `select user_id from tiktok_follow.followed where followed_id = ? and id >= ? limit ?`
	err = l.svcCtx.FollowDB.Select(&res, sqlStr, in.GetUserId(), in.GetCursor(), in.GetPageSize()+1)
	if err != nil {
		logx.Errorw("mysql query failed",
			logx.Field("err", err),
		)
		resp.FollowerIds = []int64{}
		return resp, err
	}
	logx.Infof("get follower list success, %v", res)
	resp.FollowerIds = make([]int64, min(len(res), int(in.GetPageSize())))
	for i := range res {
		id, err := strconv.Atoi(res[i])
		if err != nil {
			resp.FollowerIds = []int64{}
			return resp, err
		}
		resp.FollowerIds[i] = int64(id)
	}
	logx.Info("get next cursor from mysql")
	if len(res) > int(in.GetPageSize()) {
		// 有下一页
		next, err := strconv.Atoi(res[len(res)-1])
		if err != nil {
			logx.Errorw("get next cursor failed",
				logx.Field("err", err),
			)
			resp.NextCursor = 0
			return resp, err
		}
		logx.Info("get next cursor success")
		resp.NextCursor = int64(next)
	}
	return resp, nil
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
