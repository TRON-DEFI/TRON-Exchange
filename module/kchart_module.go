package module

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/util"
	"math"
	"strconv"
	"time"
)

var TableMin1 = "exchange_kchart_min1"
var TableMin5 = "exchange_kchart_min5"
var TableMin15 = "exchange_kchart_min15"
var TableMin30 = "exchange_kchart_min30"
var TableH1 = "exchange_kchart_h1"
var TableH4 = "exchange_kchart_h4"
var TableD1 = "exchange_kchart_d1"
var TableD5 = "exchange_kchart_d5"
var TableW = "exchange_kchart_w"
var TableM = "exchange_kchart_m"

//PRECISION price's precision
const PRICE_PRECISION = 6

//PRICE_FACTOR price's factor for store and get from db
const PRICE_FACTOR = 1000000

//BuildKChartMin5 ...
func BuildKChartMin5(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartData(TableMin1, TableMin5, time.Minute*5, fromTime)
}

//BuildKChartMin15 ...
func BuildKChartMin15(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartData(TableMin5, TableMin15, time.Minute*15, fromTime)
}

//BuildKChartMin30 ...
func BuildKChartMin30(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartData(TableMin15, TableMin30, time.Minute*30, fromTime)
}

//BuildKChartH1 ...
func BuildKChartH1(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartData(TableMin30, TableH1, time.Hour, fromTime)
}

//BuildKChartH4 ...
func BuildKChartH4(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartData(TableH1, TableH4, time.Hour*4, fromTime)
}

//BuildKChartD1 ...
func BuildKChartD1(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartData(TableH4, TableD1, time.Hour*24, fromTime)
}

//BuildKChartD5 ...
func BuildKChartD5(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartData(TableD1, TableD5, time.Hour*24*5, fromTime)
}

//BuildKChartW ...
func BuildKChartW(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartData(TableD1, TableW, time.Hour*24*7, fromTime)
}

//BuildKChartM ...
func BuildKChartM(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	return buildKChartDataMonth(TableD1, TableM, fromTime)
}

