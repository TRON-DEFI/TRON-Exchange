package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/common/redis"
	"github.com/wlcy/tradehome-service/util"
	"time"
)

//GetExchangeZeroPrice 获取所有交易对0点价格并更新redis
//获取规则：定时器0点获取交易表每隔交易对最近一笔交易价格，更新对应rediskey中
func GetExchangeZeroPrice(exchangeID int64) (int64, error) {
	now := time.Now().UTC()
	d, _ := time.ParseDuration("-24h")
	d24h := now.Add(d).UTC().UnixNano() / 1000000
	strSQL := fmt.Sprintf(`select price 
	from trxmarket.market_transaction 
	where pair_id=%v and create_time>=%v 
	order by create_time desc limit 1`, exchangeID, d24h)
	log.Info(strSQL)
	zeroPrice := int64(0)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf(err, "GetExchangeZeroPrice error")
		return zeroPrice, errors.New("internal error")
	}
	if dataPtr == nil {
		log.Error("GetExchangeZeroPrice dataPtr is nil", nil)
		return zeroPrice, errors.New("internal error")
	}

	//填充数据
	for dataPtr.NextT() {
		zeroPrice = util.ConvertDBValueToInt64(dataPtr.GetField("price"))
	}
	return zeroPrice, nil
}

//InsertZeroPrice 更新交易对id零点价格，数据库和redis同时
func InsertZeroPrice(exchangeID, zeroPrice int64) error {
	err := updateExchangeZeroPrice(zeroPrice, exchangeID)
	if err != nil {
		log.Errorf(err, "update zero price in redis for exchangeID[%v] error", exchangeID)
	}
	strSQL := fmt.Sprintf(`update trxmarket.market_price set zero_price = %v where pair_id=%v `, zeroPrice, exchangeID)
	log.Info(strSQL)
	_, _, err = mysql.ExecuteSQLCommand(strSQL)
	if err != nil {
		log.Errorf(err, "InsertZeroPrice error,  sql:%s", strSQL)
		return err
	}
	log.Debugf("InsertZeroPrice success, id: [%v],zeroPrice:[%v]", exchangeID, zeroPrice)
	return nil
}

//updateExchangeZeroPrice 更新0点价格等信息
func updateExchangeZeroPrice(zeroPrice int64, exchangeID int64) error {
	log.Debugf("updateExchangeZeroPrice start update exchangeID: [%v],zeroPrice: [%v],", exchangeID, zeroPrice)
	key := fmt.Sprintf(redis.ExchangeIDRealPriceRedisKey, exchangeID)

	exchangeRealPrice := redis.GetExchangeRealPrice(exchangeID)
	if exchangeRealPrice == nil {
		log.Error("get updateExchangeZeroPrice in redis nil", nil)
		return errors.New("get updateExchangeZeroPrice in redis nil")
	}

	exchangeRealPrice.ZeroPrice = fmt.Sprintf("%v", zeroPrice)

	buffer, err := json.Marshal(exchangeRealPrice)
	if err != nil {
		log.Errorf(err, "Marshal error")
		return err
	}

	//写入redis 并不设置有效期
	err = redis.RedisCli.Set(key, string(buffer), 0).Err()
	if err != nil {
		log.Errorf(err, "set key:[%v] value:[%v] error",
			key, string(buffer))
		return err
	}

	log.Debugf("updateExchangeZeroPrice update done:[%#v]", exchangeRealPrice)
	return nil
}
