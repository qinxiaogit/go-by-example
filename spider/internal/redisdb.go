package internal

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/qinxiaogit/go-by-example/spider/config"
	"strconv"
	"time"
)

/**
useage:
RedisCliPool()
*/
var cliPool *redis.Pool
func NewRedisCliPool(maxIdle , maxActive , idleTimeOut int, host string, port int) *redis.Pool {
	return &redis.Pool{
		MaxIdle: maxIdle,
		MaxActive:maxActive,
		IdleTimeout:time.Duration(idleTimeOut),
		Dial:func()(redis.Conn,error){
			c,err := redis.Dial("tcp",host+":"+strconv.Itoa(port))
			if err!=nil{
				return nil,err
			}
			return c,nil
		},
	}	
}
//
func RedisCliPool()*redis.Pool{
	if cliPool!=nil{
		return cliPool
	}
	return NewRedisCliPool(10,100,20,config.GetConfig().Redis.Host,config.GetConfig().Redis.Port)
}

const (
	hotdataKey = "hotword:%s:%d"
)
//getHotWordKey
func GetHotWordKey(hot_type string,year int) string{
	return fmt.Sprintf(hotdataKey,hot_type,year)
}