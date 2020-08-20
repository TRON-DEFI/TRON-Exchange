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

var _exchangeBuffer *exchangeBuffer
var onceExchangeBuffer sync.Once

func GetExchangeBuffer() *exchangeBuffer {
	return getExchangeBuffer()
}

func getExchangeBuffer() *exchangeBuffer {
	onceExchangeBuffer.Do(func() {
		_exchangeBuffer = &exchangeBuffer{}
		_exchangeBuffer.load()
		go exchangeBufferLoader()
	})
	return _exchangeBuffer
}

func exchangeBufferLoader() {
	for {
		_exchangeBuffer.load()
		time.Sleep(60 * time.Second)

	}
}

type exchangeBuffer struct {
	sync.RWMutex
	//maxBlockID int64
	exchanges         map[string]int64               //交易地址对-交易对ID
	exchangeIDs       map[int64]int64                //交易对ID-交易对ID
	exchangeIDAddress map[int64]string               //交易对ID-first token address
	exchangeInfo      map[int64]*entity.ExchangeInfo //交易对id-精度信息
}

func (w *exchangeBuffer) GetAddrByExchange(exchangeID int64) (firstTokenAddress string, ok bool) {
	w.RLock()
	firstTokenAddress, ok = w.exchangeIDAddress[exchangeID]
	w.RUnlock()
	return
}

func (w *exchangeBuffer) GetExchangeIDByAddr(addr string) (exchangeID int64, ok bool) {
	w.RLock()
	exchangeID, ok = w.exchanges[addr]
	w.RUnlock()
	return
}

func (w *exchangeBuffer) GetExchangeDicemal(id int64) (exchangeInfo *entity.ExchangeInfo, ok bool) {
	w.RLock()
	exchangeInfo, ok = w.exchangeInfo[id]
	w.RUnlock()
	return
}

func (w *exchangeBuffer) GetExchanges() (exchanges map[string]int64) {
	w.RLock()
	exchanges = w.exchanges
	w.RUnlock()
	return
}

func (w *exchangeBuffer) GetExchangeIDs() (exchanges map[int64]int64) {
	w.RLock()
	exchanges = w.exchangeIDs
	w.RUnlock()
	return
}
func (w *exchangeBuffer) load() {
	strSQL := fmt.Sprintf(`
	select market.first_token_addr, market.second_token_addr,market.id,market.first_token_precision,market.second_token_precision
			 from trxmarket.market_pair market
			 where 1=1`)

	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf(err, "load exchange_pair error")
		return
	}
	if dataPtr == nil {
		log.Error("load exchange_pair dataPtr is nil", nil)
		return
	}
	exchanges := make(map[string]int64, 0)
	exchangeIDs := make(map[int64]int64, 0)
	exchangeInfos := make(map[int64]*entity.ExchangeInfo, 0)
	exchangeAddr := make(map[int64]string, 0)
	//填充数据
	for dataPtr.NextT() {
		exh := &entity.ExchangeInfo{}
		addressA := dataPtr.GetField("first_token_addr")
		addressB := dataPtr.GetField("second_token_addr")
		exchangeID := util.ConvertDBValueToInt64(dataPtr.GetField("id"))
		exh.ExchangeID = exchangeID
		exh.FirstPrecision = util.ConvertDBValueToInt64(dataPtr.GetField("first_token_precision"))
		exh.SecondPrecision = util.ConvertDBValueToInt64(dataPtr.GetField("second_token_precision"))
		addrPair := fmt.Sprintf("%v-%v", addressA, addressB)
		exchanges[addrPair] = exchangeID
		exchangeIDs[exchangeID] = exchangeID
		exchangeAddr[exchangeID] = addressA
		exchangeInfos[exchangeID] = exh
	}

	w.Lock()
	w.exchanges = exchanges
	w.exchangeIDs = exchangeIDs
	w.exchangeInfo = exchangeInfos
	w.exchangeIDAddress = exchangeAddr
	log.Debugf("set exchange_pair buffer data done.")
	w.Unlock()
}
