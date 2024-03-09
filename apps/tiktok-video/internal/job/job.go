package job

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
	"github.com/segmentio/kafka-go"
)

const (
	TopicVideo      = "tiktok_video_video"
	TopicVideoCount = "tiktok_video_video_count"
	TopicFavor      = "tiktok_thumbup_thumbup"
	TopicComment    = "tiktok_comment_comment"
)

type Worker struct {
	Repo         *repository.Repo
	ReaderConfig kafka.ReaderConfig
}

type Msg struct {
	Type     string                   `json:"type"`
	Database string                   `json:"database"`
	Table    string                   `json:"table"`
	IsDdl    bool                     `json:"isDdl"`
	Data     []map[string]interface{} `json:"data"`
}
