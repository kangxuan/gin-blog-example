package gredis

import (
	"encoding/json"
	"gin-blog-example/settings"
	"github.com/gomodule/redigo/redis"
	"time"
)

// RedisConn 声明Redis的连接池
var RedisConn *redis.Pool

// SetUp 初始化Redis
func SetUp() error {
	RedisConn = &redis.Pool{
		MaxIdle:     settings.RedisSetting.MaxIdle,
		MaxActive:   settings.RedisSetting.MaxActive,
		IdleTimeout: settings.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) { //提供创建和配置应用程序连接的一个函数
			// 连接redis
			c, err := redis.Dial("tcp", settings.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			// 如果密码不为空则进行验证
			if settings.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", settings.RedisSetting.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { // 检查健康功能
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Set 设置缓存
func Set(key string, data interface{}, time int) error {
	// 从连接池中获取一个连接
	conn := RedisConn.Get()
	defer conn.Close()

	// 对值进行jsonEncode
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists 判断Key是否存在
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get 根据Key获取缓存内容
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Delete 删除key
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes 根据搜索删除key
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
