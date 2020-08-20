package handler

import (
	"github.com/adschain/log"
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tradehome-engin/core/utils"
	"github.com/wlcy/tradehome-service/buffer2"
	"github.com/wlcy/tradehome-service/common/contract"
	"github.com/wlcy/tradehome-service/common/errno"
	"github.com/wlcy/tradehome-service/common/redis"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/service"
	"github.com/wlcy/tradehome-service/util"
	"math"
	"strconv"
	"time"
)

//添加交易对
func AddMarketPair(c *gin.Context) {
	//请求头密码校验
	key := c.Request.Header.Get("Key")
	if key != "Tron@123456" {
		SendResponse(c, errno.ErrValidation, nil)
		return
	}
	//处理请求参数
	req := &entity.AddMarketPairReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Errorf(err, "AddTransPair API bind AddTransactionPairReq error")
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	//保存
	if err := service.AddMarketPair(req); err != nil {
		log.Error("AddTransactionPair API service execute error", err)
		SendResponse(c, errno.InternalServerError, nil)
		return
	}
	//操作成功
	SendResponse(c, nil, nil)
}

//查询交易对列表
func MarketPairList(c *gin.Context) {
	//获取分页参数
	pageInfo, errnoResult := getPageInfo(c)
	if errnoResult != nil {
		SendResponse(c, errnoResult, nil)
		return
	}

	pairs := service.MarketPairList(pageInfo)

	SendResponse(c, nil, pairs)
}

