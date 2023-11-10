package data

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type VideoDB interface {
	GetVideoById(vid int64) (*Video, error)

	GetVideoByUser(uid int64) ([]*Video, error)

	GetVideoIdsByAuthorId(uid int64) ([]int64, error)
	GetFeedIds(int64) ([]*VideoWithTime, error)
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
		VideoId:       vid,
		AuthorId:      videoInfo.AuthorId,
		Title:         videoInfo.Title,
		PlayUrl:       videoInfo.PlayUrl,
		CoverUrl:      videoInfo.CoverUrl,
		FavoriteCount: videoCount.FavoriteCount,
		CommentCount:  videoCount.CommentCount,
	}, nil

}

func (r *MysqlImpl) GetVideoByUser(uid int64) ([]*Video, error) {
	return nil, nil
}

func (r *MysqlImpl) GetVideoIdsByAuthorId(uid int64) ([]int64, error) {
	sqlStr := `select video_id from tiktok_video.video where author_id = ?`
	var videoIds []int64
	err := r.DB.Select(&videoIds, sqlStr, uid)
	if err != nil {
		return nil, err
	}
	return videoIds, nil
}

func (r *MysqlImpl) GetFeedIds(lastTime int64) ([]*VideoWithTime, error) {
	sqlStr := `select video_id, publish_time from tiktok_video.video where publish_time < ? limit 30`
	videos := make([]*VideoWithTime, 0, 30)
	err := r.DB.Select(&videos, sqlStr, time.Unix(lastTime, 0))
	if err != nil {
		return nil, err
	}
	return videos, nil
}
