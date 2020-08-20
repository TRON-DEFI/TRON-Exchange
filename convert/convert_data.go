package convert

import (
	"fmt"
	"math"
	"strconv"

	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/buffer2"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/util"
)

//GetSchedule ...
func GetSchedule(fTokenBalanceStr string, curTurnover string) string {
	var schedule float64
	ct := util.ConvertDBValueToFloat64(curTurnover)
	tokenBalance := util.ConvertDBValueToFloat64(fTokenBalanceStr)
	if tokenBalance > 0 {
		schedule = ct / tokenBalance
	}
	log.Debugf("GetSchedule curTurnover :[%v]--cur:[%f]/tokenBalance:[%f]=schedule:[%v]", curTurnover, ct, tokenBalance, schedule)
	return fmt.Sprintf("%.4f", schedule)
}

func Convert2PirceInfMap(dataPtr *mysql.TronDBRows, volume24hMap map[int64]float64, amountA24hMap map[int64]float64) map[int64]*entity.PriceInf {
	priceInfs := make(map[int64]*entity.PriceInf, 0)
	for dataPtr.NextT() {
		priceInf := &entity.PriceInf{}
		precisionF := buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[dataPtr.GetField("first_token_addr")]
		precisionS := buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[dataPtr.GetField("second_token_addr")]
		log.Debugf("get precisionF:[%d],precisionS:[%d],firstaddress:[%v],secondaddress:[%v]", precisionF, precisionS, dataPtr.GetField("first_token_address"), dataPtr.GetField("second_token_address"))
		priceInf.LowestPrice24h = util.ConvertDBValueToFloat64(dataPtr.GetField("minp")) / math.Pow10(int(precisionS))
		priceInf.HighestPrice24h = util.ConvertDBValueToFloat64(dataPtr.GetField("maxp")) / math.Pow10(int(precisionS))
		pid, _ := strconv.ParseInt(dataPtr.GetField("pair_id"), 10, 64)
		priceInf.Volume24h = volume24hMap[pid] // / math.Pow10(int(precisionF))  //amountB 不需要除以精度
		priceInf.Amount24h = amountA24hMap[pid]
		log.Debugf("orderInfo:[%#v]", priceInf)
		priceInfs[pid] = priceInf
	}

	return priceInfs
}
func GetVolume24hMap(dataPtr *mysql.TronDBRows) map[int64]float64 {
	volume24hMap := make(map[int64]float64, 0)
	for dataPtr.NextT() {
		tokenAddressA := dataPtr.GetField("first_token_address")
		tokenAddressB := dataPtr.GetField("second_token_address")
		exchangeName := fmt.Sprintf("%v-%v", tokenAddressA, tokenAddressB)
		if exchangeID, ok := buffer2.GetExchangeBuffer().GetExchangeIDByAddr(exchangeName); ok {
			precision := int(buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[tokenAddressB])
			volume24hMap[exchangeID] = util.ConvertDBValueToFloat64(dataPtr.GetField("sv")) / (math.Pow10(precision) * 2)
		}
		//pid, _ := strconv.ParseInt(dataPtr.GetField("pair_id"), 10, 64)
		//volume24hMap[pid] = mysql.ConvertDBValueToFloat64(dataPtr.GetField("sv"))
	}
	return volume24hMap
}
func GetAmountA24hMap(dataPtr *mysql.TronDBRows) map[int64]float64 {
	amountA24hMap := make(map[int64]float64, 0)
	for dataPtr.NextT() {
		tokenAddressA := dataPtr.GetField("first_token_address")
		tokenAddressB := dataPtr.GetField("second_token_address")
		exchangeName := fmt.Sprintf("%v-%v", tokenAddressA, tokenAddressB)
		if exchangeID, ok := buffer2.GetExchangeBuffer().GetExchangeIDByAddr(exchangeName); ok {
			precision := int(buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[tokenAddressA])
			amountA24hMap[exchangeID] = util.ConvertDBValueToFloat64(dataPtr.GetField("amount")) / (math.Pow10(precision) * 2)
		}
		//pid, _ := strconv.ParseInt(dataPtr.GetField("pair_id"), 10, 64)
		//amountA24hMap[pid] = mysql.ConvertDBValueToFloat64(dataPtr.GetField("amount"))
	}
	return amountA24hMap
}

