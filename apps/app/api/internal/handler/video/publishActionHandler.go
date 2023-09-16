package video

import (
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/logic/video"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
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
