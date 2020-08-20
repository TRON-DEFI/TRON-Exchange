package service

import (
	"encoding/json"
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-engin/web/common"
	"github.com/wlcy/tradehome-service/buffer"
	"github.com/wlcy/tradehome-service/buffer2"
	"github.com/wlcy/tradehome-service/common/contract"
	"github.com/wlcy/tradehome-service/common/redis"
	"github.com/wlcy/tradehome-service/convert"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/module"
	"math"
	"strconv"
	"strings"
)

func AddMarketPair(req *entity.AddMarketPairReq) (err error) {
	logData, _ := json.Marshal(req)
	log.Debugf("AddMarketPair AddMarketPairReq :[%v]", string(logData))
	pair := entity.MarketPair{
		FisrtShortName:  req.FirstShortName,
		FirstTokenName:  req.FirstTokenName,
		FirstTokenAddr:  req.FirstTokenAddr,
		SecondTokenName: req.SecondTokenName,
		SecondShortName: req.SecondShortName,
		SecondTokenAddr: req.SecondTokenAddr,
		Price:           req.Price,
		Unit:            req.Unit,
		PairType:        req.PairType}

	pairType := pair.PairType
	pair.SecondPrecision = callGetDecimal(req.SecondTokenAddr)
	//价格乘以精度
	pair.Price = math.Pow10(int(pair.SecondPrecision)) * pair.Price
	log.Info("  price exchange ")

	if 1 == pairType { //10 token pair
		pair.FirstPrecision, err = callGet10Decimal(req.FirstTokenAddr)
		if nil != err {
			log.Errorf(err, "get 10 token decimal error")
			return fmt.Errorf("get 10 token decimal error")
		}
		pair.SecondTokenAddr = ""
	} else if 2 == pairType { //20 token pair
		pair.FirstPrecision = callGetDecimal(req.FirstTokenAddr)
		pair.SecondTokenAddr = "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb"
	} else {
		log.Error("pairType of request is not legal", nil)
		return fmt.Errorf("pair type of request if not legal, please double check")
	}

	pairID, err := module.AddMarketPair(&pair)
	if err == nil {
		//初始化Redis中的缓存价格信息
		intPrice := int64(pair.Price)
		log.Debug("  init redis price ")
		redis.UpdateExchangeRealPrice(intPrice, intPrice, pairID, 0)
	}
	return err
}
func callGetDecimal(addr string) int64 {
	log.Debug("  callGetDecimal ")
	//如果是trx address 直接返回6
	if strings.Compare(addr, "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb") == 0 {
		return 6
	} else {
		decimal, err := contract.CallGetDecimal(addr)
		if err != nil {
			log.Error("callGetDecimal error", err)
		}
		return decimal
	}
}

func callGet10Decimal(id string) (int64, error) {
	assetInfo, err := common.GetWalletClient().GetAssetIssueByID(id)
	if err != nil || assetInfo == nil {
		log.Errorf(err, "get 10 token decimal error")
		return 0, err
	}
	return int64(assetInfo.Precision), nil
}

func MarketPairList(pageInfo *entity.PageInfo) *entity.MarketPairListInfo {

	info := buffer.GetRecordBuffer().GetMarketPairList()
	/*if info == nil {
		return nil
	}
	oriLen := len(info.Rows)
	if info.Total == 0 || oriLen == 0 {
		return info
	}
	//根据分页信息截取列表
	//起始index超过尾部,返回nil
	if pageInfo.Start > oriLen-1 {
		info.Rows = nil
		info.Total = 0

		//start+limit超过总长度
	} else if pageInfo.Start+pageInfo.Limit > oriLen {
		info.Rows = info.Rows[pageInfo.Start:]
		info.Total = int64(oriLen - pageInfo.Start)
	} else {
		info.Rows = info.Rows[pageInfo.Start : pageInfo.Limit+pageInfo.Start]
		info.Total = int64(pageInfo.Limit)
	}*/
	return info
}