func buildKChartDataMonth(fromKChartTable string, toKChartTable string, fromTime time.Time) (map[string][]*entity.KChartData, error) {
	log.Debugf("In buildKChartData, from time is: %v", fromTime.UTC().Unix())
	strSQL := fmt.Sprintf(`select id, first_token_precision, second_token_precision from market_pair where 1 = 1`) //only query the valid pair??
	exchangeIDPtr, err := mysql.QueryTableData(strSQL)
	if nil != err {
		log.Errorf(err, "Query exchange id error")
		return nil, fmt.Errorf("query exchange id error")
	}
	if nil == exchangeIDPtr || 0 == exchangeIDPtr.ResNum() {
		log.Errorf(err, "exchange id list is nil or empty")
		return nil, fmt.Errorf("query exchange id error")
	}

	exIDKChartData := make(map[string][]*entity.KChartData)
	for exchangeIDPtr.NextT() {
		//traverse the exchange_id and get the corresponding data from exchange transaction table
		//and insert to the kchart table
		kchartDataList := make([]*entity.KChartData, 0)
		exchangeID := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("id"))
		firstDecimal := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("first_token_precision"))
		secondDecimal := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("second_token_precision"))
		fromDataPtr, err := queryKChartData(fromKChartTable, fromTime, time.Now(), exchangeID)
		var open, high, low, close, volume int64
		startTimeUnix := fromTime.UTC().Unix()
		toTime := time.Date(fromTime.Year(), fromTime.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		toTimeUnix := toTime.UTC().Unix()
		nowUnix := time.Now().UTC().Unix()
		// if kchart data of the query time period not exist, then use last
		if nil != err || nil == fromDataPtr || 0 == fromDataPtr.ResNum() {
			log.Errorf(err, "query exchange kchart error or kchart data not exist, will use last record value to store into DB")
			for toTimeUnix < nowUnix { // traverse the time period and insert the default value for kchart data
				lastPrice, err := queryLastPeriodClosePrice(toKChartTable, exchangeID)
				if nil != err {
					log.Errorf(err, "query last minute close price failed, will use 0 as last price")
					lastPrice = 0
				}
				insert2KChartTable(toKChartTable, startTimeUnix, exchangeID, lastPrice, lastPrice, lastPrice, lastPrice, 0)
				kchartData := buildKChartDataPiece(startTimeUnix, lastPrice, lastPrice, lastPrice, lastPrice, 0, firstDecimal, secondDecimal)
				kchartDataList = append(kchartDataList, kchartData)
				// update next time span
				startTimeUnix = toTimeUnix
				temp := time.Unix(startTimeUnix, 0).UTC()
				toTime = time.Date(temp.Year(), temp.Month()+1, 1, 0, 0, 0, 0, time.UTC)
				toTimeUnix = toTime.UTC().Unix()
			}
		} else { // exchange transactions exist
			var fromDataIndex int     // this index indicate the first or not of a new time span of 1 minute k graph data
			var traversedTransNum int // indicate that the transactions processed
			// var nowMonth time.Month
			// var monthChecked bool
			for fromDataPtr.NextT() {
				traversedTransNum++
				fromDataTime := util.ConvertDBValueToInt64(fromDataPtr.GetField("time"))
				// if !monthChecked {
				// 	ftime := time.Unix(fromDataTime, 0).UTC()
				// 	// nowMonth = ftime.Month()
				// 	monthChecked = true
				// }
				for fromDataTime >= toTimeUnix {
					if 0 == fromDataIndex { //no transaction in this time span,insert the default value to DB
						lastPrice, err := queryLastPeriodClosePrice(toKChartTable, exchangeID)
						if nil != err {
							log.Errorf(err, "query last minute close price failed, will use 0 as last price")
							lastPrice = 0
						}
						insert2KChartTable(toKChartTable, startTimeUnix, exchangeID, lastPrice, lastPrice, lastPrice, lastPrice, 0)
						kchartData := buildKChartDataPiece(startTimeUnix, lastPrice, lastPrice, lastPrice, lastPrice, 0, firstDecimal, secondDecimal)
						kchartDataList = append(kchartDataList, kchartData)
					} else {
						// insert the 1 minute k graph data to DB
						insert2KChartTable(toKChartTable, startTimeUnix, exchangeID, open, high, low, close, volume)
						kchartData := buildKChartDataPiece(startTimeUnix, open, high, low, close, volume, firstDecimal, secondDecimal)
						kchartDataList = append(kchartDataList, kchartData)
						volume = 0
						fromDataIndex = 0
					}
					// update time span to next
					startTimeUnix = toTimeUnix
					temp := time.Unix(startTimeUnix, 0).UTC()
					toTime = time.Date(temp.Year(), temp.Month()+1, 1, 0, 0, 0, 0, time.UTC)
					toTimeUnix = toTime.UTC().Unix()
				}

				h := util.ConvertDBValueToInt64(fromDataPtr.GetField("high"))
				l := util.ConvertDBValueToInt64(fromDataPtr.GetField("low"))
				v := util.Abs(util.ConvertDBValueToInt64(fromDataPtr.GetField("volume")))
				if 0 == fromDataIndex { //if the first transaction, then record open price
					open = util.ConvertDBValueToInt64(fromDataPtr.GetField("open"))
					high = h
					low = l
					volume = v
				} else { // compare the high and low price
					if high < h {
						high = h
					}
					if low > l {
						low = l
					}
					volume += v
				}
				// record the close value
				close = util.ConvertDBValueToInt64(fromDataPtr.GetField("close"))
				// if the last transactoin, then insert to DB and break the traverse
				if fromDataPtr.ResNum() == traversedTransNum {
					insert2KChartTable(toKChartTable, startTimeUnix, exchangeID, open, high, low, close, volume)
					kchartData := buildKChartDataPiece(startTimeUnix, open, high, low, close, volume, firstDecimal, secondDecimal)
					kchartDataList = append(kchartDataList, kchartData)
					break
				}
				fromDataIndex++
			}
		}
		exIDKChartData[strconv.FormatInt(exchangeID, 10)] = kchartDataList
	}
	return exIDKChartData, nil
}

