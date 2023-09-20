package video

import (
	"context"
	"errors"

	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"
	"sync"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListRequest) (resp *types.PublishListResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.PublishListResponse)

	_, err = jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	//loginUserId := claims.UserId

	// 并行请求
	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)
	defer close(errChan)

	// 1. 根据userid 查询 user 信息
	var author types.User
	go func() {
		defer wg.Done()
		userInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{UserId: req.UserId})
		if err != nil {
			logx.Errorw("get user by id failed",
				logx.Field("err", err),
			)
			errChan <- err
		}
		author = types.User{
			Id:              userInfo.GetId(),
			Name:            userInfo.GetName(),
			Avatar:          userInfo.GetAvatar(),
			BackgroundImage: userInfo.GetBackgroundImage(),
		}
	}()

	// 2. 从video模块获取video list
	go func() {
		defer wg.Done()
		res, err := l.svcCtx.VideoRpc.GetListByUserId(l.ctx, &video.GetListByUserIdRequest{
			UserId: req.UserId,
		})
		if err != nil {
			logx.Errorw("get list by user id failed or empty list",
				logx.Field("err", err),
			)
			errChan <- err
		} else if len(res.GetVideoList()) == 0 {
			errChan <- errors.New("empty video list")
		}

		resp.VideoList = make([]types.Video, len(res.GetVideoList()))

		for i, v := range res.GetVideoList() {
			var videoInfo types.Video
			videoInfo = types.Video{
				Id:       v.GetVideoId(),
				Author:   author,
				PlayUrl:  v.GetPlayUrl(),
				CoverUrl: v.GetCoverUrl(),
				Title:    v.GetTitle(),
			}
			resp.VideoList[i] = videoInfo
		}
	}()
	wg.Wait()

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
