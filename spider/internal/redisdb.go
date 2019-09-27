package internal

import (
	"github.com/gomodule/redigo/redis"
)

/**
useage:
RedisCliPool()
*/
var cliPool *redis.Pool
//NewRdisCliPool
type NewRedisCliPool(maxIdle, maxActive, idleTimeOut int, host string, port int) *redis.Pool {
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
		}

	}	
}
//
func RedisCliPool()*redis.Pool{
	if cliPool!=nil{
		return clicliPool
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