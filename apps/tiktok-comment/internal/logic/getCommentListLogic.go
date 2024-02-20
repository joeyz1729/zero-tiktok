package logic

import (
	"context"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/pb"
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

func (l *GetCommentListLogic) GetCommentList(in *pb.GetCommentListRequest) (*pb.GetCommentListResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(pb.GetCommentListResponse)
	//var err error
	//
	//var cl []*repository.Comment
	//getStr := `select * from tiktok_comment.comment where video_id = ? order by create_time desc`
	//err = l.svcCtx.CommentDB.Select(&cl, getStr, in.VideoId)
	//if err != nil {
	//	logx.Errorw("mysql get comment list failed",
	//		logx.Field("err", err),
	//	)
	//	resp.CommentList = []*pb.CommentInfo{}
	//	return resp, err
	//}
	//if len(cl) == 0 {
	//	resp.CommentList = []*pb.CommentInfo{}
	//	return resp, nil
	//}
	//
	//resp.CommentList = make([]*pb.CommentInfo, len(cl))
	//for i, c := range cl {
	//	_, m, d := c.CreateTime.Date()
	//	resp.CommentList[i] = &pb.CommentInfo{
	//		CommentId:  c.CommentId,
	//		Content:    c.Content,
	//		CreateTime: fmt.Sprintf("%s-%s", m.String(), strconv.Itoa(d)),
	//	}
	//}
	return resp, nil
}