/*func UserOrderAllList(req *entity.UserOrderListReq) *entity.UserOrderListInfo {
	if "" == req.ChannelId {
		req.ChannelId = "0"
	}
	info := buffer.GetRecordBuffer().GetUserOrder(req.UAddr, req.Status, req.Start, req.Limit, req.FTokenAddr, req.STokenAddr, req.ChannelId)
	return info
}*/

func UserOrderAllList(req *entity.UserOrderListReq) (*entity.UserOrderListInfo, error) {
	dataPtr, totalCount, err := module.UserOrderList(req.UAddr, req.Status, &entity.PageInfo{Start: req.Start, Limit: req.Limit}, req.FTokenAddr, req.STokenAddr)
	if err == nil && dataPtr != nil {
		return &entity.UserOrderListInfo{Rows: convert.Convert2UserOrderRowInfos(dataPtr), Total: totalCount}, nil
	}
	return &entity.UserOrderListInfo{}, err
}

func LatestOrderList(req *entity.MarketHisReq) (*entity.LatestOrderListInfo, error) {
	pageInfo := &entity.PageInfo{}
	pageInfo.Start = req.Start
	pageInfo.Limit = req.Limit
	resp := &entity.LatestOrderListInfo{}
	list, total, err := module.LatestOrderList(pageInfo, req.PairID)
	if err != nil {
		return resp, err
	}
	resp.Total = total
	resp.Rows = list
	return resp, nil
}

func GetOrderListByPairID(pairID int64) *entity.OrderListResp {
	orderListMap := redis.GetOrderList()
	if orderListMap != nil {
		logDataB, _ := json.Marshal(orderListMap)
		log.Debugf("GetOrderListByPairID orderListMap :[%v]", string(logDataB))
		return genOrderRecord(orderListMap[pairID])
	} else {
		data := &entity.OrderListResp{}
		buy := make([]*entity.CommonOrder, 0)
		sell := make([]*entity.CommonOrder, 0)
		data.Buy = buy
		data.Sell = sell
		data.Price = "0"
		return data
	}
}

func GetOrderListByPairID2(pairID int64) *entity.OrderListResp {

	if nil == buffer2.GetMarketPairBuffer() || nil == buffer2.GetMarketPairBuffer().GetMarketPairs() {
		log.Error("GetOrderListByPairID2: Get market pair from buffer returns nil", nil)
		data := &entity.OrderListResp{}
		buy := make([]*entity.CommonOrder, 0)
		sell := make([]*entity.CommonOrder, 0)
		data.Buy = buy
		data.Sell = sell
		data.Price = "0"
		return data
	}

	marketPair, ok := buffer2.GetMarketPairBuffer().GetMarketPairs()[pairID]
	if !ok || nil == marketPair {
		log.Error("GetOrderListByPairID2: Market pair doesn't have pair ID", nil)
		data := &entity.OrderListResp{}
		buy := make([]*entity.CommonOrder, 0)
		sell := make([]*entity.CommonOrder, 0)
		data.Buy = buy
		data.Sell = sell
		data.Price = "0"
		return data
	}

	pairType := marketPair.PairType
	log.Debugf("GetOrderListByPairID2 pairType :%v", pairType)
	orderListMap := redis.GetOrderList2(pairType)
	if orderListMap != nil {
		logDataB, _ := json.Marshal(orderListMap)
		log.Debugf("GetOrderListByPairID orderListMap :[%v]", string(logDataB))
		return genOrderRecord(orderListMap[pairID])
	} else {
		log.Debugf("GetOrderListByPairID orderListMap  nil")
		data := &entity.OrderListResp{}
		buy := make([]*entity.CommonOrder, 0)
		sell := make([]*entity.CommonOrder, 0)
		data.Buy = buy
		data.Sell = sell
		data.Price = "0"
		return data
	}
}

