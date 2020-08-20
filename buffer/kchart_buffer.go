package buffer

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/buffer2"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/module"
	"github.com/wlcy/tradehome-service/util"
	"strconv"
	"sync"
	"time"
)

var loadOnce sync.Once
var _kChartCache *kchartCache

const initialCaseSize = 4320

//StartKChartTimer ...
func StartKChartTimer() {
	GetKChartData()
	startMin1Timer()
	startMin5Timer()
	startMin15Timer()
	startMin30Timer()
	startH1Timer()
	startH4Timer()
	startDay1Timer()
	startDay5Timer()
	startWeekTimer()
	startMonthTimer()
	startAuditTimer()
}

type kchartCache struct {
	sync.RWMutex
	//exchange_id-->[granu-->[[]*entity.KGraphData]]
	kchartData map[string](map[string][]*entity.KChartData)
}

func prepareKChartCache() *kchartCache {
	loadOnce.Do(func() {
		_kChartCache = &kchartCache{}
		_kChartCache.initKChartCache(initialCaseSize)
		// go routineUpdateCache()
	})
	return _kChartCache
}

//GetKChartData ...
func GetKChartData() map[string]map[string][]*entity.KChartData {
	prepareKChartCache()
	return _kChartCache.kchartData
}

func (b *kchartCache) initKChartCache(initialCaseSize int) { //初始化不同的exchange_ID和不同的粒度
	if initialCaseSize <= 0 {
		_kChartCache.Lock()
		_kChartCache.kchartData = make(map[string]map[string][]*entity.KChartData)
		_kChartCache.Unlock()
		return
	}
	initMin1Data(initialCaseSize)
	initMin5Data(initialCaseSize)
	initMin15Data(initialCaseSize)
	initMin30Data(initialCaseSize)
	init1HData(initialCaseSize)
	init4HData(initialCaseSize)
	init1DData(initialCaseSize)
	init5DData()
	init1WData()
	init1MData()
}

// func routineUpdateCache() {
// 	startKChartTimer()
// }