func Convert2OrderHistRowInfos(dataPtr *mysql.TronDBRows) []*entity.OrderHisRowInfo {
	rowInfos := make([]*entity.OrderHisRowInfo, 0)

	for dataPtr.NextT() {
		rowInfo := &entity.OrderHisRowInfo{}
		orderID := dataPtr.GetField("order_id")
		tempPrice, _ := strconv.ParseFloat(dataPtr.GetField("price"), 10)
		pairID, _ := strconv.ParseInt(dataPtr.GetField("pair_id"), 10, 64)
		rowInfo.PairID = pairID
		//根据pairID 获取精度,处理price volume
		rowInfo.Price = DealPriceByPairID(tempPrice, pairID)
		//tempVolume, _ := strconv.ParseFloat(dataPtr.GetField("first_token_balance"), 10)
		tempVolume := util.ConvertDBValueToFloat64(dataPtr.GetField("volume"))
		//first_token_balance  存的是交易的余量，如果订单刚好匹配，则使用volume字段
		if tempVolume == 0 {
			tempVolume = GetVolumeByOrderID(orderID)
		}
		rowInfo.Volume = DealVolumeByPairID(tempVolume, pairID)

		rowInfo.Unit = dataPtr.GetField("unit")
		rowInfo.OrderTime = dataPtr.GetField("create_time")
		rowInfo.BlockID = dataPtr.GetField("block_id")

		orderType, _ := strconv.Atoi(dataPtr.GetField("order_type"))
		//根据交易单类型，进行分别设置买方地址和卖方地址
		if orderType == 1 { //买单
			rowInfo.SellerAddr = dataPtr.GetField("owner_address")
			rowInfo.BuyerAddr = dataPtr.GetField("deal_owner_address")
		} else { //卖单
			rowInfo.SellerAddr = dataPtr.GetField("deal_owner_address")
			rowInfo.BuyerAddr = dataPtr.GetField("owner_address")
		}
		rowInfo.OrderType = orderType
		rowInfos = append(rowInfos, rowInfo)
	}
	return rowInfos
}

func GetVolumeByOrderID(orderID string) float64 {
	var volume float64
	strSQL := fmt.Sprintf("select volume from trxmarket.market_transaction where order_id=%v and first_token_balance=0", orderID)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf(err, "queryAllExchangeID error")
		return volume
	}
	if dataPtr == nil {
		log.Error("queryAllExchangeID dataPtr is nil", nil)
		return volume
	}
	//填充数据
	for dataPtr.NextT() {
		volume = util.ConvertDBValueToFloat64(dataPtr.GetField("volume"))
	}
	return volume
}

func QueryResult2MarketPair(dataPtr *mysql.TronDBRows) []*entity.MarketPair {
	tps := make([]*entity.MarketPair, 0)

	for dataPtr.NextT() {
		tp := &entity.MarketPair{}
		tp.ID, _ = strconv.ParseInt(dataPtr.GetField("id"), 10, 64)
		tp.FisrtShortName = dataPtr.GetField("first_short_name")
		tp.FirstTokenName = dataPtr.GetField("first_token_name")
		tp.FirstTokenAddr = dataPtr.GetField("first_token_addr")
		tp.FirstPrecision, _ = strconv.ParseInt(dataPtr.GetField("first_token_precision"), 10, 64)
		tp.SecondShortName = dataPtr.GetField("second_short_name")
		tp.SecondTokenName = dataPtr.GetField("second_token_name")
		tp.SecondTokenAddr = dataPtr.GetField("second_token_addr")
		tp.SecondPrecision, _ = strconv.ParseInt(dataPtr.GetField("second_token_precision"), 10, 64)
		tempPrice, _ := strconv.ParseFloat(dataPtr.GetField("price"), 10)
		//precision := entity.GetPrecisionByPairID(tp.ID)
		precision := tp.SecondPrecision
		coe := math.Pow10(int(precision))
		tp.Price = tempPrice / coe
		log.Infof("ID:%v, tempPrice:%v, precision:%v, coe:%v, price:%v", tp.ID, tempPrice, precision, coe, tp.Price)
		tp.Unit = dataPtr.GetField("unit")
		tp.CreatedAt = dataPtr.GetField("create_time")
		tp.UpdateAt = dataPtr.GetField("update_time")
		tp.PairType, _ = strconv.ParseInt(dataPtr.GetField("pair_type"), 10, 64)
		tp.DefaultIdx, _ = strconv.ParseInt(dataPtr.GetField("default_idx"), 10, 64)
		tps = append(tps, tp)
	}
	return tps

}

func DealPriceByPairID(price float64, pairID int64) float64 {
	precision := buffer2.GetPrecisionBuffer().GetPairIDSecondPrecisionMap()[pairID]
	return price / math.Pow10(int(precision))

}

//计算volume curTurnover 精度
func DealVolumeByPairID(volume float64, pairID int64) float64 {
	precision := buffer2.GetPrecisionBuffer().GetPairIDFirstPrecisionMap()[pairID]
	return volume / math.Pow10(int(precision))

}
func DealPriceByAddr(price float64, addr string) float64 {
	precision := buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[addr]
	return price / math.Pow10(int(precision))

}

