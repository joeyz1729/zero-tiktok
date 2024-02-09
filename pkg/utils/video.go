package utils

import (
	"context"
	"fmt"
	"github.com/joeyz1729/zero-tiktok/apps/video/rpc/mw/minio"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

func NewFileName(userId, time int64) string {
	return fmt.Sprintf("%d.%d", userId, time)
}

func URLConvert(ctx context.Context, s, h, p, path string) (fullURL string) {
	if len(path) == 0 {
		return ""
	}
	arr := strings.Split(path, "/")
	u, err := minio.GetObjURL(ctx, arr[0], arr[1])
	if err != nil {
		logx.Errorf("get obj url failed, err: %v\n", err)
		return ""
	}
	u.Scheme = s
	u.Host = h
	u.Path = p
	return u.String()
}
