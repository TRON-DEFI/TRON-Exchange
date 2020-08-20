package module

import (
	"fmt"
	"strconv"
	"time"

	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/convert"
	"github.com/wlcy/tradehome-service/entity"
)

func UserOrderList(addr string, status string, pageInfo *entity.PageInfo, fTokenAddr string, sTokenAddr string) (*mysql.TronDBRows, int64, error) {
	strSQL := fmt.Sprintf(`select distinct mo.id,mo.owner_address,mo.order_id,mo.order_type,mo.first_token_balance,mo.price,mo.second_token_balance, `)
	strSQL += fmt.Sprintf(`mo.first_token_address,mo.second_token_address,mo.cur_turnover,mo.orderfill,mo.create_time,mo.status `)
	strSQL += fmt.Sprintf(`from market_order mo `)
	filterSQL := fmt.Sprintf(" where ")
	if addr != "" {
		filterSQL = fmt.Sprintf(filterSQL+" mo.owner_address = '%s'", addr)
	} else {
		filterSQL = fmt.Sprintf(filterSQL + " 1=1")
	}
	if status != "" {
		if status == "0" {
			status = "0,10"
		}
		filterSQL = fmt.Sprintf(filterSQL+" and mo.status in (%s)", status)
	}
	if fTokenAddr != "" {
		filterSQL = fmt.Sprintf(filterSQL+" and mo.first_token_address = '%s'", fTokenAddr)
	}
	if sTokenAddr != "" {
		filterSQL = fmt.Sprintf(filterSQL+" and mo.second_token_address = '%s'", sTokenAddr)
	}
	sortSQL := " order by mo.create_time desc "

	// 分页语句
	pageSQL := ""
	if pageInfo.Limit != 0 {
		pageSQL = fmt.Sprintf(" limit %v, %v", pageInfo.Start, pageInfo.Limit)
	}
	//fmt.Println(strSQL + filterSQL + sortSQL + pageSQL)
	log.Info(strSQL + filterSQL + sortSQL + pageSQL)
	dataPtr, err := queryList(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		//fmt.Println("UserOrderList error")
		return nil, 0, err
	}
	if nil != dataPtr {
		countSQL := fmt.Sprintf("select count(1) as total from (" + strSQL + filterSQL + ") as M")
		//fmt.Println(countSQL)
		dataCntPtr, er := mysql.QueryTableData(countSQL)
		if err != nil {
			log.Errorf(err, "querySQLCount error")
			//fmt.Printf("querySQLCount error :[%v]\n", er)
			return dataPtr, 0, er
		}
		if dataCntPtr == nil {
			log.Error("querySQLCount dataPtr is nil", nil)
			//fmt.Println("querySQLCount dataPtr is nil ")
			return dataPtr, 0, er
		}
		var totalNum = int64(0)
		for dataCntPtr.NextT() {
			tempTotal, _ := strconv.ParseInt(dataCntPtr.GetField("total"), 10, 64)
			totalNum += tempTotal
		}
		//totalCount, _ := totalCount("market_order", filterSQL)
		//fmt.Printf("UserOrderList totalcount = %v\n", totalCount)
		return dataPtr, totalNum, nil
	}
	return nil, 0, nil
}

func MarketPairList(pageInfo *entity.PageInfo) ([]*entity.MarketPair, int64, error) {
	strSQL := fmt.Sprintf(`select id,first_short_name,first_token_name,first_token_addr,first_token_precision,second_short_name,second_token_name,second_token_addr,second_token_precision,price,unit,create_time,update_time,pair_type,default_idx from market_pair `)

	filterSQL := "where is_valid = 1 "

	sortSQL := " order by default_idx "

	// 分页语句
	pageSQL := "" //fmt.Sprintf("limit %v, %v", pageInfo.Start, pageInfo.Limit)
	dataPtr, err := queryList(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		return nil, 0, err
	}
	if nil != dataPtr {
		totalCount, _ := totalCount("market_pair", filterSQL)
		return convert.QueryResult2MarketPair(dataPtr), totalCount, nil
	}
	return nil, 0, nil
}

func GetFirstTokenLogoUrl(firstTokenAddr string) string {
	strSQL := fmt.Sprintf(`select logo_url from market_token_info where address='%v' `, firstTokenAddr)
	dataPtr, err := mysql.QueryTableData(strSQL)

	logoUrl := ""
	if err != nil {
		log.Errorf(err, "GetFistTokenLogoUrl error")
		return ""
	}

	for dataPtr.NextT() {
		logoUrl = dataPtr.GetField("logo_url")
	}
	return logoUrl
}

