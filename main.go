package main

import (
	"fmt"
	"gin-blog-example/models"
	"gin-blog-example/pkg/logging"
	"gin-blog-example/routers"
	"gin-blog-example/settings"
	"log"
	"net/http"
)

func main() {
	// 初始化路由
	settings.SetUp()
	models.SetUp()
	logging.SetUp()

	r := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", settings.ServerSetting.HttpPort), //监听的 TCP 地址
		Handler:        r,                                                   //http 句柄，实质为ServeHTTP，用于处理程序响应 HTTP 请求
		ReadTimeout:    settings.ServerSetting.ReadTimeout,                  //允许读取的最大时间
		WriteTimeout:   settings.ServerSetting.WriteTimeout,                 //允许读取请求头的最大时间
		MaxHeaderBytes: 1 << 20,                                             //请求头的最大字节数
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalln("ListenAndServe failed！")
	}

	//endless.DefaultReadTimeOut = time.Duration(settings.ReadTimeout)
	//endless.DefaultWriteTimeOut = time.Duration(settings.WriteTimeout)
	//endless.DefaultMaxHeaderBytes = 1 << 20
	//endPoint := fmt.Sprintf(":%d", settings.HttpPort)
	//
	//server := endless.NewServer(endPoint, routers.InitRouter())
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Fatalln("ListenAndServe failed！")
	//}

	//// 采用http.Server - Shutdown()完成优雅重启
	//// 初始化路由
	//r := routers.InitRouter()
	//
	//s := &http.Server{
	//	Addr:           fmt.Sprintf(":%d", settings.ServerSetting.HttpPort), //监听的 TCP 地址
	//	Handler:        r,                                                   //http 句柄，实质为ServeHTTP，用于处理程序响应 HTTP 请求
	//	ReadTimeout:    settings.ServerSetting.ReadTimeout,                  //允许读取的最大时间
	//	WriteTimeout:   settings.ServerSetting.WriteTimeout,                 //允许读取请求头的最大时间
	//	MaxHeaderBytes: 1 << 20,                                             //请求头的最大字节数
	//}
	//
	//// 开启协程
	//go func() {
	//	if err := s.ListenAndServe(); err != nil {
	//		log.Printf("Listen: %s\n", err)
	//	}
	//}()
	//
	//quit := make(chan os.Signal)
	//signal.Notify(quit, os.Interrupt)
	//<-quit
	//
	//log.Println("Shutdown Server ...")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := s.Shutdown(ctx); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	//
	//log.Println("Server exiting")
}
