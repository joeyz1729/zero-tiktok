package dao

import (
	"github.com/jmoiron/sqlx"
)

type VideoDB interface {
	GetVideoById(vid int64) (Video, error)
	GetVideoByUser(uid int64) ([]Video, error)
}

type MysqlImpl struct {
	*sqlx.DB
}

func NewMysqlImpl(db *sqlx.DB) *MysqlImpl {
	return &MysqlImpl{db}
}

var _ VideoDB = (*MysqlImpl)(nil)

func (r *MysqlImpl) GetVideoById(vid int64) (video Video, err error) {
	sqlStr := `select author_id, title, play_url, cover_url  from tiktok_video.video where video_id = ? limit 1`
	if err = r.Get(&video, sqlStr, vid); err != nil {
		return Video{}, err
	}

	sqlStr = `select favorite_count, comment_count from tiktok_video.video_count where video_id = ? limit 1`
	if err = r.Get(&video, sqlStr, vid); err != nil {
		return Video{}, err
	}

	return video, nil
}

func (r *MysqlImpl) GetVideoByUser(uid int64) ([]Video, error) {
	return nil, nil
}
