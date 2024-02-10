package video

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/logic/video"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	maxSize = 10 << 20
)

func PublishActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(maxSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}

		file, header, err := r.FormFile("repository")
		defer func() { _ = file.Close() }()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		fileData, err := io.ReadAll(file)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}

		l := video.NewPublishActionLogic(r, r.Context(), svcCtx)
		resp, err := l.PublishAction(
			r.FormValue("token"),
			header.Filename,
			r.FormValue("title"),
			fileData)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