func initMin1Data(initialCaseSize int) {
	log.Debugf("begin init kchart minute buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("1min", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func initMin5Data(initialCaseSize int) {
	log.Debugf("begin init kchart 5minute buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("5min", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func initMin15Data(initialCaseSize int) {
	log.Debugf("begin init kchart 15minute buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("15min", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func initMin30Data(initialCaseSize int) {
	log.Debugf("begin init kchart 30minute buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("30min", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func init1HData(initialCaseSize int) {
	log.Debugf("begin init kchart 1hour buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("1h", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func init4HData(initialCaseSize int) {
	log.Debugf("begin init kchart 4hour buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("4h", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func init1DData(initialCaseSize int) {
	log.Debugf("begin init kchart 1Day buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("1d", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func init5DData() { //5 years data
	log.Debugf("begin init kchart 5Day buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("5d", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func init1WData() { //5 years data
	log.Debugf("begin init kchart 1Week buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("1w", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func init1MData() { //5 years data
	log.Debugf("begin init kchart 1Month buffer, time now:%v", time.Now().Format(util.DATETIMEFORMAT))
	queryParam, err := constructParam("1m", initialCaseSize)
	if nil != err {
		log.Errorf(err, "construct parameter error")
		return
	}
	queryKChartDataAll(queryParam)
}

func constructParam(granu string, initialCaseSize int) (*entity.KChartQueryParam, error) {
	toTime := time.Now()
	var fromTime time.Time
	queryParam := &entity.KChartQueryParam{}
	queryParam.Granu = granu
	queryParam.TimeEnd = strconv.FormatInt(toTime.UTC().Unix(), 10)
	if "1min" == granu {
		fromTime = toTime.Add(-time.Minute * time.Duration(initialCaseSize))
	} else if "5min" == granu {
		fromTime = toTime.Add(-time.Minute * 5 * time.Duration(initialCaseSize))
	} else if "15min" == granu {
		fromTime = toTime.Add(-time.Minute * 15 * time.Duration(initialCaseSize))
	} else if "30min" == granu {
		fromTime = toTime.Add(-time.Minute * 30 * time.Duration(initialCaseSize))
	} else if "1h" == granu {
		fromTime = toTime.Add(-time.Hour * time.Duration(initialCaseSize))
	} else if "4h" == granu {
		fromTime = toTime.Add(-time.Hour * 4 * time.Duration(initialCaseSize))
	} else if "1d" == granu {
		fromTime = toTime.Add(-time.Hour * 24 * time.Duration(initialCaseSize))
	} else if "5d" == granu {
		fromTime = toTime.Add(-time.Hour * 24 * 5 * 365)
	} else if "1w" == granu {
		fromTime = toTime.Add(-time.Hour * 24 * 7 * 260)
	} else if "1m" == granu {
		fromTime = toTime.Add(-time.Hour * 24 * 31 * 60)
	} else {
		log.Error("Granularity parameter is not valid", nil)
		return nil, fmt.Errorf("granularity parameter is not valid")
	}
	queryParam.TimeStart = strconv.FormatInt(fromTime.UTC().Unix(), 10)
	return queryParam, nil
}

func queryKChartDataAll(queryParam *entity.KChartQueryParam) {
	marketPairs := buffer2.GetMarketPairBuffer().GetMarketPairs()
	if nil == marketPairs {
		log.Info("get market pairs returns nil")
		return
	}

	for exID := range marketPairs {
		queryParam.ExchangeID = strconv.FormatInt(exID, 10)
		log.Debugf("begin get grau[%v] kchart data from db for exchange id [%v]", queryParam.Granu, queryParam.ExchangeID)
		kchartResp, err := queryKChartData(queryParam)
		if nil != err || nil == kchartResp {
			log.Errorf(err, "get grau[%v] kchart data from db for exchange id [%v] error", queryParam.Granu, queryParam.ExchangeID)
			continue
		}
		log.Debugf("got data for  grau[%v] and exchange id [%v], size is: %v", queryParam.Granu, queryParam.ExchangeID, len(kchartResp.Data))
		if nil == _kChartCache.kchartData {
			_kChartCache.Lock()
			_kChartCache.kchartData = make(map[string]map[string][]*entity.KChartData)
			_kChartCache.Unlock()
		}

		_kChartCache.Lock()
		granuKChartDataMap, ok := _kChartCache.kchartData[queryParam.ExchangeID]
		if ok { //此exchange ID
			log.Debugf("cache data exist exchange id [%v]:\n", queryParam.ExchangeID)
			if nil == granuKChartDataMap { //不存在此exchange ID则新建并加入map
				granuKChartDataMap := make(map[string][]*entity.KChartData)
				granuKChartDataMap[queryParam.Granu] = kchartResp.Data
				_kChartCache.kchartData[queryParam.ExchangeID] = granuKChartDataMap
			} else { //存在直接更新或者加入此granu
				granuKChartDataMap[queryParam.Granu] = kchartResp.Data
			}
		} else {
			granuKChartDataMap := make(map[string][]*entity.KChartData)
			granuKChartDataMap[queryParam.Granu] = kchartResp.Data
			_kChartCache.kchartData[queryParam.ExchangeID] = granuKChartDataMap
		}
		log.Debugf("finish initate cache data for exchange id [%v] and granu[%v], data size[%v] \n", queryParam.ExchangeID, queryParam.Granu, len(_kChartCache.kchartData[queryParam.ExchangeID][queryParam.Granu]))
		_kChartCache.Unlock()
	}
}

func queryKChartData(param *entity.KChartQueryParam) (*entity.KChartResp, error) {
	if nil == param || "" == param.ExchangeID || "" == param.TimeStart || "" == param.TimeEnd || "" == param.Granu {
		log.Error("Parameter(s) not valid for QueryKChartData", nil)
		return nil, fmt.Errorf("parameter(s) not valid for QueryKChartData")
	}
	granu := param.Granu
	exchangeID, err := strconv.ParseInt(param.ExchangeID, 10, 64)
	if nil != err {
		log.Errorf(err, "Parameter of exchange_id is not right ")
		return nil, fmt.Errorf("Parameter of exchange id is not right")
	}
	marketPair := buffer2.GetMarketPairBuffer().GetMarketPairs()[exchangeID]
	firstDecimal := marketPair.FirstPrecision
	secondDecimal := marketPair.SecondPrecision
	log.Debugf("begin query kchart data, granu is: %v, firstDecimal is: %v, secondDecimal is: %v", granu, firstDecimal, secondDecimal)
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

//StartMinTimer ...
func startMin1Timer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("1min", now, false)
			targetTime := tempTime.Add(-time.Minute * 2)
			log.Infof("start gather 1 minutes data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.Handle1minKChartData(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("1min", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(10 * time.Second))
		}
	}()
}

func startMin5Timer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("5min", now, false)
			targetTime := tempTime.Add(-time.Minute * 2 * 5)
			log.Infof("start gather 5 minutes data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartMin5(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("5min", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(10 * time.Second))
		}
	}()
}

func startMin15Timer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("15min", now, false)
			targetTime := tempTime.Add(-time.Minute * 2 * 15)
			log.Infof("start gather 15 minutes data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartMin15(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("15min", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(time.Minute))
		}
	}()
}

func startMin30Timer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("30min", now, false)
			targetTime := tempTime.Add(-time.Minute * 2 * 30)
			log.Infof("start gather 30 minutes data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartMin30(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("30min", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(time.Minute))
		}
	}()
}

func startH1Timer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("1h", now, false)
			targetTime := tempTime.Add(-time.Hour * 2)
			log.Infof("start gather 1 hour data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartH1(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("1h", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(time.Minute))
		}
	}()
}

