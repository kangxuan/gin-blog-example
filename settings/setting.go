package settings

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg *ini.File

	RunMode string

	HttpPort     int // HTTP_PORT
	ReadTimeout  int // 读超时
	WriteTimeout int // 写超时

	JwtSecret string
	PageSize  int
)

func init() {
	// 利用ini加载配置
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
}

func LoadBase() {
	// Section是分区，如果没有分区则为默认""
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HttpPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = int(time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second)
	WriteTimeout = int(time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second)
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}
