package job

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/segmentio/kafka-go"
)

const (
	TopicUser      = "tiktok_user_user"
	TopicUserCount = "tiktok_user_user_count"
	TopicRelation  = "tiktok_relation_relation"
	TopicVideo     = "tiktok_video_video"

	GroupUpdateWorkCount     = "db_update_work_count"
	GroupUpdateFavorCount    = "db_update_favor_count"
	GroupUpdateRelationCount = "db_update_relation_count"
	GroupCreateUser          = "es_create_user"
	GroupUpdateUserCount     = "es_update_user_count"
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
