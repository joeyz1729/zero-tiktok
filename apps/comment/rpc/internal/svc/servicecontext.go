package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/config"
	"github.com/zeromicro/go-queue/kq"

	//"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/kmq"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/model"
	"github.com/jmoiron/sqlx"
	sqlz "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	CommentModel model.CommentModel

	CommentDB *sqlx.DB

	//CommentMq queue.MessageQueue
	//CommentMq *amqp.Connection
	//CommentMq *kmq.Kmq

	KafkaPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlz.NewMysql(c.Mysql.DataSource)
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	//dial, err := kmq.NewConnection(c)
	//addMq, delMq, err := kmq.NewMq()
	//if err != nil {
	//	panic(err)
	//}
	//q := kq.MustNewQueue(c.KafkaMq, kq.WithHandle(func(k, v string) error {
	//	fmt.Printf("=> %s\n, v")
	//	return nil
	//}))
	//q := kmq.NewMq(c.KafkaMq, db)

	return &ServiceContext{
		Config:       c,
		CommentModel: model.NewCommentModel(sqlConn),
		CommentDB:    db,
		//CommentMq:    q,
		KafkaPusher: kq.NewPusher(c.KafkaMq.Brokers, c.KafkaMq.Topic),
	}
}