func buildKChartData(fromKChartTable string, toKChartTable string, timeSpan time.Duration, fromTime time.Time) (map[string][]*entity.KChartData, error) {
	log.Debugf("In buildKChartData, from time is: %v, fromTable: %v, toTable: %v", fromTime.UTC().Unix(), fromKChartTable, toKChartTable)
	strSQL := fmt.Sprintf(`select id, first_token_precision, second_token_precision from market_pair where 1 = 1`) //only query the valid pair??
	exchangeIDPtr, err := mysql.QueryTableData(strSQL)
	if nil != err {
		log.Errorf(err, "Query exchange id error")
		return nil, fmt.Errorf("query exchange id error")
	}
	if nil == exchangeIDPtr || 0 == exchangeIDPtr.ResNum() {
		log.Errorf(err, "exchange id list is nil or empty")
		return nil, fmt.Errorf("query exchange id error")
	}

	exIDKChartData := make(map[string][]*entity.KChartData)
	for exchangeIDPtr.NextT() {
		//traverse the exchange_id and get the corresponding data from exchange transaction table
		//and insert to the kchart table
		kchartDataList := make([]*entity.KChartData, 0)
		exchangeID := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("id"))
		firstDecimal := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("first_token_precision"))
		secondDecimal := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("second_token_precision"))
		fromDataPtr, err := queryKChartData(fromKChartTable, fromTime, time.Now(), exchangeID)
		var open, high, low, close, volume int64
		startTime := fromTime.UTC().Unix()
		toTime := (fromTime.Add(timeSpan)).UTC().Unix()
		now := time.Now().UTC().Unix()
		// if kchart data of the query time period not exist, then use last
		if nil != err || nil == fromDataPtr || 0 == fromDataPtr.ResNum() {
			log.Errorf(err, "query exchange kchart error or kchart data not exist, will use last record value to store into DB")
			for toTime < now { // traverse the time period and insert the default value for the 1 day k graph data
				lastPrice, err := queryLastPeriodClosePrice(toKChartTable, exchangeID)
				if nil != err {
					log.Errorf(err, "Query last minute close price failed, will use 0 as last price")
					lastPrice = 0
				}
				insert2KChartTable(toKChartTable, startTime, exchangeID, lastPrice, lastPrice, lastPrice, lastPrice, 0)
				kchartData := buildKChartDataPiece(startTime, lastPrice, lastPrice, lastPrice, lastPrice, 0, firstDecimal, secondDecimal)
				kchartDataList = append(kchartDataList, kchartData)
				// update next time span
				startTime = toTime
				temp := time.Unix(startTime, 0).UTC()
				toTime = (temp.Add(timeSpan)).UTC().Unix()
			}
		} else { // exchange transactions exist
			var fromDataIndex int     // this index indicate the first or not of a new time span of 1 minute k graph data
			var traversedTransNum int // indicate that the transactions processed
			// var nowMonth time.Month
			// var monthChecked bool
			for fromDataPtr.NextT() {
				traversedTransNum++
				fromDataTime := util.ConvertDBValueToInt64(fromDataPtr.GetField("time"))
				// if !monthChecked {
				// 	ftime := time.Unix(fromDataTime, 0).UTC()
				// 	// nowMonth = ftime.Month()
				// 	monthChecked = true
				// }
				for fromDataTime >= toTime {
					if 0 == fromDataIndex { //no transaction in this time span,insert the default value to DB
						lastPrice, err := queryLastPeriodClosePrice(toKChartTable, exchangeID)
						if nil != err {
							log.Errorf(err, "query last minute close price failed, will use 0 as last price")
							lastPrice = 0
						}
						insert2KChartTable(toKChartTable, startTime, exchangeID, lastPrice, lastPrice, lastPrice, lastPrice, 0)
						kchartData := buildKChartDataPiece(startTime, lastPrice, lastPrice, lastPrice, lastPrice, 0, firstDecimal, secondDecimal)
						kchartDataList = append(kchartDataList, kchartData)
					} else {
						// insert the 1 minute k graph data to DB
						insert2KChartTable(toKChartTable, startTime, exchangeID, open, high, low, close, volume)
						kchartData := buildKChartDataPiece(startTime, open, high, low, close, volume, firstDecimal, secondDecimal)
						kchartDataList = append(kchartDataList, kchartData)
						volume = 0
						fromDataIndex = 0
					}
					// update time span to next
					startTime = toTime
					temp := time.Unix(startTime, 0).UTC()
					toTime = (temp.Add(timeSpan)).UTC().Unix()
				}

				h := util.ConvertDBValueToInt64(fromDataPtr.GetField("high"))
				l := util.ConvertDBValueToInt64(fromDataPtr.GetField("low"))
				v := util.Abs(util.ConvertDBValueToInt64(fromDataPtr.GetField("volume")))
				if 0 == fromDataIndex { //if the first transaction, then record open price
					open = util.ConvertDBValueToInt64(fromDataPtr.GetField("open"))
					high = h
					low = l
					volume = v
				} else { // compare the high and low price
					if high < h {
						high = h
					}
					if low > l {
						low = l
					}
					volume += v
				}
				// record the close value
				close = util.ConvertDBValueToInt64(fromDataPtr.GetField("close"))
				// if the last transactoin, then insert to DB and break the traverse
				if fromDataPtr.ResNum() == traversedTransNum {
					insert2KChartTable(toKChartTable, startTime, exchangeID, open, high, low, close, volume)
					kchartData := buildKChartDataPiece(startTime, open, high, low, close, volume, firstDecimal, secondDecimal)
					kchartDataList = append(kchartDataList, kchartData)
					break
				}
				fromDataIndex++
			}
		}
		exIDKChartData[strconv.FormatInt(exchangeID, 10)] = kchartDataList
	}
	return exIDKChartData, nil
}

