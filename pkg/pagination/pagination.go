package pagination

import (
	"encoding/base64"
	"encoding/json"
)

type Page struct {
	NextId        uint64 `json:"next_id"`
	NextTimeAtUTC int64  `json:"next_time_at_utc"`
	PageSize      int64  `json:"page_size"`
}

type PageToken string

func (p Page) Encode() PageToken {
	b, err := json.Marshal(p)
	if err != nil {
		return PageToken("")
	}
	return PageToken(base64.StdEncoding.EncodeToString(b))
}

func (t PageToken) Decode() Page {
	var result Page
	if len(t) == 0 {
		return Page{}
	}
	bytes, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return Page{}
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return Page{}
	}
	return result
}
