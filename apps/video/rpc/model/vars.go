package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

const (
	VideoFeedPrefix    = "tiktok:video::feed"    // +nil zset	(vid, timestamp)
	VideoInfoPrefix    = "tiktok:video:info:"    // +vid, hash	(info)
	VideoPublishPrefix = "tiktok:video:publish:" // +uid set (vid)

	FieldInfoTitle     = "title"
	FieldInfoPlayUrl   = "playurl"
	FieldInfoCoverUrl  = "coverurl"
	FieldInfoAuthorId  = "authorid"
	FieldCountFavorite = "favorite"
	FieldCountComment  = "comment"
)
