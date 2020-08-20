package service

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/buffer"
	"github.com/wlcy/tradehome-service/buffer2"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/module"
	"github.com/wlcy/tradehome-service/util"
	"strconv"
	"time"
)

//QueryKChartData ...
func QueryKChartData(param *entity.KChartQueryParam) (*entity.KChartResp, error) {
	if nil == param || "" == param.ExchangeID || "" == param.TimeStart || "" == param.TimeEnd || "" == param.Granu {
		log.Error("Parameter(s) not valid for QueryKChartData", nil)
		return nil, fmt.Errorf("parameter(s) not valid for QueryKChartData")
	}
	granu := param.Granu
	exchangeID, err := strconv.ParseInt(param.ExchangeID, 10, 64)
	if nil != err {
		log.Errorf(err, "Parameter of exchange_id is not right")
		return nil, fmt.Errorf("parameter of exchange id is not right")
	}
	if nil == buffer2.GetMarketPairBuffer() || nil == buffer2.GetMarketPairBuffer().GetMarketPairs() {
		log.Errorf(err, "Get market pair from buffer returns nil")
		return nil, fmt.Errorf("get market pair from buffer returns nil")
	}

	marketPair, ok := buffer2.GetMarketPairBuffer().GetMarketPairs()[exchangeID]
	if !ok || nil == marketPair {
		log.Error("Market pair doesn't have pair ID", nil)
		return nil, fmt.Errorf("market pair doesn't have request pair ID[%v]", exchangeID)
	}

	firstDecimal := marketPair.FirstPrecision
	secondDecimal := marketPair.SecondPrecision
	log.Debugf("begin query kchart data, exchange_id is: %v,  granu is: %v, firstDecimal is: %v, secondDecimal is: %v", exchangeID, granu, firstDecimal, secondDecimal)
	if "1min" == granu {
		return module.GetKChartData(module.TableMin1, param, firstDecimal, secondDecimal)
	} else if "5min" == granu {
		return module.GetKChartData(module.TableMin5, param, firstDecimal, secondDecimal)
	} else if "15min" == granu {
		return module.GetKChartData(module.TableMin15, param, firstDecimal, secondDecimal)
	} else if "30min" == granu {
		return module.GetKChartData(module.TableMin30, param, firstDecimal, secondDecimal)
	} else if "1h" == granu {
		return module.GetKChartData(module.TableH1, param, firstDecimal, secondDecimal)
	} else if "4h" == granu {
		return module.GetKChartData(module.TableH4, param, firstDecimal, secondDecimal)
	} else if "1d" == granu {
		return module.GetKChartData(module.TableD1, param, firstDecimal, secondDecimal)
	} else if "5d" == granu {
		return module.GetKChartData(module.TableD5, param, firstDecimal, secondDecimal)
	} else if "1w" == granu {
		return module.GetKChartData(module.TableW, param, firstDecimal, secondDecimal)
	} else if "1m" == granu {
		return module.GetKChartData(module.TableM, param, firstDecimal, secondDecimal)
	} else {
		log.Error("Granularity parameter is not valid for QueryKChartData", nil)
		return nil, fmt.Errorf("granularity parameter is not valid")
	}
}

