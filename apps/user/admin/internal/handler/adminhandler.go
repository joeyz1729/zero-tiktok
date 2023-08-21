package handler

import (
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/user/admin/internal/logic"
	"github.com/YiZou89/zero-tiktok/apps/user/admin/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/user/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAdminLogic(r.Context(), svcCtx)
		resp, err := l.Admin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
