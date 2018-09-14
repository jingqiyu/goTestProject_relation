package main

import (
	"os"
	"os/signal"
	"relation/conf"
	"relation/db"
	"relation/handlers/user"
	"relation/logger"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
)

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGUSR1, syscall.SIGUSR2)
	go server()
	<-c
	db.CloseDB()
	//todo 释放资源相关逻辑放在这里处理
	os.Exit(0)

}

func server() {

	configor.Load(&conf.Config, "conf/config.yaml") //加载配置文件
	logger.Init()                                   //初始化日志文件

	//db.InitMysql(conf.Config) //初始化Mysql
	///db.InitRedisPool(conf.Config) //初始化Redis连接池
	//db.GetDB().SetLogger(logger.Info) //设置mysql日志输出

	router := gin.Default()
	rUser := router.Group("/user")
	{
		/*rUser.GET("", user.GetUserById)
		rUser.POST("add", user.CreateUserById)
		rUser.GET("all", user.GetUsers)
		rUser.GET("list", user.GetUsersSlice)
		rUser.GET("phone", user.GetUserByPhone)
		rUser.POST("login", user.Login)*/
		rUser.POST("test",user.ToCheck)
	}

	router.Run(":8088") // listen and serve on 0.0.0.0:8080
}