//GetKChartData ...
func GetKChartData(param *entity.KChartQueryParam) (*entity.KChartResp, error) {
	timeStart, err := strconv.ParseInt(param.TimeStart, 10, 64)
	if nil != err {
		log.Errorf(err, "Parameter of time_start is not right")
		return nil, fmt.Errorf("parameter of time_start is not right")
	}
	timeEnd, err := strconv.ParseInt(param.TimeEnd, 10, 64)
	if nil != err {
		log.Errorf(err, "Parameter of time_end is not right")
		return nil, fmt.Errorf("parameter of time_end is not right")
	}
	fTime := time.Unix(timeStart, 0).UTC()
	tTime := time.Unix(timeEnd, 0).UTC()

	_, fromTime, err := util.GetBenchmarkTimeAndSpan(param.Granu, fTime, false)
	if nil != err {
		log.Errorf(err, "GetBenchmarkTimeAndSpan for fromTime error")
		return nil, fmt.Errorf("GetBenchmarkTimeAndSpan for fromTime error")
	}

	_, toTime, err := util.GetBenchmarkTimeAndSpan(param.Granu, tTime, false)
	if nil != err {
		log.Errorf(err, "GetBenchmarkTimeAndSpan for toTime error")
		return nil, fmt.Errorf("GetBenchmarkTimeAndSpan for toTime error")
	}
	allKChartData := buffer.GetKChartData()
	granuKchartData, ok := allKChartData[param.ExchangeID]
	if !ok {
		log.Infof("Data for ExchangeID[%v] is not in cache\n", param.ExchangeID)
		return QueryKChartData(param) //query from DB
	}
	kchartData, ok := granuKchartData[param.Granu]
	if !ok {
		log.Infof("Data for ExchangeID[%v] Granu[%v] is not in cache\n", param.ExchangeID, param.Granu)
		return QueryKChartData(param) //query from DB
	}

	if nil == kchartData || 0 == len(kchartData) {
		log.Infof("Data for ExchangeID[%v] Granu[%v] in cache is empty\n", param.ExchangeID, param.Granu)
		return QueryKChartData(param) //query from DB
	}

	oldestDataPieceUnixTime, err := strconv.ParseInt(kchartData[0].Time, 10, 64)
	if nil != err {
		log.Infof("Parse kchart time[%v] in cache data error:\n", kchartData[0].Time)
		return QueryKChartData(param) //query from DB
	}

	//1. totime< firstDataTime, 全从DB中查
	//2. fromTime >= firstDataTime, 全从Cache中查
	//3. fromTime < firstDataTime && toTime >= fromDataTime, 则fromTime～firstDataTime-granu从DB中查，firstDataTime～toTime从cache中查
	log.Debugf("oldest data in cache is: [%v]\n", oldestDataPieceUnixTime)
	if toTime.UTC().Unix() < oldestDataPieceUnixTime {
		log.Debug("time end less than cache data, get kchart data from DB")
		return QueryKChartData(param) //query from DB
	}

	if fromTime.UTC().Unix() >= oldestDataPieceUnixTime {
		log.Debug("time start in cache data, get kchart data from cache")
		return getDataFromCache(fromTime, toTime, kchartData, param), nil
	}
	log.Debug("will get kchart data partly from cache, partly from DB")
	cacheResp := getDataFromCache(time.Unix(oldestDataPieceUnixTime, 0).UTC(), toTime, kchartData, param)
	param.TimeEnd = kchartData[0].Time
	dbResp, err := QueryKChartData(param)
	if nil != err { //查询数据库出错则只返回cache中的数据
		log.Errorf(err, "QueryKChartData from DB error")
		return cacheResp, nil
	}
	if nil == dbResp || nil == dbResp.Data { //数据库返回数据为空
		log.Error("QueryKChartData from DB returns no data", nil)
		return cacheResp, nil
	}

	//合并数据
	dbResp.TimeEnd = cacheResp.TimeEnd
	dbData := dbResp.Data
	cacheData := cacheResp.Data
	log.Debugf("data from DB size[%v], data from cache size[%v]\n", len(dbData), len(cacheData))
	for i := 0; i < len(cacheData); i++ {
		dbData = append(dbData, cacheData[i])
	}
	dbResp.Data = dbData
	return dbResp, nil
}

func getDataFromCache(fromTime time.Time, toTime time.Time, cacheData []*entity.KChartData, param *entity.KChartQueryParam) *entity.KChartResp {
	kcharList := make([]*entity.KChartData, 0)
	for i := 0; i < len(cacheData); i++ {
		piece := cacheData[i]
		pieceUnixTime, _ := strconv.ParseInt(piece.Time, 10, 64)
		if pieceUnixTime < fromTime.UTC().Unix() || pieceUnixTime > toTime.UTC().Unix() {
			continue
		}
		kcharList = append(kcharList, piece)
	}

	kchartResp := &entity.KChartResp{}
	kchartResp.ExchangeID = param.ExchangeID
	kchartResp.Granu = param.Granu
	kchartResp.TimeEnd = strconv.FormatInt(toTime.UTC().Unix(), 10)
	kchartResp.TimeStart = strconv.FormatInt(fromTime.UTC().Unix(), 10)
	kchartResp.Data = kcharList
	return kchartResp
}
