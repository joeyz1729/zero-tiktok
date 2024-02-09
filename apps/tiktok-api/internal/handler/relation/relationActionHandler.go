package relation

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/logic/relation"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RelationActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RelationActionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := relation.NewRelationActionLogic(r.Context(), svcCtx)
		resp, err := l.RelationAction(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