//Handle1minKgraphData ...
func Handle1minKChartData(fromTime time.Time) (map[string][]*entity.KChartData, error) {
	log.Infof("Begin handle1minKChartData, from time is: %v", fromTime.UTC().Format(util.DATETIMEFORMAT))
	strSQL := fmt.Sprintf(`select id, first_token_precision, second_token_precision from market_pair where 1 = 1`) //only query the valid pair??
	exchangeIDPtr, err := mysql.QueryTableData(strSQL)
	if nil != err {
		log.Errorf(err, "Query exchange id error")
		return nil, fmt.Errorf("query exchange id error")
	}
	if nil == exchangeIDPtr || 0 == exchangeIDPtr.ResNum() {
		log.Error("exchange id list is nil or empty", nil)
		return nil, fmt.Errorf("query exchange id error")
	}

	exIDKChartData := make(map[string][]*entity.KChartData)
	for exchangeIDPtr.NextT() {
		//traverse the exchange_id and get the corresponding data from exchange transaction table
		//and insert to the 1min kgraph table
		kchartDataList := make([]*entity.KChartData, 0)
		exchangeID := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("id"))
		firstDecimal := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("first_token_precision"))
		secondDecimal := util.ConvertDBValueToInt64(exchangeIDPtr.GetField("second_token_precision"))
		transPtr, err := queryExchangeTransactions(fromTime, exchangeID)
		var open, high, low, close, volume int64
		startTime := fromTime.UTC().Unix()
		toTime := (fromTime.Add(time.Minute)).UTC().Unix()
		now := time.Now().UTC().Unix()
		// if transactons of the query time period not exist, then use default value to fufill the 1 minutes k graph data
		if nil != err || nil == transPtr || 0 == transPtr.ResNum() {
			//if there is no transactions in this 1 minutes period, then use the last minute's close price as the open, high, low and close price for this minute
			lastPrice, err := queryLastPeriodClosePrice(TableMin1, exchangeID)
			if nil != err {
				log.Errorf(err, "Query last minute close price failed, will use 0 as last price")
				lastPrice = 0
			}
			if nil != err {
				log.Errorf(err, "Query exchange transaction error")
			}
			// traverse the time period and insert the last minute value for the 1 minutes k graph data
			for toTime < now {
				insert2KChartTable(TableMin1, startTime, exchangeID, lastPrice, lastPrice, lastPrice, lastPrice, 0)
				kchartData := buildKChartDataPiece(startTime, lastPrice, lastPrice, lastPrice, lastPrice, 0, firstDecimal, secondDecimal)
				kchartDataList = append(kchartDataList, kchartData)
				// update next time span
				startTime = toTime
				temp := time.Unix(startTime, 0).UTC()
				toTime = (temp.Add(time.Minute)).UTC().Unix()
			}
		} else { // exchange transactions exist
			var transIndex int        // this index indicate the first or not of a new time span of transaction
			var traversedTransNum int // indicate that the transactions processed in a time span
			for transPtr.NextT() {
				traversedTransNum++
				createTime := util.ConvertDBValueToInt64(transPtr.GetField("create_time"))
				for createTime >= toTime*1000 { //transaction's create time is unix timestamp in mili second
					if 0 == transIndex { //no transaction in this time span,insert the default value to DB
						lastPrice, err := queryLastPeriodClosePrice(TableMin1, exchangeID)
						if nil != err {
							log.Errorf(err, "Query last minute close price failed, will use 0 as last price")
							lastPrice = 0
						}
						insert2KChartTable(TableMin1, startTime, exchangeID, lastPrice, lastPrice, lastPrice, lastPrice, 0)
						kchartData := buildKChartDataPiece(startTime, lastPrice, lastPrice, lastPrice, lastPrice, 0, firstDecimal, secondDecimal)
						kchartDataList = append(kchartDataList, kchartData)
					} else {
						// insert the 1 minute kchart data to DB
						insert2KChartTable(TableMin1, startTime, exchangeID, open, high, low, close, volume)
						kchartData := buildKChartDataPiece(startTime, open, high, low, close, volume, firstDecimal, secondDecimal)
						kchartDataList = append(kchartDataList, kchartData)
						volume = 0
						transIndex = 0
					}
					// update time span to next
					startTime = toTime
					temp := time.Unix(startTime, 0).UTC()
					toTime = (temp.Add(time.Minute)).UTC().Unix()
				}

				p := util.ConvertDBValueToInt64(transPtr.GetField("price"))
				v := util.ConvertDBValueToInt64(transPtr.GetField("amountA"))
				if 0 == transIndex { //if the first transaction, then record open price
					open = p
					high = p
					low = p
					volume = v
				} else { // compare the high and low price
					if high < p {
						high = p
					}
					if low > p {
						low = p
					}
					volume += v
				}
				// record the close value
				close = p
				// if the last transactoin, then insert to DB and break the traverse
				if transPtr.ResNum() == traversedTransNum {
					insert2KChartTable(TableMin1, startTime, exchangeID, open, high, low, close, volume)
					kchartData := buildKChartDataPiece(startTime, open, high, low, close, volume, firstDecimal, secondDecimal)
					kchartDataList = append(kchartDataList, kchartData)
					break
				}
				transIndex++
			}
		}
		exIDKChartData[strconv.FormatInt(exchangeID, 10)] = kchartDataList
	}
	return exIDKChartData, nil
}