//最近交易记录
func LatestOrderList(c *gin.Context) {
	req := &entity.MarketHisReq{}
	if err := c.ShouldBind(req); err != nil {
		log.Error(" API bind MarketHisReq error", err)
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	if req.Start < 0 {
		req.Start = 0
	}
	if req.Limit > 50 {
		req.Limit = 50
	}

	key := GetQueryKey("LatestOrderList", req)
	data := LoadBuffer(key, 3)
	if nil != data {
		SendResponse(c, nil, data)
		return
	}

	resp, err := service.LatestOrderList(req)
	if err != nil {
		SendResponse(c, errno.InternalServerError, nil)
		return
	}

	if len(resp.Rows) > 0 {
		StoreBuffer(key, resp)
	}
	SendResponse(c, nil, resp)
}

func UserOrderList(c *gin.Context) {
	//获取参数
	req, errnoResult := getUserOrderListReq(c)
	if errnoResult != nil {
		SendResponse(c, errnoResult, nil)
		return
	}
	key := GetQueryKey("UserOrderList", req)
	data := LoadBuffer(key, 3)
	if nil != data {
		SendResponse(c, nil, data)
		return
	}
	//获取数据
	info, err := service.UserOrderAllList(req)
	if err != nil {
		log.Error("UserOrderList API service execute error", err)
		SendResponse(c, errno.InternalServerError, nil)
		return
	}

	if len(info.Rows) > 0 {
		StoreBuffer(key, info)
	}
	SendResponse(c, nil, info)
}

//GetAllOrderList 获取redis中所有的盘口信息
func GetAllOrderList(c *gin.Context) {
	log.Debugf("Hello /api/exchange/common/orderall")
	SendResponse(c, nil, redis.GetOrderList())
	return
}

//GetAllOrderList 获取redis中所有的盘口信息
func GetAllOrderList2(c *gin.Context) {
	orderList10 := redis.GetOrderList2(1) //查询10盘口信息
	orderList20 := redis.GetOrderList2(2) //查询20盘口信息
	if nil == orderList20 {
		SendResponse(c, nil, orderList10)
		return
	}
	if nil == orderList10 {
		SendResponse(c, nil, orderList20)
		return
	}
	for exchangeID, orderList := range orderList10 { //合并10，20盘口信息
		orderList20[exchangeID] = orderList
	}

	SendResponse(c, nil, orderList20)
	return
}

//GetSmartOrderInfo 获取智能合约中的订单数据
func GetSmartOrderInfo(c *gin.Context) {

	orderID := c.Param("orderID")
	order := util.ConvertDBValueToInt64(orderID)

	orderInfo, err := contract.CallOrderInfo(order)
	if err != nil {
		log.Errorf(err, "GetSmartOrderInfo  error")
		SendResponse(c, errno.InternalServerError, nil)
		return
	}
	data := &contract.SmartOrderInfoTwo{}
	data.OrderID = orderInfo.OrderID
	data.Status = orderInfo.Status
	data.OrderType = orderInfo.OrderType
	data.OwnerAddress = orderInfo.OwnerAddress
	data.TokenA = orderInfo.TokenA
	data.TokenB = orderInfo.TokenB
	precisionA := int(buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[orderInfo.TokenA])
	precisionB := int(buffer2.GetPrecisionBuffer().GetTokenAddrPrecisionMap()[orderInfo.TokenB])
	log.Infof("GetSmartOrderInfo tokenA:%v, tokenB:%v, precisionA:%v, precisionB:%v", data.TokenA, data.TokenB, precisionA, precisionB)
	data.AmountA = float64(orderInfo.AmountA) / math.Pow10(precisionA)
	data.AmountB = float64(orderInfo.AmountB) / math.Pow10(precisionB)
	data.OrderFills = float64(orderInfo.OrderFills) / math.Pow10(precisionA)
	data.Turnover = float64(orderInfo.Turnover) / math.Pow10(precisionB)
	data.ChannelId = orderInfo.ChannelId

	SendResponse(c, nil, data)
	return

}

//GetSmartOrderInfo 获取智能合约中的订单数据
func GetSmartOrderInfo2(c *gin.Context) {
	var err error
	var orderInfo interface{}
	orderID := c.Param("orderID")
	order := util.ConvertDBValueToInt64(orderID)
	utils.TestNet = false
	contract.EventLogURL = "https://event.tronscan.org/api/v1"

	// utils.NetName = utils.NetShasta
	// config.SmartNode = "54.236.37.243:50051"
	// config.SmartOwnerAddr = "TVXH2dRgh7abQoQx5UyzakWheqPXvaKprL"
	// config.SmartOwnerAddr10 = "TVXH2dRgh7abQoQx5UyzakWheqPXvaKprL"
	// config.SmartContractAddr = "TPb1vyaZv59Mte7ABaRqTGCk4TaX1MzBUT"
	// config.SmartContractAddr10 = "THxN2BbBHp5nam7ZQWPrgLoAju3KTDZcNg"
	// config.SmartPrivateKey = "5B47A104364C66DB7C5544034AE3362915E4B6A77A08676935DCD95D290CCE6E"
	contract.InitTradeSmart2(order)
	if order >= 1000000000000 {
		orderInfo, err = contract.CallOrderInfo10(order)
	} else {
		orderInfo, err = contract.CallOrderInfo(order)
	}
	if err != nil {
		log.Errorf(err, "GetSmartOrderInfo  error")
		SendResponse(c, errno.InternalServerError, nil)
		return
	}
	SendResponse(c, nil, orderInfo)
	return

}

//盘口列表
/*func GetOrderList(c *gin.Context) {
	startTime := time.Now().UTC().UnixNano()
	log.Print("GetOrderList request")
	pairIDStr := c.Param("pairID")
	log.Debugf("GetOrderList pairID :[%v]", pairIDStr)
	if pairIDStr == "" {
		handler.SendResponse(c, errno.ErrParam, nil)
		return
	}
	pairID, err := strconv.ParseInt(pairIDStr, 10, 64)
	if err != nil {
		handler.SendResponse(c, errno.ErrParam, nil)
		return
	}
	key := GetQueryKey("GetOrderList", pairID)

	data := LoadBuffer(key, 3)
	if nil != data {
		handler.SendResponse(c, nil, data)
		endTime := time.Now().UTC().UnixNano()
		//log.Infof("get lastest order list cost time(milli seconds): %v\n", (endTime-startTime)/1000000)
		log.Debugf("get order list from buffer cost time(milli seconds): %v, result:[%v]\n", (endTime-startTime)/1000000, utils.ToJSONStr(data))
		return
	}
	ret := service.GetOrderListByPairID(pairID)
	handler.SendResponse(c, nil, ret)
	endTime := time.Now().UTC().UnixNano()
	StoreBuffer(key, ret)
	log.Debugf("get order list cost time(milli seconds): %v\n", (endTime-startTime)/1000000)
	return
}*/

//盘口列表
func GetOrderList2(c *gin.Context) {
	startTime := time.Now().UTC().UnixNano()
	log.Info("GetOrderList2 request")
	pairIDStr := c.Param("pairID")
	log.Debugf("GetOrderList2 pairID :[%v]", pairIDStr)
	if pairIDStr == "" {
		SendResponse(c, errno.ErrBind, nil)
		log.Debugf("GetOrderList2 pairIDStr empty")
		return
	}
	pairID, err := strconv.ParseInt(pairIDStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ErrBind, nil)
		log.Debugf("GetOrderList2 ParseInt error\n")
		return
	}

	key := GetQueryKey("GetOrderList", pairID)
	log.Debugf("get order list2 start LoadBuffer\n")
	data := LoadBuffer(key, 3)
	if nil != data {
		SendResponse(c, nil, data)
		endTime := time.Now().UTC().UnixNano()
		//log.Infof("get lastest order list cost time(milli seconds): %v\n", (endTime-startTime)/1000000)
		log.Debugf("get order list2 from buffer cost time(milli seconds): %v, result:[%v]\n", (endTime-startTime)/1000000, utils.ToJSONStr(data))
		return
	}
	ret := service.GetOrderListByPairID2(pairID)
	SendResponse(c, nil, ret)
	endTime := time.Now().UTC().UnixNano()
	StoreBuffer(key, ret)
	log.Debugf("get order list2 cost time(milli seconds):Price:%v, %v, ret:[%v]\n", ret.Price, (endTime-startTime)/1000000, *ret)
	return
}

