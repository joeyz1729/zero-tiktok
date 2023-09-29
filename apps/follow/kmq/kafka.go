package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/follow/dao/cache"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/YiZou89/zero-tiktok/apps/follow/dao"
	mysqldb "github.com/YiZou89/zero-tiktok/apps/follow/dao/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
)

var mdb *mysqldb.FollowDB
var rdb *cache.FollowCache

func main() {

	var c kq.KqConf
	conf.MustLoad("config.yaml", &c)
	var err error
	d, err := sqlx.Connect("mysql",
		"root:Zy_9908091729@tcp(localhost:3306)/tiktok_follow?parseTime=true&charset=utf8")
	if err != nil {
		panic(err)
	}
	defer d.Close()
	mdb = mysqldb.NewFollowDB(d)

	redisAddr := "127.0.0.1:6379"
	rd := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		//Password: c.RedisDB.Password,
		//DB:       c.RedisDB.DB,
		//PoolSize: c.RedisDB.PoolSize,
	})
	defer rd.Close()
	_, err = rd.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	rdb = cache.NewFollowCache(rd)

	q := kq.MustNewQueue(c, kq.WithHandle(func(k, v string) error {
		err = biz(v)
		if err != nil {
			logx.Error(err)
		}
		return nil
	}))
	defer q.Stop()
	q.Start()
}

func biz(v string) (err error) {
	fmt.Printf("=> %s\n", v)
	var relation dao.Action
	err = json.Unmarshal([]byte(v), &relation)
	if err != nil {
		logx.Errorw("json unmarshal data from kafka mq failed",
			logx.Field("err", err))
		return err
	}
	logx.Info(relation)
	if relation.ActionType == int32(1) {
		return addRelation(relation.UserId, relation.ToUserId)

	} else {
		return delRelation(relation.UserId, relation.ToUserId)
	}
}

func addRelation(uid, tid int64) (err error) {
	// TODO
	// 添加事务操作防止重复消费
	ok, err := mdb.CheckRelation(context.Background(), uid, tid)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	_, _, err = mdb.AddRelation(context.Background(), uid, tid)
	if err != nil {
		return err
	}
	err = rdb.DelRelation(context.Background(), uid, tid)
	if err != nil {
		return err
	}
	return nil
}

func delRelation(uid, tid int64) (err error) {
	_, _, err = mdb.DelRelation(context.Background(), uid, tid)
	if err != nil {
		return err
	}
	err = rdb.DelRelation(context.Background(), uid, tid)
	if err != nil {
		return err
	}
	return nil
}
