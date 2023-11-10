package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/model"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
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
	gid := dtmgrpc.MustGenGid(DtmServer)
	msg := dtmgrpc.NewMsgGrpc(DtmServer, gid).Add(userServer+"", userReq).Add(videoServer+"", videoReq)
	db, err := sql.Open("mysql")

	if in.ActionType == int32(1) {
		err = l.svcCtx.FavorRepo.CreateFavoriteRecord(l.ctx, &model.Favorite{UserId: in.UserId, VideoId: in.VideoId})
	} else {
		err = l.svcCtx.FavorRepo.DeleteFavoriteRecord(l.ctx, &model.Favorite{UserId: in.UserId, VideoId: in.VideoId})
	}
	if err != nil {
		logx.Error("update favorite failed:", err)
		return nil, err
	}

	//rpc 更新计数
	var errCh = make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := l.svcCtx.UserRpc.UpdateFavoriteInfo(l.ctx, &user.UpdateFavoriteInfoRequest{
			UserId:     in.UserId,
			VideoId:    in.VideoId,
			ActionType: in.ActionType == int32(1),
			AuthorId:   in.AuthorId,
		})
		if err != nil {
			logx.Errorw("update user info",
				logx.Field("err", err))
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		_, err := l.svcCtx.VideoRpc.UpdateFavoriteCount(l.ctx, &video.UpdateFavoriteCountRequest{
			UserId:     in.UserId,
			VideoId:    in.VideoId,
			ActionType: in.ActionType == int32(1),
		})
		if err != nil {
			logx.Errorw("update video info",
				logx.Field("err", err))
			errCh <- err
		}
	}()
	wg.Wait()
	select {
	case <-errCh:
		return nil, err
	default:
	}
	resp.Code = http.StatusOK
	resp.Msg = "update record success"
	return resp, nil
}
