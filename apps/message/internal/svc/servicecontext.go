package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/message/internal/config"
	"github.com/jmoiron/sqlx"
)

type ServiceContext struct {
	Config config.Config

	MessageDB *sqlx.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sqlx.Connect("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		MessageDB: db,
	}
}
