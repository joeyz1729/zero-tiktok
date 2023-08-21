package relation

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationFriendLogic {
	return &RelationFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationFriendLogic) RelationFriend(req *types.Douyin_relation_friend_list_request) (resp *types.Douyin_relation_friend_list_response, err error) {
	// todo: add your logic here and delete this line

	return
}
