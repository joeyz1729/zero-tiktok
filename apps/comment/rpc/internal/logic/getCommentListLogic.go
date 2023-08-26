package logic

import (
	"context"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/model"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentListLogic) GetCommentList(in *model.GetCommentListRequest) (*model.GetCommentListResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.GetCommentListResponse)
	var err error

	var cl []*model.Comment
	getStr := `select * from tiktok_comment.comment where video_id = ? order by create_time desc`
	err = l.svcCtx.CommentDB.Select(&cl, getStr, in.VideoId)
	if err != nil {
		logx.Errorw("mysql get comment list failed",
			logx.Field("err", err),
		)
		resp.CommentLen = 0
		resp.CommentList = []*model.CommentInfo{}
		return resp, err
	}
	if len(cl) == 0 {
		resp.CommentLen = 0
		resp.CommentList = []*model.CommentInfo{}
		return resp, nil
	}

	resp.CommentLen = int64(len(cl))
	resp.CommentList = make([]*model.CommentInfo, len(cl))
	for i, c := range cl {
		_, m, d := c.CreateTime.Date()
		resp.CommentList[i] = &model.CommentInfo{
			CommentId:  c.CommentId,
			AuthorId:   c.UserId,
			Content:    c.Content,
			CreateTime: fmt.Sprintf("%s-%s", m.String(), strconv.Itoa(d)),
		}
	}
	return resp, nil
}