func PriceInfs(timeLimit int64) map[int64]*entity.PriceInf {
	strSQL := fmt.Sprintf(`
	select trx.pair_id,pairs.first_token_addr,pairs.second_token_addr,min(trx.price) minp,max(trx.price) maxp 
	from  trxmarket.market_transaction trx
    left join trxmarket.market_pair pairs on trx.pair_id=pairs.id
	where trx.create_time > %v and trx.status=0
	group by trx.pair_id,pairs.first_token_addr,pairs.second_token_addr`, timeLimit)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)

	if err != nil {
		log.Errorf(err, "PriceInfs error")
		return nil
	}
	if dataPtr == nil {
		log.Error("PriceInfs dataPtr is nil", nil)
		return nil
	}
	return convert.Convert2PirceInfMap(dataPtr, volume24h(timeLimit), amountA24h(timeLimit))
}

func amountA24h(timeLimit int64) map[int64]float64 {
	//strSQL := fmt.Sprintf(`select pair_id,sum(amountA) as amount from  market_transaction where create_time > %v group by pair_id`, timeLimit)
	strSQL := fmt.Sprintf(`select  first_token_address,second_token_address,sum(orderfill) as amount from  market_order where create_time > %v group by first_token_address,second_token_address`, timeLimit)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)

	if err != nil {
		log.Errorf(err, "volume24h error")
		return nil
	}
	if dataPtr == nil {
		log.Error("volume24h dataPtr is nil", nil)
		return nil
	}
	return convert.GetAmountA24hMap(dataPtr)
}

func volume24h(timeLimit int64) map[int64]float64 {
	//strSQL := fmt.Sprintf(`select pair_id,sum(amountB) sv from  market_transaction where create_time > %v and trx_hash != '' group by pair_id`, timeLimit)
	strSQL := fmt.Sprintf(`select  first_token_address,second_token_address,sum(cur_turnover) as sv from  market_order where create_time > %v group by first_token_address,second_token_address`, timeLimit)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)

	if err != nil {
		log.Errorf(err, "volume24h error")
		return nil
	}
	if dataPtr == nil {
		log.Error("volume24h dataPtr is nil", nil)
		return nil
	}
	return convert.GetVolume24hMap(dataPtr)
}

func AddChannelId(channelId *entity.ChannelId) (pairID int64, err error) {
	curTime := time.Now().UTC().UnixNano() / 1000000
	filterSQL := fmt.Sprintf(" where order_id = %v", channelId.OrderId)
	totalCount, _ := totalCount("market_order_channel", filterSQL)
	if totalCount > 0 {
		updateSQL := fmt.Sprintf(`update market_order_channel set
			trx_hash = '%v', channel_id = '%v',create_time = '%d' where order_id = '%v'`,
			channelId.Hash, channelId.ChannelId, curTime, channelId.OrderId)
		updateID, _, err := mysql.ExecuteSQLCommand(updateSQL)
		if err != nil {
			log.Errorf(err, "AddChannelId error,sql:%s", updateSQL)
			return 0, err
		}
		log.Debugf("AddChannelId success:duplicated order_id[%v]", channelId.OrderId)
		return updateID, nil
	}
	strSQL := fmt.Sprintf(`insert into market_order_channel 
	(order_id,trx_hash, channel_id,create_time)
	 values('%v','%v','%v','%d')`,
		channelId.OrderId, channelId.Hash, channelId.ChannelId, curTime)
	log.Info(strSQL)
	insertID, _, err := mysql.ExecuteSQLCommand(strSQL)
	if err != nil {
		log.Errorf(err, "AddMarketPair error, sql:%s", strSQL)
		return 0, err
	}
	log.Debugf("AddChannelId success, id: [%v]", insertID)
	return insertID, nil
}

func LatestOrderList(pageInfo *entity.PageInfo, pairID int64) ([]*entity.OrderHisRowInfo, int64, error) {
	strSQL := fmt.Sprintf(`
	select mt.order_id,mt.order_type,mt.pair_id, mt.price,mt.order_type , mt.block_id,mt.trx_hash,mt.owner_address,mt.deal_owner_address,
	mt.create_time as create_time, mt.amountA as volume from trxmarket.market_transaction mt `)

	filterSQL := fmt.Sprintf(" where 1=1 and status=0 and mt.pair_id=%v ", pairID)
	sortSQL := "order by mt.create_time desc"
	// 分页语句
	pageSQL := ""
	if pageInfo.Limit != 0 {
		pageSQL = fmt.Sprintf(" limit %v, %v", pageInfo.Start, pageInfo.Limit)
	}

	dataPtr, err := queryList(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		return nil, 0, err
	}
	if nil != dataPtr {
		totalCount, _ := totalCount("market_transaction", filterSQL)
		return convert.Convert2OrderHistRowInfos(dataPtr), totalCount, nil
	}
	return nil, 0, nil
}