// GetKChartData get kchart data according to the query parameters
func GetKChartData(fromTable string, param *entity.KChartQueryParam, firstDecimal int64, secondDecimal int64) (*entity.KChartResp, error) {
	exchangeID, err := strconv.ParseInt(param.ExchangeID, 10, 64)
	if nil != err {
		log.Errorf(err, "Parameter of exchange_id is not right")
		return nil, fmt.Errorf("parameter of exchange id is not right")
	}
	timeStart, err := strconv.ParseInt(param.TimeStart, 10, 64)
	if nil != err {
		log.Errorf(err, "Parameter of time_start is not right")
		return nil, fmt.Errorf("parameter of time_start is not right")
	}
	timeEnd, err := strconv.ParseInt(param.TimeEnd, 10, 64)
	if nil != err {
		log.Errorf(err, "Parameter of time_end is not right :[%v]\n", err)
		return nil, fmt.Errorf("parameter of time_end is not right")
	}
	fromTime := time.Unix(timeStart, 0).UTC()
	toTime := time.Unix(timeEnd, 0).UTC()
	kgraphDataPtr, err := queryKChartData(fromTable, fromTime, toTime, exchangeID)
	if nil != err {
		return nil, err
	}

	kchartResp := &entity.KChartResp{}
	kchartResp.ExchangeID = param.ExchangeID
	kchartResp.Granu = param.Granu
	kchartResp.TimeEnd = param.TimeEnd
	kchartResp.TimeStart = param.TimeStart
	kchartData := make([]*entity.KChartData, 0)
	for kgraphDataPtr.NextT() {
		t := util.ConvertDBValueToInt64(kgraphDataPtr.GetField("time"))
		o := util.ConvertDBValueToInt64(kgraphDataPtr.GetField("open"))
		h := util.ConvertDBValueToInt64(kgraphDataPtr.GetField("high"))
		l := util.ConvertDBValueToInt64(kgraphDataPtr.GetField("low"))
		c := util.ConvertDBValueToInt64(kgraphDataPtr.GetField("close"))
		v := util.ConvertDBValueToInt64(kgraphDataPtr.GetField("volume"))
		row := buildKChartDataPiece(t, o, h, l, c, v, firstDecimal, secondDecimal)
		kchartData = append(kchartData, row)
	}
	kchartResp.Data = kchartData
	return kchartResp, nil
}

