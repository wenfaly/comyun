package main

import (
	"comyun/dao/mongo"
	"comyun/dao/mysql"
	"comyun/dao/redis"
	"comyun/logger"
	"comyun/pkg/snowflake"
	"comyun/routes"
	"comyun/settings"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){
	//项目脚手架配置
	//配置文件
	var err error
	if err = settings.Init();err!= nil{
		fmt.Printf("Init settings failed, err:%v\n",err)
		return
	}

	//初始化日志
	if err = logger.Init("web.mode");err!= nil{
		fmt.Printf("Init logger failed, err:%v\n",err)
		return
	}
	defer zap.L().Sync()

	//初始化Mysql连接
	if err = mysql.Init();err!= nil{
		fmt.Printf("Init mysql failed, err:%v\n",err)
		return
	}
	defer mysql.Close()

	//初始化Mongo连接
	var ctx context.Context
	if err,ctx = mongo.Init();err != nil{
		fmt.Printf("Init mongo failed, err:%v\n",err)
		return
	}
	defer mongo.Close(ctx)

	//初始化Redis连接
	if err = redis.Init(); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	//雪花算法初始化
	//后续如果有其他算法统一进行初始化
	if err = snowflake.Init(
		viper.GetString("web.start_time"),
		int64(viper.GetInt("web.machine_id")),
		);err != nil{
		fmt.Printf("Init snowflake failed, err:%v\n",err)
		return
	}

	//路由设置
	r := routes.SetUp()
	//优雅启动
	server := &http.Server{
		Addr: fmt.Sprintf(":%d",viper.GetInt("web.port")),
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			zap.L().Error("listen:"+err.Error())
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ")
	}

	zap.L().Info("Server exiting")
}
