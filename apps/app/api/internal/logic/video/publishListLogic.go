package video

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"

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

	// 1. 从video模块获取video list
	res, err := l.svcCtx.VideoRpc.GetListByUserId(l.ctx, &video.GetListByUserIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logx.Errorw("get list by user id failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		return resp, nil
	}
	if len(res.GetVideoList()) == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "no video data"
		return resp, nil
	}

	resp.VideoList = make([]types.Video, len(res.GetVideoList()))

	for i, v := range res.GetVideoList() {
		// 2. 根据video list内容从user模块获取author的详细信息
		userInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{UserId: v.GetAuthorId()})
		if err != nil {
			logx.Errorw("get user by id failed",
				logx.Field("err", err),
			)
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "internal server error"
			return resp, nil
		}
		// 3. 从like模块读取点赞信息

		// 4. 从comment模块读取评论信息

		var videoInfo types.Video
		videoInfo = types.Video{
			Id: v.GetVideoId(),
			Author: types.User{
				Id:              userInfo.GetId(),
				Name:            userInfo.GetName(),
				Avatar:          userInfo.GetAvatar(),
				BackgroundImage: userInfo.GetBackgroundImage(),
			},
			PlayUrl:  v.GetPlayUrl(),
			CoverUrl: v.GetCoverUrl(),
			Title:    v.GetTitle(),
		}
		resp.VideoList[i] = videoInfo
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
