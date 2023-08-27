package logic

import (
	"context"
	"time"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FeedLogic) Feed(in *model.FeedRequest) (*model.FeedResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.FeedResponse)
	var err error

	var fl []*model.VideoInfo
	getStr := `select * from tiktok_video.video where create_time > ? order by create_time desc limit 30 `

	err = l.svcCtx.VideoDB.Select(&fl, getStr, time.Unix(0, in.LatestTime))
	if err != nil {
		logx.Errorw("get feed failed",
			logx.Field("err", err),
		)
		resp.VideoLen = int64(0)
		resp.VideoList = []*model.VideoInfo{}
		resp.NextTime = time.Now().UnixNano()
		return resp, err
	}
	if len(fl) == 0 {
		resp.VideoLen = int64(0)
		resp.VideoList = []*model.VideoInfo{}
		resp.NextTime = time.Now().UnixNano()
		return resp, err
	}

	resp.VideoLen = int64(len(fl))
	resp.VideoList = make([]*model.VideoInfo, len(fl))
	for i, v := range fl {
		videoInfo := &model.VideoInfo{
			VideoId:  v.VideoId,
			AuthorId: v.AuthorId,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		}
		resp.VideoList[i] = videoInfo
	}

	var nextTime time.Time
	timeStr := `select create_time from tiktok_video.video order by create_time desc limit 1 offset ?`
	err = l.svcCtx.VideoDB.Select(&nextTime, timeStr, len(fl)-1)
	if err != nil {
		logx.Errorw("get next time failed",
			logx.Field("err", err),
		)
		resp.NextTime = time.Now().UnixNano()
		return resp, err
	}
	resp.NextTime = nextTime.UnixNano()
	return resp, nil
}
