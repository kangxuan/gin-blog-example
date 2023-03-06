package settings

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg *ini.File
	err error
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int           // 最大空闲连接数
	MaxActive   int           // 在给定时间内，允许分配的最大连接数（当为0时，没有限制）
	IdleTimeout time.Duration // 在给定时间内将会保持空闲状态，达到这个时间限制则会关闭连接（当为0时，没有限制）
	ExpireTime  int           // 通用的到期时间
}

var RedisSetting = &Redis{}

func SetUp() {
	// 利用ini加载配置
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadServer()
	LoadApp()
	LoadDatabase()
	LoadRedis()
}

// LoadServer 加载ServerSetting
func LoadServer() {
	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}
	// 这两句一定要加上，因为这里是按纳秒传进去的，需要乘以秒的单位，如果按纳秒服务在纳秒之内不返回数据则会失败。
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second
}

// LoadApp 加载AppSetting
func LoadApp() {
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}
}

// LoadDatabase 加载DatabaseSetting
func LoadDatabase() {
	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}
}

// LoadRedis 加载RedisSetting
func LoadRedis() {
	err = Cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
