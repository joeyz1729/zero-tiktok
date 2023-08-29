package logic

import (
	"context"

	"github.com/YiZou89/zero-tiktok/apps/message/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/message/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ActionLogic) Action(in *model.ActionRequest) (*model.ActionResponse, error) {
	// todo: add your logic here and delete this line
	resp := new(model.ActionResponse)
	var err error
	insertStr := `insert into tiktok_message.message(user_id, to_user_id, content) value(?, ?, ?)`
	_, err = l.svcCtx.MessageDB.Exec(insertStr, in.UserId, in.ToUserId, in.Content)
	if err != nil {
		logx.Errorw("insert message record failed",
			logx.Field("err", err),
		)
		resp.Msg = "insert message record failed"
		return resp, err
	}

	resp.Msg = "success"
	return resp, nil
}