// queryKChartData query the kchart data from kchart table
func queryKChartData(fromTable string, fromTime time.Time, toTime time.Time, exchangeID int64) (*mysql.TronDBRows, error) {
	log.Debugf("begin queryKChartData from %v from time is: [%v], to time is: [%v]\n", fromTable, fromTime.Format(util.DATETIMEFORMAT), toTime.Format(util.DATETIMEFORMAT))
	frTime := fromTime.UTC().Unix()
	tTime := toTime.UTC().Unix()
	strSQL := fmt.Sprintf(`select * from trxmarket.%v where exchange_id = %v and time >= %v and time < %v`, fromTable, exchangeID, frTime, tTime)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf(err, "Query data from %v error", fromTable)
		return nil, fmt.Errorf("query data from %v error", fromTable)
	}
	if dataPtr == nil {
		log.Error("return data from %v is nil", nil)
		return nil, fmt.Errorf("return data from %v is nil", fromTable)
	}
	return dataPtr, nil
}

// func constructKChartData(startTime, open, high, low, close, volume int64) *entity.KChartData {
// 	kchartData := &entity.KChartData{
// 		Time:   strconv.FormatInt(startTime, 10),
// 		Open:   strconv.FormatInt(open, 10),
// 		High:   strconv.FormatInt(high, 10),
// 		Low:    strconv.FormatInt(low, 10),
// 		Close:  strconv.FormatInt(close, 10),
// 		Volume: strconv.FormatInt(volume, 10),
// 	}
// 	return kchartData
// }

//buildKChartDataPiece build one piece of kgraph data
func buildKChartDataPiece(timeStamp int64, open int64, high int64, low int64, close int64, volume int64, firstDecimal int64, secondDecimal int64) *entity.KChartData {
	log.Debugf("before build: timeStamp: %v, open: %v , high: %v , low: %v, close: %v, volume: %v, firstDecimal: %v, secondDecimal: %v", timeStamp, open, high, low, close, volume, firstDecimal, secondDecimal)
	row := &entity.KChartData{}
	row.Time = strconv.FormatInt(timeStamp, 10)
	divisor1 := math.Pow(float64(10), float64(secondDecimal))
	divisor2 := math.Pow(float64(10), float64(firstDecimal))
	log.Debugf("in buildKChartDataPiece, divisor1 is %v", divisor1)
	log.Debugf("in buildKChartDataPiece, divisor1 is %v", divisor2)

	row.Open = strconv.FormatFloat(float64(open)/divisor1, 'f', int(secondDecimal), 64)
	row.High = strconv.FormatFloat(float64(high)/divisor1, 'f', int(secondDecimal), 64)
	row.Low = strconv.FormatFloat(float64(low)/divisor1, 'f', int(secondDecimal), 64)
	row.Close = strconv.FormatFloat(float64(close)/divisor1, 'f', int(secondDecimal), 64)
	row.Volume = strconv.FormatFloat(float64(volume)/divisor2, 'f', int(firstDecimal), 64)
	log.Debugf("after build: timeStamp: %v, open: %v , high: %v , low: %v, close: %v, volume: %v", row.Time, row.Open, row.High, row.Low, row.Close, row.Volume)
	return row
}

