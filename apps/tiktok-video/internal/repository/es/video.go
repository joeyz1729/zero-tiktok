package es

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository/dto"
	"strconv"
)

var (
	VideoIndex = "tiktok_video"
)

func CreateVideo(ctx context.Context, data map[string]interface{}, esClient *elasticsearch.TypedClient) error {
	videoId := data["id"].(string)
	_, err := esClient.Index(VideoIndex).Id(videoId).Document(data).Do(ctx)
	return err
}

func GetVideoById(ctx context.Context, vid int64, esClient *elasticsearch.TypedClient) (*dto.Video, error) {
	resp, err := esClient.Get(VideoIndex, strconv.FormatInt(vid, 10)).Do(ctx)
	data, err := resp.Source_.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var video dto.Video
	err = json.Unmarshal(data, &video)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func GetVideosByAuthor(ctx context.Context, uid int64, esClient *elasticsearch.TypedClient) ([]*dto.Video, error) {
	resp, err := esClient.Search().Index(VideoIndex).
		Query(&types.Query{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"author_id": {Query: strconv.FormatInt(uid, 10)},
			},
		}).Do(ctx)
	if err != nil {
		return nil, err
	}
	videos := make([]*dto.Video, len(resp.Hits.Hits))
	for i := range videos {
		data, err := resp.Hits.Hits[i].Source_.MarshalJSON()
		if err != nil {
			return nil, err
		}
		var video dto.Video
		if err = json.Unmarshal(data, &video); err != nil {
			return nil, err
		}
		videos[i] = &video
	}
	return videos, nil
}

func UpdateVideoCount(ctx context.Context, data map[string]interface{}) error {
	return nil
}
