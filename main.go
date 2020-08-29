package main

import (
	"flag"
	"runtime"

	"github.com/gin-gonic/gin"

	"fmt"
	"github.com/sirupsen/logrus"
	"zlexample/biz"
	"zlexample/com"
	"zlexample/model"
	"zlutils/bind"
	"zlutils/caller"
	"zlutils/code"
	"zlutils/guard"
	"zlutils/logger"
	"zlutils/metric"
	"zlutils/mysql"
	"zlutils/xray"
)

var (
	workerNum int
	ginMode   string
	address   string
	//consulAddress string
	//consulPrefix  string
)

const addr = ":12345" //TODO: 服务地址

func init() {
	flag.IntVar(&workerNum, "worker", runtime.NumCPU(), "runtime MAXPROCS value")
	flag.StringVar(&ginMode, "gm", gin.DebugMode, fmt.Sprintf("gin mode: %s %s %s", gin.DebugMode, gin.TestMode, gin.ReleaseMode))
	flag.StringVar(&address, "addr", addr, "server address")
	//flag.StringVar(&consulAddress, "ca", ":8500", "consul address")
	//flag.StringVar(&consulPrefix, "cp", "test/service/"+com.ProjectName, "consul prefix")
	flag.Parse()
	logrus.WithFields(logrus.Fields{
		"workerNum": workerNum,
		"ginMode":   ginMode,
		"address":   address,
		//"consulAddress": consulAddress,
		//"consulPrefix":  consulPrefix,
	}).Info()

	runtime.GOMAXPROCS(workerNum)
	gin.SetMode(ginMode)

	caller.Init(com.ProjectName)
	//consul.Init(consulAddress, consulPrefix)

	//非中间件监控等
	guard.DoBeforeCtx, guard.DoAfter = xray.DoBeforeCtx, xray.DoAfter //注释掉这一行后xray就不会记录函数（仍会记录sql、rpc）
	guard.InitDefaultMetric(com.ProjectName)
	logger.InitDefaultMetric(com.ProjectName)
	mysql.InitDefaultMetric(com.ProjectName)
}

func main() {
	//中间件等
	metricsPath := fmt.Sprintf("/%s/metrics", com.ProjectName)
	router := gin.New()
	router.Use(guard.Mid())                 //捕获中间件panic，万一中间件有问题，也能够响应
	router.GET(metricsPath, metric.Metrics) //忽略metrics的记录，因为太多且干扰视线
	router.Use(
		code.MidRespCounterErr(com.ProjectName),
		metric.MidRespCounterLatency(com.ProjectName),
		xray.Mid(
			com.ProjectName,
			nil,
			code.RespIsServerErr,
			code.RespIsClientErr),
		logger.MidDebug(),
		logger.MidInfo())
	router.Use(guard.Mid()) //捕获业务panic，保证xray能正常close、能正常打日志、正常响应、正常监控

	//TODO: 业务初始化代码，例如日志配置、数据库连接、默认配置等
	//logger.WatchByConsul("log_watch")
	logger.Init(logger.Config{Level: logrus.DebugLevel})
	model.Init()

	//绑定路由
	bindRouter(router)

	//服务启动
	err := router.Run(addr)
	//srv := endless.NewServer(address, router)//endless能优雅重启，但只有*nix上才能用
	//err := srv.ListenAndServe()
	logrus.Infof("shutdown %v", err)
}

//TODO 绑定路由
func bindRouter(router *gin.Engine) {
	base := router.Group(com.ProjectName, //项目名作为路径前缀
		code.MidRespWithTraceId(false), //响应trace_id帮助定位请求
	)
	//指定不同权限的路径，基本上就是4种：站内、站外、内部调用、管理后台
	{
		app := base.Group("app",
			code.MidRespWithErr(true), //测试环境响应携带错误详细信息，正式环境不响应，帮助排查bug
			//session.MidUser(),         //检查并提取用户信息
		)
		app.POST("book/add/:name", bind.Wrap(biz.AppAddBook))
	}

	{
		admin := base.Group("admin",
			code.MidRespWithErr(false), //响应携带错误详细信息
			//session.MidOperator(),      //检查并提取操作者信息
		)
		admin.GET("book/get/:id", bind.Wrap(model.BookGet))
		admin.GET("book/list", bind.Wrap(model.BookList))
		admin.POST("book/add", bind.Wrap(model.BookAdd))
		admin.POST("book/edit/:id", bind.Wrap(model.BookEdit))
		admin.POST("book/delete/:id", bind.Wrap(model.BookDelete))
	}
	{
		outside := base.Group("outside",
			code.MidRespWithErr(true), //测试环境响应携带错误详细信息，正式环境不响应，帮助排查bug
		)
		outside.GET("hi", bind.Wrap(biz.Outside)) //站外例子比较简单
	}
	{
		rpc := base.Group("rpc",
			code.MidRespWithErr(false), //响应携带错误详细信息
		)
		rpc.GET("book", bind.Wrap(biz.RpcGetBook))

	}
}
