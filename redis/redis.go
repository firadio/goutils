package redis

import (
	"fmt"
	"time"

	red "github.com/gomodule/redigo/redis"
)

type Class struct {
	pool *red.Pool
}

func New(sHostPort string, sPassword string, iDatabase int) *Class {
	this := &Class{}
	this.pool = &red.Pool{
		MaxIdle:     256,
		MaxActive:   100,
		IdleTimeout: time.Duration(120),
		Dial: func() (red.Conn, error) {
			return red.Dial(
				"tcp",
				sHostPort, //redis.dev:6379
				red.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				red.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				red.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				red.DialPassword(sPassword), //密码
				red.DialDatabase(iDatabase), //数据库编号
			)
		},
	}
	return this
}

func RedisExec_Do(redis_pool *red.Pool, cmd string, key string, args ...string) (interface{}, error) {
	con := redis_pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}

func (redis *Class) LLEN(key string) int64 {
	ret, err := RedisExec_Do(redis.pool, "LLEN", key)
	if err != nil {
		fmt.Printf("LLEN(%s) %s\r\n", key, err)
		return 0
	}
	return ret.(int64)
}

func (redis *Class) LPUSH(key string, element ...string) int64 {
	ret, err := RedisExec_Do(redis.pool, "LPUSH", key, element...)
	if err != nil {
		fmt.Printf("LPUSH(%s) %s\r\n", key, err)
		return 0
	}
	return ret.(int64)
}

func (redis *Class) POP(key string, cmd string) string {
	ret, err := RedisExec_Do(redis.pool, cmd, key)
	if err != nil {
		fmt.Printf("%s(%s) %s\r\n", cmd, key, err)
		return ""
	}
	if ret == nil {
		return ""
	}
	return string(ret.([]byte))
}

func (redis *Class) LPOP(key string) string {
	return redis.POP(key, "LPOP")
}

func (redis *Class) RPOP(key string) string {
	return redis.POP(key, "RPOP")
}
