package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/util"
	"time"
)

//GetOrderPool 从redis中获取订单
func GetOrderPool() map[int64]*OrderListRedis {
	if !redisExists(OrderListRedisKey) {
		log.Debugf("GetOrderPool no such key in redis:[%v]", OrderListRedisKey)
		return nil
	}
	val := redisGet(OrderListRedisKey)
	log.Debugf("GetOrderPool get data length = [%v]", len(val))
	if len(val) == 0 {
		return nil
	}

	orderList := make(map[int64]*OrderListRedis, 0)
	err := json.Unmarshal([]byte(val), &orderList)
	if nil != err || len(orderList) == 0 {
		log.Errorf(err, "Unmarshal redis content error, reids key:[%v]", OrderListRedisKey)
		return nil
	}
	log.Debugf("GetOrderPool done")
	return orderList
}

//GetOrderList 获取redis中的订单池
func GetOrderList() map[int64]*OrderListRedis {
	if !redisExists(OrderListRedisKey) {
		log.Debugf("GetOrderList no such key in redis:[%v]", OrderListRedisKey)
		return nil
	}
	val := redisGet(OrderListRedisKey)
	log.Debugf("GetOrderList get data length = [%v]", len(val))
	if len(val) == 0 {
		return nil
	}
	orderList := make(map[int64]*OrderListRedis, 0)
	err := json.Unmarshal([]byte(val), &orderList)
	if nil != err || len(orderList) == 0 {
		log.Errorf(err, "Unmarshal redis content failed:[%v], redis key:[%v]", err, OrderListRedisKey)
		return nil
	}
	return orderList
}

//GetOrderList 获取redis中的订单池
func GetOrderList2(pairType int64) map[int64]*OrderListRedis {
	var val string
	if 1 == pairType {
		if !redisExists(OrderListRedisKey10) {
			log.Debugf("GetOrderList no such key in redis:[%v]", OrderListRedisKey10)
			return nil
		}
		val = redisGet(OrderListRedisKey10)
		log.Debugf("GetOrderList get trc10 data length = [%v]", len(val))
		if len(val) == 0 {
			return nil
		}
	} else if 2 == pairType {
		if !redisExists(OrderListRedisKey) {
			log.Debugf("GetOrderList no such key in redis:[%v]", OrderListRedisKey)
			return nil
		}
		val = redisGet(OrderListRedisKey)
		log.Debugf("GetOrderList get trc20 data length = [%v]", len(val))
		if len(val) == 0 {
			return nil
		}
	}
	orderList := make(map[int64]*OrderListRedis, 0)
	err := json.Unmarshal([]byte(val), &orderList)
	if nil != err || len(orderList) == 0 {
		log.Errorf(err, "Unmarshal redis content failed:[%v], redis key:[%v]", err, OrderListRedisKey)
		return nil
	}
	return orderList
}

//GetExchangeRealPrice 从redis中获取订单价格信息
func GetExchangeRealPrice(exchangeID int64) *ExchangeIDRealPrice {
	key := fmt.Sprintf(ExchangeIDRealPriceRedisKey, exchangeID)
	if !redisExists(key) {
		log.Debugf("GetExchangeRealPrice no such key in redis:[%v]", key)
		return nil
	}
	val := redisGet(key)
	log.Debugf("GetExchangeRealPrice get data length = [%v]", len(val))
	if len(val) == 0 {
		return nil
	}

	exchangeRealPrice := &ExchangeIDRealPrice{}
	err := json.Unmarshal([]byte(val), &exchangeRealPrice)
	if nil != err || exchangeRealPrice == nil {
		log.Errorf(err, "Unmarshal redis content failed:[%v], redis key:[%v]", err, key)
		return nil
	}
	if exchangeRealPrice.RealPrice == "" {
		exchangeRealPrice.RealPrice = "0"
	}
	log.Debugf("GetExchangeRealPrice get value = [%#v]", exchangeRealPrice)
	return exchangeRealPrice
}

