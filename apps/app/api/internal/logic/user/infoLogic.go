package user

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	// token 校验
	resp = new(types.UserInfoResponse)
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("parse token failed, or token is invalid",
			logx.Field("err", err))
		return resp, nil
	}
	userId := claims.UserId // 该用户id
	logx.Infof("userid:%d\n", userId)
	toUserId := req.UserId // 查询用户id
	resp.UserInfo = types.User{}
	errChan := make(chan error, 2)
	defer close(errChan)
	var wg sync.WaitGroup
	wg.Add(2)
	var userRes = &user.UserInfoResponse{}
	go func() {
		defer wg.Done()
		userRes, err = l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{UserId: toUserId})
		if err != nil {
			errChan <- err
		}
		resp.UserInfo.Id = toUserId
		resp.UserInfo.Name = userRes.User.Name
		resp.UserInfo.Avatar = userRes.User.Avatar
		resp.UserInfo.BackgroundImage = userRes.User.BackgroundImage
		resp.UserInfo.Signature = userRes.User.Signature
		resp.UserInfo.FollowCount = userRes.User.FollowCount
		resp.UserInfo.FollowerCount = userRes.User.FollowerCount
		resp.UserInfo.TotalFavorited = userRes.User.TotalFavorited
		resp.UserInfo.FavoriteCount = userRes.User.FavoriteCount

	}()
	var followRes = &follow.GetRelationResponse{}
	go func() {
		defer wg.Done()
		followRes, err = l.svcCtx.FollowRpc.GetRelation(l.ctx, &follow.GetRelationRequest{UserId: userId, ToUserId: toUserId})
		if err != nil {
			errChan <- err
		}
	}()
	wg.Wait()
	select {
	case result := <-errChan:
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "internal server error"
		logx.Error(result.Error())
		return nil, err
	default:
		resp.UserInfo.IsFollow = followRes.IsFollowing
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "success"
		return resp, nil
	}

}
