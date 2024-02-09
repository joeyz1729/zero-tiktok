package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/joeyz1729/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/favorite/rpc/model"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/user"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/video"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

var (
	ErrRepeatedOperation = errors.New("repeated operation")
	DtmServer            = "etcd://localhost:2379/dtmservice"
)

type ActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ActionLogic) Action(in *model.ActionRequest) (*model.ActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.ActionResponse)
	// 检查关系是否存在
	exist, err := l.svcCtx.FavorRepo.CheckFavor(l.ctx, in.UserId, in.VideoId)
	if err != nil {
		logx.Error("is favorite video failed: " + err.Error())
		return nil, err
	}
	// 是否重复
	if exist && in.ActionType == int32(1) || !exist && in.ActionType == int32(2) {
		return nil, ErrRepeatedOperation
	}
	// 分布式事务更新开始
	userServer, err := l.svcCtx.Config.UserRpc.BuildTarget()
	if err != nil {
		return nil, err
	}
	videoServer, err := l.svcCtx.Config.VideoRpc.BuildTarget()
	if err != nil {
		return nil, err
	}
	userReq := &user.UpdateFavoriteInfoRequest{
		UserId:     in.UserId,
		ActionType: in.ActionType == int32(1),
		AuthorId:   in.AuthorId,
	}
	videoReq := &video.UpdateFavoriteCountRequest{
		UserId:     in.UserId,
		VideoId:    in.VideoId,
		ActionType: in.ActionType == int32(1),
	}
	// 生成唯一事务id
	gid := dtmgrpc.MustGenGid(DtmServer)
	// 添加分布式事务中的调用
	msg := dtmgrpc.NewMsgGrpc(DtmServer, gid).
		Add(userServer+"/tiktok-user.User/UpdateFavoriteInfo", userReq).
		Add(videoServer+"/video.Video/UpdateFavoriteCount", videoReq)
	// 本地调用
	err = msg.DoAndSubmitDB(userServer+"/tiktok-user.User/FollowPrepare", l.svcCtx.BarrierDB, func(tx *sql.Tx) error {
		var e error
		if in.ActionType == int32(1) {
			e = l.svcCtx.FavorRepo.CreateFavoriteRecord(l.ctx, &model.Favorite{UserId: in.UserId, VideoId: in.VideoId})
		} else {
			e = l.svcCtx.FavorRepo.DeleteFavoriteRecord(l.ctx, &model.Favorite{UserId: in.UserId, VideoId: in.VideoId})
		}
		if e != nil {
			logx.Error("update favorite failed:", err)
		}
		return e
	})
	if err != nil {
		logx.Error("dtm transaction failed", err)
		return nil, err
	}
	resp.Code = http.StatusOK
	resp.Msg = "update record success"
	return resp, nil
}
