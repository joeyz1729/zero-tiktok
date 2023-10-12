package dao

import (
	"github.com/jmoiron/sqlx"
)

type VideoDB interface {
	GetVideoById(vid int64) (*Video, error)

	GetVideoByUser(uid int64) ([]*Video, error)
}

type MysqlImpl struct {
	*sqlx.DB
}

func NewMysqlImpl(db *sqlx.DB) *MysqlImpl {
	return &MysqlImpl{db}
}

var _ VideoDB = (*MysqlImpl)(nil)

func (r *MysqlImpl) GetVideoById(vid int64) (video *Video, err error) {
	var videoInfo VideoInfo
	var videoCount VideoCount
	sqlStr := `select video_id, author_id, title, play_url, cover_url  from tiktok_video.video where video_id = ? limit 1`
	if err = r.Get(&videoInfo, sqlStr, vid); err != nil {
		return nil, err
	}
	sqlStr = `select favorite_count, comment_count from tiktok_video.video_count where video_id = ? limit 1`
	if err = r.Get(&videoCount, sqlStr, vid); err != nil {
		return nil, err
	}
	return &Video{
		vid,
		videoInfo.AuthorId,
		videoInfo.Title,
		videoInfo.PlayUrl,
		videoInfo.CoverUrl,
		videoCount.FavoriteCount,
		videoCount.CommentCount,
	}, nil

}

func (r *MysqlImpl) GetVideoByUser(uid int64) ([]*Video, error) {
	return nil, nil
}
