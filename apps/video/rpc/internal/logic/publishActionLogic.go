package logic

import (
	"context"
	"database/sql"
	"github.com/YiZou89/zero-tiktok/pkg/snowflake"
	"time"

	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishActionLogic) PublishAction(in *model.PublishActionRequest) (*model.PublishActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.PublishActionResponse)
	// snowflake gen video id
	vid, err := snowflake.GenID()
	// gen url

	_, err = l.svcCtx.VideoModel.Insert(l.ctx, &model.Video{
		VideoId:  int64(vid),
		AuthorId: in.GetUserId(),
		Title:    in.GetTitle(),
		Data: sql.NullString{
			String: string(in.GetData()),
			Valid:  true,
		},
		PlayUrl:     "need cover url",
		CoverUrl:    "need cover url",
		PublishTime: time.Now(),
	})
	if err != nil {

		return resp, err
	}

	resp.VideoId = int64(vid)
	return resp, nil

}