//计算volume curTurnover 精度
func DealVolumeByAddr(volume float64, addr string) float64 {
	precision := buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[addr]
	return volume / math.Pow10(int(precision))

}
func Convert2UserOrderRowInfos(dataPtr *mysql.TronDBRows) []*entity.UserOrderRowInfo {
	rowInfos := make([]*entity.UserOrderRowInfo, 0)

	for dataPtr.NextT() {
		rowInfo := &entity.UserOrderRowInfo{}
		//通过sencond token address 获取 精度
		secondTokenAddr := dataPtr.GetField("second_token_address")
		if secondTokenAddr == "" {
			rowInfo.PairType = 1
		} else {
			rowInfo.PairType = 2
		}
		rowInfo.OrderID = util.ConvertStringToInt64(dataPtr.GetField("order_id"), 0)
		firstTokenAddr := dataPtr.GetField("first_token_address")
		tempPrice, _ := strconv.ParseFloat(dataPtr.GetField("price"), 10)
		tempVolume, _ := strconv.ParseFloat(dataPtr.GetField("first_token_balance"), 10)
		//amount := mysql.ConvertDBValueToFloat64(dataPtr.GetField("orderfill")) //成交量
		tempTurnover, _ := strconv.ParseFloat(dataPtr.GetField("cur_turnover"), 10)
		rowInfo.Price = DealPriceByAddr(tempPrice, secondTokenAddr)
		rowInfo.Volume = DealVolumeByAddr(tempVolume, firstTokenAddr)
		//rowInfo.Volume = module.DealVolumeByAddr(amount, firstTokenAddr)

		rowInfo.CurTurnover = DealVolumeByAddr(tempTurnover, secondTokenAddr) //成交额

		//rowInfo.Schedule = module.GetSchedule(dataPtr.GetField("first_token_balance"), dataPtr.GetField("cur_turnover"))
		rowInfo.Schedule = GetSchedule(dataPtr.GetField("first_token_balance"), dataPtr.GetField("orderfill")) //挂单量，，成交量

		rowInfo.ID, _ = strconv.ParseInt(dataPtr.GetField("id"), 10, 64)
		rowInfo.OrderTime = dataPtr.GetField("create_time")
		//rowInfo.OrderID = dataPtr.GetField("order_id")
		rowInfo.OrderType, _ = strconv.ParseInt(dataPtr.GetField("order_type"), 10, 64)
		rowInfo.OrderStatus = util.ConvertStringToInt64(dataPtr.GetField("status"), 0)
		rowInfo.FisrtShortName, rowInfo.SecondShortName = addShortName(dataPtr.GetField("first_token_address"), secondTokenAddr)
		rowInfos = append(rowInfos, rowInfo)
	}
	return rowInfos
}

/*
func Convert2UserOrderRowInfos(dataPtr *mysql.TronDBRows) []*entity.UserOrderRowInfoDetail {
	rowInfos := make([]*entity.UserOrderRowInfoDetail, 0)

	for dataPtr.NextT() {
		rowInfo := &entity.UserOrderRowInfoDetail{}
		//通过sencond token address 获取 精度
		rowInfo.UAddr = dataPtr.GetField("owner_address")
		secondTokenAddr := dataPtr.GetField("second_token_address")
		firstTokenAddr := dataPtr.GetField("first_token_address")
		rowInfo.FTokenAddr = firstTokenAddr
		rowInfo.STokenAddr = secondTokenAddr
		tempPrice, _ := strconv.ParseFloat(dataPtr.GetField("price"), 10)
		tempVolume, _ := strconv.ParseFloat(dataPtr.GetField("first_token_balance"), 10)
		//amount := mysql.ConvertDBValueToFloat64(dataPtr.GetField("orderfill")) //成交量
		tempTurnover, _ := strconv.ParseFloat(dataPtr.GetField("cur_turnover"), 10)
		rowInfo.Price = DealPriceByAddr(tempPrice, secondTokenAddr)
		rowInfo.Volume = DealVolumeByAddr(tempVolume, firstTokenAddr)
		//rowInfo.Volume = module.DealVolumeByAddr(amount, firstTokenAddr)

		rowInfo.CurTurnover = DealVolumeByAddr(tempTurnover, secondTokenAddr) //成交额

		//rowInfo.Schedule = module.GetSchedule(dataPtr.GetField("first_token_balance"), dataPtr.GetField("cur_turnover"))
		rowInfo.Schedule = GetSchedule(dataPtr.GetField("first_token_balance"), dataPtr.GetField("orderfill")) //挂单量，，成交量

		rowInfo.ID, _ = strconv.ParseInt(dataPtr.GetField("id"), 10, 64)
		rowInfo.OrderTime = dataPtr.GetField("create_time")
		rowInfo.OrderID = dataPtr.GetField("order_id")
		rowInfo.OrderType, _ = strconv.ParseInt(dataPtr.GetField("order_type"), 10, 64)
		//rowInfo.OrderStatus = mysql.ConvertStringToInt64(dataPtr.GetField("status"), 0)
		rowInfo.OrderStatus = dataPtr.GetField("status")
		rowInfo.FisrtShortName, rowInfo.SecondShortName = addShortName(dataPtr.GetField("first_token_address"), secondTokenAddr)
		rowInfo.ChannelId = dataPtr.GetField("channel_id")
		if "" == rowInfo.ChannelId {
			rowInfo.ChannelId = "0"
		}
		rowInfos = append(rowInfos, rowInfo)
	}
	return rowInfos
}*/

func addShortName(firstTokenAddress string, secondTokenAddress string) (fisrtShortName string, secondShortName string) {
	key := firstTokenAddress + secondTokenAddress
	info := buffer2.GetMarketPairBuffer().GetMarketPairInfo()[key]
	if info != nil {
		return info.FirstTokenName, info.SecondTokenName
	}
	return "", ""
}
