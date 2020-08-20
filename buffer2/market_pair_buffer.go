package buffer2

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/util"
	"sync"
	"time"
)

var _marketPairBuffer *marketPairBuffer
var onceMarketPairBuffer sync.Once

type marketPairBuffer struct {
	sync.RWMutex
	marketPairInfo map[string]*entity.MarketPairInfo
	marketPairs    map[int64]*entity.MarketPair
}

func GetMarketPairBuffer() *marketPairBuffer {
	onceMarketPairBuffer.Do(func() {
		_marketPairBuffer = &marketPairBuffer{}
		_marketPairBuffer.loadMarketPairBuffer()
		go marketPairBufferLoader()

	})
	return _marketPairBuffer
}

//对外获取数据函数
func (b *marketPairBuffer) GetMarketPairInfo() map[string]*entity.MarketPairInfo {
	return b.marketPairInfo
}

//对外获取数据函数,获取全信息
func (b *marketPairBuffer) GetMarketPairs() map[int64]*entity.MarketPair {
	return b.marketPairs
}

// 定时重载数据
func marketPairBufferLoader() {
	for {
		_marketPairBuffer.loadMarketPairBuffer()
		time.Sleep(30 * time.Second)
	}
}

//加载数据
func (b *marketPairBuffer) loadMarketPairBuffer() {
	b.Lock()
	b.marketPairInfo, b.marketPairs = marketPairInfo()
	b.Unlock()
}

func marketPairInfo() (map[string]*entity.MarketPairInfo, map[int64]*entity.MarketPair) {
	strSQL := fmt.Sprintf(`select * from market_pair where is_valid = 1`)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)

	if err != nil {
		log.Errorf(err, "MarketPairInfo error")
		return nil, nil
	}
	if dataPtr == nil {
		log.Error("MarketPairInfo dataPtr is nil", nil)
		return nil, nil
	}

	return data2MarketPairInfo(dataPtr)
}

func data2MarketPairInfo(dataPtr *mysql.TronDBRows) (map[string]*entity.MarketPairInfo, map[int64]*entity.MarketPair) {
	data := make(map[string]*entity.MarketPairInfo, 0)
	dataFull := make(map[int64]*entity.MarketPair, 0)
	for dataPtr.NextT() {
		mpi := &entity.MarketPairInfo{}
		mpi.FirstTokenName = dataPtr.GetField("first_short_name")
		mpi.SecondTokenName = dataPtr.GetField("second_short_name")
		key := dataPtr.GetField("first_token_addr") + dataPtr.GetField("second_token_addr")
		data[key] = mpi

		mpf := &entity.MarketPair{}
		mpf.ID = util.ConvertStringToInt64(dataPtr.GetField("id"), 0)
		mpf.FirstTokenName = dataPtr.GetField("first_token_name")
		mpf.FirstTokenAddr = dataPtr.GetField("first_token_addr")
		mpf.SecondTokenName = dataPtr.GetField("second_token_name")
		mpf.SecondTokenAddr = dataPtr.GetField("second_token_addr")
		mpf.Price = util.ConvertDBValueToFloat64(dataPtr.GetField("price"))
		mpf.Unit = dataPtr.GetField("unit")
		mpf.FirstPrecision = util.ConvertStringToInt64(dataPtr.GetField("first_token_precision"), 6)
		mpf.SecondPrecision = util.ConvertStringToInt64(dataPtr.GetField("second_token_precision"), 6)
		mpf.FisrtShortName = dataPtr.GetField("first_short_name")
		mpf.SecondShortName = dataPtr.GetField("second_short_name")
		mpf.CreatedAt = dataPtr.GetField("create_time")
		mpf.UpdateAt = dataPtr.GetField("update_time")
		mpf.PairType = util.ConvertStringToInt64(dataPtr.GetField("pair_type"), 0)
		log.Debugf("add market pair id[%v]\n", mpf.ID)
		dataFull[mpf.ID] = mpf
	}
	log.Debugf("init or update finished for market pair")
	return data, dataFull

}
