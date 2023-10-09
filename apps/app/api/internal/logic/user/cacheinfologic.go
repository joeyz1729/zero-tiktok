package user

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"
	"sync"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CacheInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCacheInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CacheInfoLogic {
	return &CacheInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CacheInfoLogic) CacheInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	// todo: add your logic here and delete this line
	// token 校验
	resp = new(types.UserInfoResponse)
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed, or token is invalid",
			logx.Field("err", err))
		return resp, nil
	}
	userId := claims.UserId // 该用户id
	toUserId := req.UserId  // 查询用户id
	resp.UserInfo = types.User{}
	errChan := make(chan error, 2)
	defer close(errChan)
	var wg sync.WaitGroup
	wg.Add(2)

	// user 模块
	go func() {
		defer wg.Done()
		var userRes *user.UserInfoResponse
		userRes, err = l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{UserId: toUserId})

		//userRes, err = l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{UserId: toUserId})
		if err != nil {
			errChan <- err
		} else {
			resp.UserInfo.Id = toUserId
			resp.UserInfo.Name = userRes.User.Name
			resp.UserInfo.Avatar = userRes.User.Avatar
			resp.UserInfo.BackgroundImage = userRes.User.BackgroundImage
			resp.UserInfo.Signature = userRes.User.Signature

			resp.UserInfo.FollowCount = userRes.User.FollowCount
			resp.UserInfo.FollowerCount = userRes.User.FollowerCount

			resp.UserInfo.FavoriteCount = userRes.User.FavoriteCount
			resp.UserInfo.TotalFavorited = userRes.User.TotalFavorited
			resp.UserInfo.WorkCount = userRes.User.WorkCount
		}
	}()

	// 调用follow模块
	go func() {
		defer wg.Done()
		var followRes *follow.GetFollowCountResponse
		followRes, err = l.svcCtx.FollowRpc.GetFollowCount(l.ctx, &follow.GetFollowCountRequest{UserId: userId, ToUserId: toUserId})
		if err != nil {
			errChan <- err
		} else {
			//resp.UserInfo.FollowCount = followRes.FollowCount
			//resp.UserInfo.FollowerCount = followRes.FollowerCount
			resp.UserInfo.IsFollow = followRes.IsFollowing
		}
	}()

	wg.Wait()
	select {
	case result := <-errChan:
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		logx.Error(result.Error())
		return &types.UserInfoResponse{}, nil
	default:
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
