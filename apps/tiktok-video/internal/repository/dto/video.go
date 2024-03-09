package dto

type Video struct {
	ID           int64  `json:"id,string"`
	AuthorID     int64  `json:"author_id,string"`
	Title        string `json:"title"`
	PlayURL      string `json:"play_url"`
	CoverURL     string `json:"cover_url"`
	ThumbupCount int64  `json:"thumbup_count,string"`
	CommentCount int64  `json:"comment_count,string"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}
