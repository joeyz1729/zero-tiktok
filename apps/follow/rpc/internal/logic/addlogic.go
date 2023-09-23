package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddLogic) Add(in *model.AddRequest) (*model.AddResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.AddResponse)
	var err error
	uid, tid := in.UserId, in.ToUserId

	//// 1. 添加到 bloom filter
	//bloomKey := data.BloomPrefix + uidStr + ":" + tidStr
	//err = l.svcCtx.Filter.AddCtx(l.ctx, []byte(bloomKey))
	//if err != nil {
	//	logx.Errorw("[bloom filter] add failed",
	//		logx.Field("err", err),
	//	)
	//	resp.Code = http.StatusInternalServerError
	//	resp.Msg = "add bloom filter failed"
	//	return resp, err
	//}
	//logx.Info("add bloom filter success")
	err = l.svcCtx.FollowDB.AddRelation(l.ctx, uid, tid)

	err = l.svcCtx.FollowCache.RemRelation(l.ctx, uid, tid)

	return resp, err
}
