package video

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/logic/video"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PublishActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := video.NewPublishActionLogic(r, r.Context(), svcCtx)
		resp, err := l.PublishAction()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
