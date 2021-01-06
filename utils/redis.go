package utils

import (
	"time"

	red "github.com/gomodule/redigo/redis"
)

var redis_pool *red.Pool

func InitRedis(sHostPort string, sPassword string, iDatabase int) {
	redis_pool = &red.Pool{
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
				red.DialDatabase(iDatabase), //数据库编号：2
			)
		},
	}
}

func RedisExec_Do(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
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

func RedisExec_ZIncr(name string, member string) {
	key := "go:" + name
	//redis_conn.Do("ZINCRBY", key, "1", member)
	//redis_client.ZIncrBy(ctx, key, 1, member)
	RedisExec_Do("ZINCRBY", key, "1", member)
}

func Redis_ipcount_service(service string, userIpAddr string) {
	iTimestamp := TimestampInt64()
	if false {
		// 单独IP的访问频次时间分布记录不具有分析的意义
		RedisExec_ZIncr("count:userip_allsrv:"+userIpAddr, itoa64(iTimestamp))
	}
	if false {
		// 统计每个service被访问的次数（按分钟分组）
		RedisExec_ZIncr("count:service_allip:time60:"+itoa64(iTimestamp/60), service)
		RedisExec_ZIncr("count:service_allip:time300:"+itoa64(iTimestamp/300), service)
		RedisExec_ZIncr("count:service_allip:time600:"+itoa64(iTimestamp/600), service)
		RedisExec_ZIncr("count:service_allip:time900:"+itoa64(iTimestamp/900), service)
		RedisExec_ZIncr("count:service_allip:time1800:"+itoa64(iTimestamp/1800), service)
	}
	if true {
		// 统计每个service被访问的次数（按小时分组）
		RedisExec_ZIncr("count:service_allip:time3600:"+itoa64(iTimestamp/3600), service)
		// 统计每个ipaddr被访问的次数（按小时分组）
		RedisExec_ZIncr("count:ipaddr_allsrv:time3600:"+itoa64(iTimestamp/3600), userIpAddr)
	}
}
