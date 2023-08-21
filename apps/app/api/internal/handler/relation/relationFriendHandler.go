package relation

import (
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/logic/relation"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RelationFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Douyin_relation_friend_list_request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := relation.NewRelationFriendLogic(r.Context(), svcCtx)
		resp, err := l.RelationFriend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
