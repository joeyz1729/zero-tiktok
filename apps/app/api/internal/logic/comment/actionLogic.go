package comment

import (
	"context"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/comment"
	"github.com/YiZou89/zero-tiktok/pkg/jwtx"
	"net/http"

	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActionLogic) Action(req *types.CommentActionRequest) (resp *types.CommentActionResponse, err error) {
	// 需要返回评价内容和作者信息
	// 使用UserRpc和CommentRpc
	resp = new(types.CommentActionResponse)
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		logx.Errorw("jwt parse token failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusUnauthorized
		resp.StatusMsg = "invalid token"
		return resp, nil
	}

	if req.ActionType == int32(1) {
		_, err = l.svcCtx.CommentRpc.AddComment(l.ctx, &comment.AddCommentRequest{
			UserId:      claims.UserId,
			VideoId:     req.VideoId,
			CommentText: req.CommentText,
		})
		if err != nil {
			logx.Errorw("add comment rpc failed",
				logx.Field("err", err),
			)
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMsg = err.Error()
			return resp, nil
		}
		resp.StatusCode = http.StatusOK
		resp.StatusMsg = "add comment success"
		return resp, nil
	}

	// delete comment
	_, err = l.svcCtx.CommentRpc.DelComment(l.ctx, &comment.DelCommentRequest{
		UserId:    claims.UserId,
		VideoId:   req.VideoId,
		CommentId: req.CommentId,
	})
	if err != nil {
		logx.Errorw("del comment rpc failed",
			logx.Field("err", err),
		)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMsg = err.Error()
		return resp, nil
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMsg = "del comment success"

	return resp, nil
}
