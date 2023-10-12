package dao

type Video struct {
	VideoId  int64  `db:"video_id"`
	AuthorId int64  `db:"author_id"`
	Title    string `db:"title"`
	PlayUrl  string `db:"play_url"`
	CoverUrl string `db:"cover_url"`

	FavoriteCount int64 `db:"favorite_count"`
	CommentCount  int64 `db:"comment_count"`
}

type VideoDetail struct {
	VideoInfo
	VideoCount
}

type VideoInfo struct {
	VideoId  int64  `db:"video_id"`
	AuthorId int64  `db:"author_id"`
	Title    string `db:"title"`
	PlayUrl  string `db:"play_url"`
	CoverUrl string `db:"cover_url"`
}

type VideoCount struct {
	FavoriteCount int64 `db:"favorite_count"`
	CommentCount  int64 `db:"comment_count"`
}
