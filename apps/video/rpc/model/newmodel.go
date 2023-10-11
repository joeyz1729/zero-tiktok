package model

const (
	VideoInfoPrefix    = "tiktok:video:info:"    // +vid, hash	(info)
	VideoCountPrefix   = "tiktok:video:count:"   // +vid, hash	(count)
	VideoFeedPrefix    = "tiktok:video:feed::"   // +nil zset	(vid, timestamp)
	VideoPublishPrefix = "tiktok:video:publish:" // +uid set (vid)
)