//统一获取分页信息
func getUserOrderListReq(c *gin.Context) (*entity.UserOrderListReq, *errno.Errno) {
	req := &entity.UserOrderListReq{}
	if err := c.ShouldBindQuery(req); err != nil {
		log.Error(" API bind UserOrderListReq error", err)
		return nil, errno.ErrBind
	}
	if req.Start < 0 {
		req.Start = 0
	}
	//限制条数最多为50
	if req.Limit > 50 {
		req.Limit = 50
	}
	return req, nil
}

//统一获取分页信息
func getPageInfo(c *gin.Context) (*entity.PageInfo, *errno.Errno) {
	req := &entity.PageInfo{}
	if err := c.ShouldBindQuery(req); err != nil {
		log.Error(" API bind PageInfo error", err)
		return nil, errno.ErrBind
	}
	if req.Start < 0 {
		req.Start = 0
	}
	//限制条数最多为50
	if req.Limit > 50 {
		req.Limit = 50
	}
	return req, nil
}

//GetTopPriceByExchangeID 根据交易对id获取盘口最高价和最低价
func GetTopPriceByExchangeID(c *gin.Context) {
	pairIDStr := c.Param("exchangeID")
	if pairIDStr == "" {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	pairID, err := strconv.ParseInt(pairIDStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	SendResponse(c, nil, service.GetTopPriceByPairID(pairID))
}

//GetTopPriceByExchangeID 根据交易对id获取盘口最高价和最低价
func GetTopPriceByExchangeID2(c *gin.Context) {
	pairIDStr := c.Param("exchangeID")
	if pairIDStr == "" {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	pairID, err := strconv.ParseInt(pairIDStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	SendResponse(c, nil, service.GetTopPriceByPairID2(pairID))
}

func AddChannelId(c *gin.Context) {
	//请求头密码校验
	key := c.Request.Header.Get("Key")
	if key != "Tron@123456" {
		SendResponse(c, errno.ErrValidation, nil)
		return
	}
	//处理请求参数
	req := &entity.AddChannelIdReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error("AddChannelId API bind AddChannelIdReq error", err)
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	if "" == req.ChannelId {
		req.ChannelId = "0"
	}
	//保存
	if err := service.AddChannelId(req); err != nil {
		log.Error("AddChannelId API service execute error", err)
		SendResponse(c, errno.InternalServerError, nil)
		return
	}
	//操作成功
	SendResponse(c, nil, nil)
}
