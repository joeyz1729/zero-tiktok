package es

import (
	"context"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository/dto"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

const (
	UserIndex = "tiktok_user"
)

func GetUser(userId int64, esClient *elasticsearch.TypedClient) (*dto.User, error) {
	resp, err := esClient.Get(UserIndex, strconv.Itoa(int(userId))).Do(context.TODO())
	if err != nil {
		logx.Errorw("[ESGetUser]", logx.Field("err", err), logx.Field("userId", userId))
		return nil, err
	}
	b, err := resp.Source_.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var detail dto.User
	extra.RegisterFuzzyDecoders()
	err = jsoniter.Unmarshal(b, &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func GetUserByName(username string, esClient *elasticsearch.TypedClient) (*dto.User, error) {
	resp, err := esClient.Search().
		Index(UserIndex).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Filter: []types.Query{
						{
							Term: map[string]types.TermQuery{
								"username": {Value: username},
							},
						},
					},
				},
			},
		}).Do(context.TODO())

	if err != nil {
		return nil, err
	}
	if len(resp.Hits.Hits) != 1 {
		return nil, errors.New("invalid record count")
	}
	str, err := resp.Hits.Hits[0].Source_.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var detail dto.User
	extra.RegisterFuzzyDecoders()
	err = jsoniter.Unmarshal(str, &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func CreateUser(ctx context.Context, data map[string]interface{}, esClient *elasticsearch.TypedClient) error {
	userId := data["id"].(string)
	_, err := esClient.Index(UserIndex).Id(userId).Document(data).Do(ctx)
	return err
}

func UpdateUserCount(ctx context.Context, data map[string]interface{}, esClient *elasticsearch.TypedClient) error {
	userId := data["id"].(string)
	_, err := esClient.Update(UserIndex, userId).Doc(data).Do(ctx)
	return err
}
