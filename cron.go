package main

import (
	"gin-blog-example/pkg/logging"
	"github.com/robfig/cron"
	"time"
)

func main() {
	logging.Info("cron start...")

	// 根据本地时间创建一个新（空白）的 Cron job runner
	c := cron.New()
	// AddFunc 会向 Cron job runner 添加一个 func ，以按给定的时间表运行
	c.AddFunc("* * * * * *", func() {
		logging.Info("执行一次")
	})
	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for true {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
