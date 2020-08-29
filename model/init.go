package model

import (
	"github.com/lun-zhang/gorm"
	"zlutils/mysql"
	"zlutils/redis"
)

var dbConn *gorm.DB
var ReConn *redis.Client

func Init() {
	//不关心值，只关心用值做事情，通常是初始化的时候
	//consul.GetJson("mysql", func(my mysql.Config) {
	//	dbConn = mysql.New(my)
	//})
	//consul.GetJson("redis", func(redisUrl string) {
	//	ReConn = redis.New(redisUrl)
	//})

	dbConn = mysql.New(mysql.Config{
		Url: "root:123@/zlexample?charset=utf8&parseTime=True&loc=Local",
	})
	ReConn = redis.New("redis://localhost:6379")
}
