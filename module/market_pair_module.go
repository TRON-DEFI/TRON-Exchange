package module

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/buffer2"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/entity"
	"strconv"
	"time"
)

func AddMarketPair(pair *entity.MarketPair) (pairID int64, err error) {
	curTime := time.Now().UTC().UnixNano() / 1000000
	strSQL := fmt.Sprintf(`insert into market_pair 
	(first_short_name,first_token_name, first_token_addr,first_token_precision,second_short_name,second_token_name,second_token_addr,second_token_precision, price, unit, pair_type, create_time,update_time)
	 values('%v','%v','%v','%v', '%v','%v', '%v', '%v', '%f', '%v', %v,'%d','%d')`,
		pair.FisrtShortName, pair.FirstTokenName, pair.FirstTokenAddr, pair.FirstPrecision, pair.SecondShortName, pair.SecondTokenName, pair.SecondTokenAddr, pair.SecondPrecision, pair.Price, pair.Unit, pair.PairType, curTime, curTime)
	log.Info(strSQL)
	insertID, _, err := mysql.ExecuteSQLCommand(strSQL)
	if err != nil {
		log.Errorf(err, "AddMarketPair error, sql:%s", strSQL)
		return 0, err
	}
	log.Debugf("AddMarketPair success, id: [%v]", insertID)
	//存储精度
	cachePrecision(pair.SecondPrecision, pair.SecondTokenAddr, insertID, pair.FirstPrecision, pair.FirstTokenAddr)
	return insertID, nil
}

func cachePrecision(sPrecision int64, sTokenAddr string, pairID int64, fPrecision int64, fTokenAddr string) {
	buffer2.GetPrecisionBuffer().GetPairIDSecondPrecisionMap()[pairID] = sPrecision
	buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[sTokenAddr] = sPrecision

	buffer2.GetPrecisionBuffer().GetPairIDFirstPrecisionMap()[pairID] = fPrecision
	buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[fTokenAddr] = fPrecision
}

// 查询列表信息
func queryList(strSQL, filterSQL, sortSQL, pageSQL string) (*mysql.TronDBRows, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Info(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	//fmt.Println(strFullSQL)
	if err != nil {
		log.Errorf(err, "queryList error")
		return nil, err
	}
	if dataPtr == nil {
		log.Error("queryList dataPtr is nil", nil)
		return nil, err
	}
	return dataPtr, nil
}

//查询总条数
func totalCount(tableName, filterSql string) (int64, error) {
	strFullSQL := "select count(1) as total from " + tableName + " mt " + filterSql
	//fmt.Println(strFullSQL)
	log.Info(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf(err, "querySQLCount error")
		return 0, err
	}
	if dataPtr == nil {
		log.Error("querySQLCount dataPtr is nil", nil)
		return 0, err
	}
	var totalNum = int64(0)
	for dataPtr.NextT() {
		tempTotal, _ := strconv.ParseInt(dataPtr.GetField("total"), 10, 64)
		totalNum += tempTotal
	}

	return totalNum, nil
}