func startH4Timer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("4h", now, false)
			targetTime := tempTime.Add(-time.Hour * 8)
			log.Infof("start gather 4 hours data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartH4(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("4h", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(time.Minute))
		}
	}()
}

//StartDay1Timer ...
func startDay1Timer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("1d", now, false)
			targetTime := tempTime.Add(-time.Hour * 24 * 2)
			log.Infof("start gather 1 day data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartD1(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("1d", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(time.Hour))
		}
	}()
}

//StartDay5Timer ...
func startDay5Timer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("5d", now, false)
			targetTime := tempTime.Add(-time.Hour * 24 * 2 * 5)
			log.Infof("start gather 5 days data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartD5(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("5d", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(24 * time.Hour))
		}
	}()
}

//StartWeekTimer ...
func startWeekTimer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("1w", now, false)
			targetTime := tempTime.Add(-time.Hour * 24 * 2 * 7)
			log.Infof("start gather 1 week data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartW(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("1w", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(24 * time.Hour))
		}
	}()
}

//StartMonthTimer ...
func startMonthTimer() {
	go func() {
		for {
			now := time.Now().UTC()
			_, tempTime, _ := util.GetBenchmarkTimeAndSpan("1m", now, false)
			targetTime := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
			log.Infof("start gather 1 month data for kgraph,\n now: %v\n tempTime:%v\n targetTime: %v\n",
				now.Format(util.DATETIMEFORMAT), tempTime.Format(util.DATETIMEFORMAT), targetTime.Format(util.DATETIMEFORMAT))
			exIDKChartData, err := module.BuildKChartM(targetTime)
			if nil == err && nil != exIDKChartData {
				if nil != _kChartCache {
					_kChartCache.Lock()
					updateKChartData("1m", exIDKChartData)
					_kChartCache.Unlock()
				}
			}
			time.Sleep(time.Duration(24 * time.Hour))
		}
	}()
}