//GetTopPriceByPairID ...
func GetTopPriceByPairID(pairID int64) *entity.TopPrice {
	topPrice := &entity.TopPrice{}
	topPrice.ExchangeID = pairID
	orderListMap := redis.GetOrderList()
	for exchangeID, orderList := range orderListMap {
		if pairID == exchangeID {
			precision := int(buffer2.GetPrecisionBuffer().GetPairIDSecondPrecisionMap()[pairID])
			if orderList != nil && len(orderList.Buy) > 0 {
				topPrice.BuyHighPrice = float64(orderList.Buy[0].Price) / math.Pow10(precision)
			}
			if orderList != nil && len(orderList.Sell) > 0 {
				topPrice.SellLowPrice = float64(orderList.Sell[0].Price) / math.Pow10(precision)
			}
			price, _ := strconv.ParseInt(orderList.Price, 10, 64)
			topPrice.Price = float64(price) / math.Pow10(precision)
			break
		}
	}
	return topPrice
}

//GetTopPriceByPairID ...
func GetTopPriceByPairID2(pairID int64) *entity.TopPrice {

	if nil == buffer2.GetMarketPairBuffer() || nil == buffer2.GetMarketPairBuffer().GetMarketPairs() {
		log.Error("Get market pair from buffer returns nil", nil)
		return nil
	}
	marketPair, ok := buffer2.GetMarketPairBuffer().GetMarketPairs()[pairID]
	if !ok || nil == marketPair {
		log.Error("Market pair doesn't have pair ID", nil)
		return nil
	}

	pairType := marketPair.PairType

	topPrice := &entity.TopPrice{}
	orderListMap := redis.GetOrderList2(pairType)
	for exchangeID, orderList := range orderListMap {
		if pairID == exchangeID {
			precision := int(buffer2.GetPrecisionBuffer().GetPairIDSecondPrecisionMap()[pairID])
			if orderList != nil && len(orderList.Buy) > 0 {
				topPrice.BuyHighPrice = float64(orderList.Buy[0].Price) / math.Pow10(precision)
			}
			if orderList != nil && len(orderList.Sell) > 0 {
				topPrice.SellLowPrice = float64(orderList.Sell[0].Price) / math.Pow10(precision)
			}
			price, _ := strconv.ParseInt(orderList.Price, 10, 64)
			topPrice.Price = float64(price) / math.Pow10(precision)
			break
		}
	}
	return topPrice
}

