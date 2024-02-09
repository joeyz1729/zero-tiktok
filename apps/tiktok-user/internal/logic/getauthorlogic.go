package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/follow"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuthorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAuthorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorLogic {
	return &GetAuthorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAuthorLogic) GetAuthor(in *pb.GetAuthorRequest) (*pb.GetAuthorResponse, error) {
	// todo: add your logic here and delete this line
	author, err := l.svcCtx.UserRepo.GetUserInfo(in.AuthorId)
	if err != nil {
		return nil, err
	}
	relation := false
	if in.UserId != in.AuthorId && in.UserId != 0 {
		followRes, err := l.svcCtx.FollowRpc.GetRelation(l.ctx, &follow.GetRelationRequest{
			UserId:   in.UserId,
			ToUserId: in.AuthorId,
		})
		if err != nil {
			return nil, err
		}
		relation = followRes.IsFollowing
	}

	return &pb.GetAuthorResponse{
		Id:              author.UserId,
		Name:            author.Username,
		Avatar:          "no avatar",
		BackgroundImage: "no background image",
		Signature:       "no signature",

		FavoriteCount:  author.FavoriteCount,
		TotalFavorited: author.TotalFavorited,
		WorkCount:      author.WorkCount,

		FollowerCount: author.FollowerCount,
		FollowCount:   author.FollowedCount,
		IsFollow:      relation,
	}, nil

}