func startAuditTimer() {
	go func() {
		for {
			//audit 程序启动的时候不启动，一段时间后启动
			time.Sleep(time.Duration(24 * time.Hour))
			//TODO
			cacheData := GetKChartData()
			if nil == cacheData {
				continue
			}
			for exID, granuKChartData := range cacheData {
				if nil == granuKChartData {
					continue
				}
				for granu, kchartDataList := range granuKChartData {
					if nil == kchartDataList {
						continue
					}
					log.Debugf("auditing exchange id-[%v]-granu[%v] data, size is:%v \n", exID, granu, len(kchartDataList))
					if len(kchartDataList) > initialCaseSize {
						audited := make([]*entity.KChartData, 0)
						for i := len(kchartDataList) - initialCaseSize; i < len(kchartDataList); i++ {
							audited = append(audited, kchartDataList[i])
						}
						granuKChartData[granu] = audited
					}
				}
			}
		}
	}()
}

func updateKChartData(granu string, exIDKChartData map[string][]*entity.KChartData) {
	log.Debugf("Begin to update cache for granu[%v]\n", granu)
	if nil == exIDKChartData {
		return
	}

	for exID, kchartDataList := range exIDKChartData {
		granuKChartData, ok := _kChartCache.kchartData[exID]
		if ok { //有这个exchangeID
			kchartDataListCache, ok := granuKChartData[granu]
			if ok { //有这个granu,则查找cache中最大的time，如果大于更新中最小的，则删除，使用更新的替换
				if len(kchartDataList) == 0 {
					continue
				}
				kchartData := kchartDataList[0]
				time1, err := strconv.ParseInt(kchartData.Time, 10, 64)
				if nil != err {
					log.Errorf(err, "parse time[%v] form update data error", kchartData.Time)
					break
				}

				delCounter := 0
				log.Debugf("Cache data for granu[%v] exchangeID[%v] size is: %v \n", granu, exID, len(kchartDataListCache))
				for i := len(kchartDataListCache) - 1; i >= 0; i-- {
					kcharDataCache := kchartDataListCache[i]
					timeCache, err1 := strconv.ParseInt(kcharDataCache.Time, 10, 64)
					if nil != err1 {
						log.Errorf(err1, "parse time[%v] form cache data error", kcharDataCache.Time)
						break
					}
					if time1 <= timeCache { //删掉缓存中存在的已经更新的数据
						delCounter++
						kchartDataListCache = kchartDataListCache[:i]
					} else { //把新查到的数据更新到缓存中
						log.Debugf("deleted cache data size is [%v], now cache data size is [%v]\n", delCounter, len(kchartDataListCache))
						log.Debugf("update data for granu[%v] exchangeID[%v] size is: %v \n", granu, exID, len(kchartDataList))
						for index := 0; index < len(kchartDataList); index++ {
							kchartDataListCache = append(kchartDataListCache, kchartDataList[index])
						}
						log.Debugf("after update cache data for granu[%v] exchangeID[%v] size is: %v \n", granu, exID, len(kchartDataListCache))
						granuKChartData[granu] = kchartDataListCache
						_kChartCache.kchartData[exID] = granuKChartData
						break
					}
				}

			} else { //没有granu， 则直接把此granu的数据加入
				if nil == granuKChartData {
					granuKChartData = make(map[string][]*entity.KChartData)
					granuKChartData[granu] = kchartDataList
					_kChartCache.kchartData[exID] = granuKChartData
				} else {
					granuKChartData[granu] = kchartDataList
					_kChartCache.kchartData[exID] = granuKChartData
				}
			}
		} else { //没有此exchange id，则把此exchange 对应的granu数据加入cache
			granuKChartData := make(map[string][]*entity.KChartData)
			granuKChartData[granu] = kchartDataList
			if _kChartCache.kchartData == nil {
				_kChartCache.kchartData = make(map[string]map[string][]*entity.KChartData)
			}
			_kChartCache.kchartData[exID] = granuKChartData
		}
	}
}