func genOrderRecord(oriData *redis.OrderListRedis) *entity.OrderListResp {
	//根据pairID获取精度 对所有 的价格做处理
	data := &entity.OrderListResp{}
	buy := make([]*entity.CommonOrder, 0)
	sell := make([]*entity.CommonOrder, 0)
	if oriData != nil {
		oriDataB, _ := json.Marshal(oriData)
		log.Debugf("genOrderRecord oriData :[%v]", string(oriDataB))
		if oriData.Buy != nil && len(oriData.Buy) > 0 {
			for _, OriBuyOrder := range oriData.Buy {
				order := &entity.CommonOrder{
					//OrderID:          OriBuyOrder.OrderID,
					BlockID:          OriBuyOrder.BlockID,
					ExchangeID:       OriBuyOrder.ExchangeID,
					OwnerAddress:     OriBuyOrder.OwnerAddress,
					BsFlag:           OriBuyOrder.BsFlag,
					BuyTokenAddress:  OriBuyOrder.TokenAddressA,
					BuyTokenAmount:   changeAmountAByPrecision(OriBuyOrder.ExchangeID, OriBuyOrder.TokenAmountA),
					SellTokenAddress: OriBuyOrder.TokenAddressB,
					SellTokenAmount:  changeAmountBByPrecision(OriBuyOrder.ExchangeID, OriBuyOrder.TokenAmountB),
					Price:            float64(OriBuyOrder.Price) / math.Pow10(int(buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[OriBuyOrder.TokenAddressB])),
					TrxHash:          OriBuyOrder.TrxHash,
					OrderTime:        OriBuyOrder.OrderTime,
					IsCancel:         OriBuyOrder.IsCancel,
					Status:           OriBuyOrder.Status,
					Amount:           changeAmountAByPrecision(OriBuyOrder.ExchangeID, OriBuyOrder.AmountA),
					CurTurnover:      changeAmountBByPrecision(OriBuyOrder.ExchangeID, OriBuyOrder.TurnOver),
				}
				buy = append(buy, order)
			}
		}
		if oriData.Sell != nil && len(oriData.Sell) > 0 {
			for _, OriBuyOrder := range oriData.Sell {
				order := &entity.CommonOrder{
					//OrderID:          OriBuyOrder.OrderID,
					BlockID:          OriBuyOrder.BlockID,
					ExchangeID:       OriBuyOrder.ExchangeID,
					OwnerAddress:     OriBuyOrder.OwnerAddress,
					BsFlag:           OriBuyOrder.BsFlag,
					BuyTokenAddress:  OriBuyOrder.TokenAddressA,
					BuyTokenAmount:   changeAmountAByPrecision(OriBuyOrder.ExchangeID, OriBuyOrder.TokenAmountA),
					SellTokenAddress: OriBuyOrder.TokenAddressB,
					SellTokenAmount:  changeAmountBByPrecision(OriBuyOrder.ExchangeID, OriBuyOrder.TokenAmountB),
					Price:            float64(OriBuyOrder.Price) / math.Pow10(int(buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[OriBuyOrder.TokenAddressB])),
					TrxHash:          OriBuyOrder.TrxHash,
					OrderTime:        OriBuyOrder.OrderTime,
					IsCancel:         OriBuyOrder.IsCancel,
					Status:           OriBuyOrder.Status,
					Amount:           changeAmountAByPrecision(OriBuyOrder.ExchangeID, OriBuyOrder.AmountA),
					CurTurnover:      changeAmountBByPrecision(OriBuyOrder.ExchangeID, OriBuyOrder.TurnOver),
				}
				sell = append(sell, order)
			}
		}
		if len(buy) > 50 {
			buy = buy[0:50]
		}
		if len(sell) > 50 {
			sell = sell[0:50]
		}
		data.Buy = buy
		data.Sell = sell
		data.Price = oriData.Price
	} else {
		data.Buy = buy
		data.Sell = sell
		data.Price = "0"
		log.Debugf("genOrderRecord oriData = nil")
	}

	return data
}
func getCurTurnover(pairID int64, amount int64, buyTokenAmount int64) float64 {

	return changeAmountAByPrecision(pairID, amount-buyTokenAmount)
}
func changeAmountAByPrecision(pairID int64, amountA int64) float64 {
	precision := int(buffer2.GetPrecisionBuffer().GetPairIDFirstPrecisionMap()[pairID])

	return float64(amountA) / math.Pow10(precision)
}

func changeAmountBByPrecision(pairID int64, amountB int64) float64 {
	precision := int(buffer2.GetPrecisionBuffer().GetPairIDSecondPrecisionMap()[pairID])
	return float64(amountB) / math.Pow10(precision)

}

func AddChannelId(req *entity.AddChannelIdReq) (err error) {
	logData, _ := json.Marshal(req)
	log.Debugf("AddChannelId AddChannelIdReq :[%v]", string(logData))
	channelId := entity.ChannelId{
		Hash:      req.Hash,
		OrderId:   req.OrderId,
		ChannelId: req.ChannelId,
	}
	//log.Info("  price exchange ")
	pairID, err := module.AddChannelId(&channelId)
	log.Debugf("pairID :[%v]", string(pairID))
	// if err == nil {
	// 	//初始化Redis中的缓存价格信息
	// 	intPrice := int64(pair.Price)
	// 	log.Debug("  init redis price ")
	// 	exchange.UpdateExchangeRealPrice(intPrice, intPrice, pairID, 0)
	// }
	return err
}
