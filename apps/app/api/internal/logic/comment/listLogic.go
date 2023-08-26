package comment

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/comment"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.CommentListRequest) (resp *types.CommentListResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.CommentListResponse)
	//claims, err := jwtx.ParseToken(req.Token)
	//if err != nil {
	//	logx.Errorw("parse jwt token failed",
	//		logx.Field("err", err),
	//	)
	//	resp.StatusCode = http.StatusUnauthorized
	//	resp.StatusMsg = "invalid token"
	//	resp.CommentList = []types.Comment{}
	//	return resp, nil
	//}
	//uid := claims.UserId
	commentListRes := new(comment.GetCommentListResponse)
	commentListRes, err = l.svcCtx.CommentRpc.GetCommentList(l.ctx, &comment.GetCommentListRequest{
		VideoId: req.VideoId,
	})

	if err != nil {
		logx.Errorw("comment list rpc failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = "comment rpc err"
		resp.CommentList = []types.Comment{}
		return resp, nil
	}
	if commentListRes.CommentLen == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "empty list"
		resp.CommentList = []types.Comment{}
		return resp, nil
	}

	resp.CommentList = make([]types.Comment, commentListRes.CommentLen)
	for i, c := range commentListRes.CommentList {
		cmt := types.Comment{
			Id:         c.CommentId,
			Content:    c.Content,
			CreateDate: c.CreateTime,
		}
		userRes, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
			UserId: c.AuthorId,
		})
		if err != nil {
			logx.Errorw("get user by id failed",
				logx.Field("err", err),
			)
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = "user rpc err"
			resp.CommentList = []types.Comment{}
			return resp, nil
		}
		cmt.UserInfo = types.User{
			Id:              c.AuthorId,
			Name:            userRes.Name,
			Avatar:          userRes.Avatar,
			BackgroundImage: userRes.BackgroundImage,
			Signature:       userRes.Signature,
		}
		resp.CommentList[i] = cmt
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "success"
	return resp, nil
}
