package buffer

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/redis"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/module"
	"github.com/wlcy/tradehome-service/util"
)

var _recordBuffer *recordBuffer
var onceRecordBuffer sync.Once

type recordBuffer struct {
	mtxMarket       sync.RWMutex
	mtxLOrder       sync.RWMutex
	mtxUOrder       sync.RWMutex
	marketPairList  *entity.MarketPairListInfo
	latestOrderList *entity.LatestOrderListInfo
	userOrderList   *entity.UserOrderListInfoDetail
}

func GetRecordBuffer() *recordBuffer {
	onceRecordBuffer.Do(func() {
		_recordBuffer = &recordBuffer{}
		_recordBuffer.loadMarketPairListBuffer()
		go marketPairListBufferLoader()
	})
	return _recordBuffer
}

//对外获取数据函数
func (b *recordBuffer) GetMarketPairList() *entity.MarketPairListInfo {
	return b.marketPairList
}

// 定时重载数据
func marketPairListBufferLoader() {
	for {
		_recordBuffer.loadMarketPairListBuffer()
		time.Sleep(3 * time.Second)
	}
}

func (b *recordBuffer) loadMarketPairListBuffer() {
	//缓存所有
	retList := marketPairList(&entity.PageInfo{0, 0})
	if nil != retList && len(retList.Rows) > 0 {
		b.mtxMarket.Lock()
		b.marketPairList = retList
		b.mtxMarket.Unlock()
	}
}

//获取交易对
func marketPairList(info *entity.PageInfo) *entity.MarketPairListInfo {
	marketPairs, totalCount, err := module.MarketPairList(info)
	if err != nil {
		log.Errorf(err, "marketPairList buffer load error")
		return nil
	}
	if marketPairs == nil {
		return &entity.MarketPairListInfo{Rows: nil, Total: 0}
	}
	rows := make([]*entity.MarketPairRowInfo, 0)
	rowsIdx := make([]*entity.MarketPairRowInfo, 0)
	for _, marketPair := range marketPairs {
		_, price, gain := getPriceAndGinInfo(marketPair.ID)
		priceInf := getPriceInf(marketPair.ID)
		row := &entity.MarketPairRowInfo{
			ID:              marketPair.ID,
			FirstShortName:  marketPair.FisrtShortName,
			FirstTokenName:  marketPair.FirstTokenName,
			FirstPrecision:  marketPair.FirstPrecision,
			FirstTokenAddr:  marketPair.FirstTokenAddr,
			SecondShortName: marketPair.SecondShortName,
			SecondTokenName: marketPair.SecondTokenName,
			SecondPrecision: marketPair.SecondPrecision,
			SecondTokenAddr: marketPair.SecondTokenAddr,
			//Volume:          float64(volume),
			//Price:    price,
			Price:      float64(price) / math.Pow10(int(marketPair.SecondPrecision)),
			Gain:       fmt.Sprintf("%f", gain),
			Unit:       marketPair.Unit,
			PairType:   marketPair.PairType,
			DefaultIdx: marketPair.DefaultIdx,
		}
		if priceInf != nil {
			row.Volume = priceInf.Amount24h
			row.HighestPrice24h = priceInf.HighestPrice24h
			row.LowestPrice24h = priceInf.LowestPrice24h
			row.Volume24h = priceInf.Volume24h
		}
		if row.PairType == 1 && row.SecondTokenName == "TRX" {
			row.SecondTokenAddr = "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb"
		}
		logoUrl := module.GetFirstTokenLogoUrl(marketPair.FirstTokenAddr)
		row.LogoUrl = logoUrl
		if row.DefaultIdx > 0 {
			rowsIdx = append(rowsIdx, row)
		} else {
			rows = append(rows, row)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Volume24h > rows[j].Volume24h })
	//sort.Slice(rowsIdx, func(i, j int) bool { return rowsIdx[i].DefaultIdx < rowsIdx[j].DefaultIdx })
	log.Debug("marketPairList success")
	rowsIdx = append(rowsIdx, rows...)
	return &entity.MarketPairListInfo{Rows: rowsIdx, Total: totalCount}
}

//获取24h内的最高价和最低价
func getPriceInf(pairID int64) *entity.PriceInf {
	priceInfs := GetPirceInfBuffer().GetPriceInfs()
	if priceInfs == nil {
		return nil
	}
	return priceInfs[pairID]
}

//获取涨幅 = 最新成交价/0点价格 - 1
func getMarketPairGain(realPrice float64, zeroPrice float64) float64 {
	if zeroPrice == 0 {
		zeroPrice = 1
	}
	return (realPrice - zeroPrice) / zeroPrice
}

//获取 当前交易对的成交量、最新价格、涨幅
func getPriceAndGinInfo(pairID int64) (int64, float64, float64) {
	priceInfo := redis.GetExchangeRealPrice(pairID)
	log.Debugf("getPriceAndGinInfo:[%#v]", priceInfo)
	if priceInfo != nil {
		// realPrice, _ := strconv.ParseFloat(priceInfo.RealPrice, 64)
		// zeroPrice, _ := strconv.ParseFloat(priceInfo.ZeroPrice, 64)
		realPrice := util.ConvertDBValueToFloat64(priceInfo.RealPrice)
		zeroPrice := util.ConvertDBValueToFloat64(priceInfo.ZeroPrice)
		return priceInfo.Amount, realPrice, getMarketPairGain(realPrice, zeroPrice)
	}
	return 0, 0, 0
}
