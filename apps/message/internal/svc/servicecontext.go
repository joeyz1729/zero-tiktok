package svc

import (
	"github.com/jmoiron/sqlx"
	"github.com/joeyz1729/zero-tiktok/apps/message/internal/config"
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
