package main

import (
	"fmt"
	"gin-blog-example/routers"
	"gin-blog-example/settings"
	"log"
	"net/http"
	"time"
)

func main() {
	// 初始化路由
	r := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", settings.HttpPort), //监听的 TCP 地址
		Handler:        r,                                     //http 句柄，实质为ServeHTTP，用于处理程序响应 HTTP 请求
		ReadTimeout:    time.Duration(settings.ReadTimeout),   //允许读取的最大时间
		WriteTimeout:   time.Duration(settings.WriteTimeout),  //允许读取请求头的最大时间
		MaxHeaderBytes: 1 << 20,                               //请求头的最大字节数
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalln("ListenAndServe failed！")
	}
}