//queryLastMinuteClosePrice
func queryLastPeriodClosePrice(tableName string, pairID int64) (int64, error) {
	// temp := time.Unix(fromTime, 0).UTC()
	// lastMinute := (temp.Add(-timeSpan)).UTC().Unix()

	strSQL := fmt.Sprintf(`select close from %v where exchange_id = %v order by time desc limit 1`, tableName, pairID)
	lastTimeSpanPricePtr, err := mysql.QueryTableData(strSQL)
	if nil != err {
		log.Errorf(err, "Query last time span price error")
		return 0, fmt.Errorf("query last time span price")
	}
	if nil == lastTimeSpanPricePtr || 0 == lastTimeSpanPricePtr.ResNum() {
		log.Errorf(err, "last time span price is nil or empty")
		return 0, fmt.Errorf("last time span price is nil or empty")
	}

	var price int64
	for lastTimeSpanPricePtr.NextT() {
		price = util.ConvertDBValueToInt64(lastTimeSpanPricePtr.GetField("close"))
	}
	return price, nil
}

func insert2KChartTable(tableName string, ts int64, exchangeID int64, open int64, high int64, low int64, close int64, volume int64) error {
	strInsertSQL := fmt.Sprintf(`insert into %v(time, exchange_id, open, high, low, close, volume)
	values('%v', '%v','%v', '%v', '%v', '%v','%v')`, tableName, ts, exchangeID, open, high, low, close, volume)
	_, _, err := mysql.ExecuteSQLCommand(strInsertSQL)
	if err != nil {
		// log.Errorf("Insert kchart data to %v fail:[%v]  sql:%s, will try update.", tableName, err, strInsertSQL)
		//try update if primary key duplicated
		strUpdateSQL := fmt.Sprintf(`UPDATE %v 
		SET open='%v',high='%v',low='%v',close='%v',volume='%v' where time='%v' and exchange_id='%v'`, tableName, open, high, low, close, volume, ts, exchangeID)
		_, _, err1 := mysql.ExecuteSQLCommand(strUpdateSQL)
		if nil != err1 {
			log.Errorf(err, "UPDATE kchart data to %v error, sql:%s", tableName, strUpdateSQL)
			return err1
		}
		log.Debugf("UPDATE kchart data to %v success, time: [%v]", tableName, ts)
		return nil
	}
	log.Debugf("Insert kchart data to %v success, time: [%v]", tableName, ts)
	return nil
}

// QueryExchangeTransactions ...
func queryExchangeTransactions(fromTime time.Time, exchangeID int64) (*mysql.TronDBRows, error) {
	frTime := fromTime.UTC().Unix() * 1000 //milisecond

	strSQL := fmt.Sprintf(`select amountA, price, create_time from market_transaction
						   where pair_id = %v and create_time >= %v and status=0 order by create_time`,
		exchangeID, frTime)
	//fmt.Printf("strSQL: %v: ", strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf(err, "Query exchange transaction error")
		return nil, fmt.Errorf("query exchange transaction error")
	}
	if dataPtr == nil {
		log.Errorf(err, "exchange transaction list is nil")
		return nil, fmt.Errorf("query exchange transaction error")
	}
	return dataPtr, nil
}
