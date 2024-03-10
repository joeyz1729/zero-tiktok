package job

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/segmentio/kafka-go"
)

const (
	TopicUser      = "tiktok_user_user"
	TopicUserCount = "tiktok_user_user_count"
	TopicRelation  = "tiktok_relation_relation"
	TopicFavor     = "tiktok_thumbup_thumbup"
	TopicVideo     = "tiktok_video_video"

	GroupUpdateWorkCount     = "db_update_work_count"
	GroupUpdateFavorCount    = "db_update_user_favor_count"
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

func Start(ctx context.Context, c kafka.ReaderConfig, repo *repository.Repo) {
	go CreateUserWorker(ctx, c, repo).Start(ctx)
	go UserCountWorker(ctx, c, repo).Start(ctx)
	go RelationWorker(ctx, c, repo).Start(ctx)
	go VideoPublishWorker(ctx, c, repo).Start(ctx)
	go FavoriteCountWorker(ctx, c, repo).Start(ctx)
}
