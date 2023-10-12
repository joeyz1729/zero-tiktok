package logic

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/dao"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/go-redis/redis/v8"
	"strconv"

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
	fl, err := l.svcCtx.VideoCache.ZRevRangeByScore(l.ctx, "tiktok:video:time", &redis.ZRangeBy{Max: strconv.Itoa(int(in.LatestTime)), Count: 30}).Result()
	if err != nil {
		logx.Errorw("[redis] get video by time failed",
			logx.Field("err", err),
		)
		return resp, err
	}
	logx.Info("[redis] get video ids success")
	if len(fl) == 0 {
		logx.Info("empty list")
		resp.VideoList = []*model.VideoInfo{}
		resp.VideoLen = 0
		return resp, nil
	}
	logx.Info("[mysql] get video info from database")

	//getStr := `select * from tiktok_video.video where create_time > ? order by create_time desc limit 30 `
	//err = l.svcCtx.VideoDB.Select(&fl, getStr, time.Unix(0, in.LatestTime))
	//if err != nil {
	//	logx.Errorw("get feed failed",
	//		logx.Field("err", err),
	//	)
	//	resp.VideoLen = int64(0)
	//	resp.VideoList = []*model.VideoInfo{}
	//	resp.NextTime = time.Now().UnixNano()
	//	return resp, err
	//}
	//if len(fl) == 0 {
	//	resp.VideoLen = int64(0)
	//	resp.VideoList = []*model.VideoInfo{}
	//	resp.NextTime = time.Now().UnixNano()
	//	return resp, err
	//}
	//
	resp.VideoLen = int64(len(fl))
	resp.VideoList = make([]*model.VideoInfo, len(fl))
	for i, vid := range fl {
		videoId, err := strconv.Atoi(vid)
		if err != nil {
			continue
		}
		logx.Infof("get video info, vid:%s", vid)
		v := new(dao.Video)
		getStr := `select * from tiktok_video.video where video_id = ?`
		err = l.svcCtx.VideoDB.Get(v, getStr, videoId)
		if err != nil {
			logx.Errorw("get video info failed",
				logx.Field("err", err),
			)
			continue
		}
		videoInfo := &model.VideoInfo{
			VideoId:  int64(videoId),
			AuthorId: v.AuthorId,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		}
		resp.VideoList[i] = videoInfo
	}

	//var nextTime time.Time
	//timeStr := `select create_time from tiktok_video.video order by create_time desc limit 1 offset ?`
	//err = l.svcCtx.VideoDB.Select(&nextTime, timeStr, len(fl)-1)
	//if err != nil {
	//	logx.Errorw("get next time failed",
	//		logx.Field("err", err),
	//	)
	//	resp.NextTime = time.Now().UnixNano()
	//	return resp, err
	//}
	//resp.NextTime = nextTime.UnixNano()

	logx.Info("[redis] get next timestamp")
	nextVid := fl[len(fl)-1]
	nextTime, err := l.svcCtx.VideoCache.ZScore(l.ctx, "tiktok:video:time", nextVid).Result()
	if err != nil {
		logx.Errorw("[mysql] get next timestamp failed",
			logx.Field("err", err),
		)
		return resp, err
	}
	resp.NextTime = int64(nextTime)
	return resp, nil
}