//UpdateExchangeRealPrice 更新实时价格等信息
func UpdateExchangeRealPrice(zeroPrice, realPrice int64, exchangeID, amount int64) error {
	log.Debugf("UpdateExchangeRealPrice start update exchangeID: [%v],realPrice: [%v],zeroPrice: [%v],amount: [%v]", exchangeID, realPrice, zeroPrice, amount)
	key := fmt.Sprintf(ExchangeIDRealPriceRedisKey, exchangeID)
	if redisExists(key) {
		exchangeRealPrice := GetExchangeRealPrice(exchangeID)
		if exchangeRealPrice == nil {
			log.Error("get exchangeRealPrice in redis nil", nil)
			return errors.New("get exchangeRealPrice in redis nil")
		}
		exchangeRealPrice.Amount = exchangeRealPrice.Amount + amount
		exchangeRealPrice.RealPrice = fmt.Sprintf("%v", realPrice)
		if zeroPrice > 0 {
			exchangeRealPrice.ZeroPrice = fmt.Sprintf("%v", zeroPrice)
		}
		buffer, err := json.Marshal(exchangeRealPrice)
		if err != nil {
			log.Errorf(err, "Marshal error")
			return err
		}
		//写入redis 并不设置有效期
		cmd := RedisCli.Set(key, string(buffer), 0)
		if cmd.Err() != nil {
			log.Errorf(cmd.Err(), "set key:[%v] value:[%v]  Error:[%v]", key, string(buffer))
			return cmd.Err()
		}

		log.Debugf("UpdateExchangeRealPrice update done:[%#v]", exchangeRealPrice)
	} else {
		val := &ExchangeIDRealPrice{}
		val.ExchangeID = exchangeID
		val.RealPrice = fmt.Sprintf("%v", realPrice)
		if zeroPrice > 0 {
			val.ZeroPrice = fmt.Sprintf("%v", zeroPrice)
		}

		val.Amount = amount
		buffer, err := json.Marshal(val)
		if err != nil {
			log.Errorf(err, "Marshal error")
			return err
		}
		//写入redis 并不设置有效期
		cmd := RedisCli.Set(key, string(buffer), 0)
		if cmd.Err() != nil {
			log.Errorf(err, "set key:[%v] value:[%v]", key, string(buffer))
			return cmd.Err()
		}
	}
	return nil
}

//UpdateOrderblockOffset 更新实时处理的orderID 所属区块
func UpdateOrderblockOffset(blockID int64) error {
	log.Debugf("UpdateOrderblockOffset start update blockID: [%v]", blockID)

	//写入redis 并不设置有效期
	cmd := RedisCli.Set(OrderBlockIDOffset, blockID, 0)
	if cmd.Err() != nil {
		log.Errorf(cmd.Err(), "UpdateOrderblockOffset set key:[%v] value:[%v]",
			OrderBlockIDOffset, blockID)
		return cmd.Err()
	}

	return nil
}

//GetOrderblockOffset 查询当前已经处理的orderID 所属区块
func GetOrderblockOffset() int64 {

	val := redisGet(OrderBlockIDOffset)
	log.Debugf("GetOrderblockOffset start get blockID result: [%v]", val)

	return util.ConvertDBValueToInt64(val)
}

//----------------------------------BASE OPERATOR----------------------------------------------

//MaxErrCnt 最大重试次数
var MaxErrCnt = 10

func redisGet(key string) string {
	tryCnt := 0
	for {
		ret := RedisCli.Get(key)
		if nil == ret || nil != ret.Err() {
			tryCnt++

			if tryCnt > 10 {
				log.Errorf(ret.Err(), "redisGet(%v) error:[%#v]", key, ret)
				return ""
			}
			continue
		}
		return ret.Val()
	}
}

func redisExpire(key string, ttl time.Duration) bool {
	errCnt := 0
	for {
		ret := RedisCli.Expire(key, ttl)

		if nil == ret || nil != ret.Err() {
			errCnt++
		} else {
			return true
		}

		if errCnt > MaxErrCnt {
			log.Errorf(ret.Err(), "redis Expire [%v] [%v] try [%v] time failed", key, ttl, errCnt)
			break
		}
	}

	return false
}

func redisExists(key string) bool {
	errCnt := 0
	for {
		ret := RedisCli.Exists(key)

		if nil == ret || nil != ret.Err() {
			errCnt++
		} else {
			return ret.Val()
		}

		if errCnt > MaxErrCnt {
			log.Errorf(ret.Err(), "redis Exists [%v] try [%v] time failed", key, errCnt)
			break
		}
	}

	return false
}

//GetRedisKey 根据key获取redis value
func getRedisKey(key string) string {
	/*_, err := redisCli.Ping().Result()
	//fmt.Println(pong, err)
	if err != nil {
		return ""
	}*/
	val, err := RedisCli.Get(key).Result()
	if err != nil {
		return ""
	}
	//fmt.Printf("get keys:[%v], value:[%v]\n", key, val)
	return val

}

//getHRedisKey 根据key获取redis value
func getHRedisKey(key, subKey string) string {
	/*_, err := redisCli.Ping().Result()
	//fmt.Println(pong, err)
	if err != nil {
		return ""
	}*/
	val, err := RedisCli.HGet(key, subKey).Result()
	if err != nil {
		return ""
	}
	//fmt.Printf("get keys:[%v], value:[%v]\n", key, val)
	return val

}

//getAllHRedisKey 根据key获取redis value
func getAllHRedisKey(key string) map[string]string {
	/*_, err := redisCli.Ping().Result()
	//fmt.Println(pong, err)
	if err != nil {
		return ""
	}*/
	val, err := RedisCli.HGetAll(key).Result()
	if err != nil {
		return nil
	}
	//fmt.Printf("get keys:[%v], value:[%v]\n", key, val)
	return val

}
