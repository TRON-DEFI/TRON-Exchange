package task

import (
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/buffer"
	"github.com/wlcy/tradehome-service/buffer2"
	redisClient "github.com/wlcy/tradehome-service/common/redis"
	"gopkg.in/redis.v5"
)

func Async() {
	go buffer.GetPirceInfBuffer()
	go buffer.StartKChartTimer()
	go buffer2.GetPrecisionBuffer()
	buffer.StartTimer()
}

func lockKey(key string) bool {
	value, err := redisClient.RedisCli.Get(key).Result()
	if err == redis.Nil {
		return true
	} else if err != nil {
		log.Errorf(err, "lock redis get value error")
		return false
	}
	if value == "" {
		return true
	} else {
		return false
	}
}
