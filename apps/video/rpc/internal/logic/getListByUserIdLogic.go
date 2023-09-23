package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"strconv"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GetListByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListByUserIdLogic {
	return &GetListByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetListByUserIdLogic) GetListByUserId(in *model.GetListByUserIdRequest) (*model.GetListByUserIdResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetListByUserIdResponse)
	// 从redis中查询user发布的video ids，
	vidRes, err := l.svcCtx.VideoCache.SMembers(l.ctx, "tiktok:video:user:"+fmt.Sprintf("%d", in.UserId)).Result()
	if err != nil || err != redis.Nil {
		return resp, err
	}
	if err == nil {

		// 根据ids从数据库查询详细信息
		val, err := mr.MapReduce(func(source chan<- string) {
			for _, vid := range vidRes {
				source <- vid
			}
		}, func(idStr string, writer mr.Writer[*model.Video], cancel func(error)) {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				cancel(err)
				return
			}
			v, err := l.svcCtx.VideoModel.FindOneByVideoId(l.ctx, id)
			if err != nil {
				cancel(err)
				return
			}
			writer.Write(v)

		}, func(pipe <-chan *model.Video, writer mr.Writer[[]*model.VideoInfo], cancel func(error)) {
			var result []*model.VideoInfo
			for p := range pipe {
				vf := &model.VideoInfo{
					VideoId:  p.VideoId,
					AuthorId: in.GetUserId(),
					PlayUrl:  p.PlayUrl,
					CoverUrl: p.CoverUrl,
					Title:    p.Title,
				}
				result = append(result, vf)
			}
			writer.Write(result)
		})

		if err != nil {
			logx.Errorw("[mr] get video list by ids failed",
				logx.Field("err", err))
			return resp, err
		}
		resp.VideoList = val
		return resp, nil
	}

	// redis 不存在video:user:uid 的 key，从mysql中读取并放入redis

	res, err := l.svcCtx.VideoModel.FindVideosByUserId(l.ctx, in.UserId)
	if err != nil {
		return resp, err
	}
	if len(res) == 0 {
		return resp, errors.New("empty list")
	}

	// 包装响应，并放入redis中
	resp.VideoList = make([]*video.VideoInfo, len(res))
	ids := make([]int64, len(res))
	for i, v := range res {
		ids[i] = v.Id
		vi := &video.VideoInfo{
			VideoId:  v.VideoId,
			AuthorId: v.AuthorId,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		}
		resp.VideoList[i] = vi
	}
	go func() {
		n, err := l.svcCtx.VideoCache.SAdd(l.ctx, "tiktok:video:user:"+fmt.Sprintf("%d", in.UserId), ids).Result()
		if err != nil {
			logx.Errorw("add into redis failed",
				logx.Field("err", err))
			return
		} else {
			logx.Info("%d records have been added", n)
		}
	}()

	return resp, nil
}
