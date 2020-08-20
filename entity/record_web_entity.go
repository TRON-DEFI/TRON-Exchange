package entity

// 添加交易对请求参数entity
type AddMarketPairReq struct {
	FirstTokenName  string  `json:"firstTokenName" form:"firstTokenName" `                     //eg:dice's token name
	FirstTokenAddr  string  `json:"firstTokenAddr" form:"firstTokenAddr" binding:"required"`   //eg:dice's token ContractAddress
	SecondTokenName string  `json:"secondTokenName" form:"secondTokenName" `                   //eg:the token name of trx
	SecondTokenAddr string  `json:"secondTokenAddr" form:"secondTokenAddr"`                    //eg:the token ContractAddress of trx
	Price           float64 `json:"price" form:"price" binding:"required"`                     //eg: 0.000245
	Unit            string  `json:"unit" form:"unit" binding:"required"`                       //eg:DICE
	FirstShortName  string  `json:"firstShortName" form:"firstShortName" binding:"required"`   //第一个token简称
	SecondShortName string  `json:"secondShortName" form:"secondShortName" binding:"required"` //第二个token简称
	PairType        int64   `json:"pairType" form:"pairType"`                                  //交易对类型1或2
}

//交易对列表每行信息
type MarketPairRowInfo struct {
	ID              int64   `json:"id"`              //交易对ID
	Volume          float64 `json:"volume"`          //交易量eg:111.11
	Gain            string  `json:"gain"`            //涨幅eg:-20.3
	Price           float64 `json:"price"`           //价格eg:0.00245
	FirstPrecision  int64   `json:"precision1"`      //第一个token精度
	SecondPrecision int64   `json:"precision2"`      //第二个token精度
	FirstTokenName  string  `json:"tokenName1"`      //第一个token全称
	SecondTokenName string  `json:"tokenName2"`      //第二个token全称
	FirstShortName  string  `json:"shortName1"`      //第一个token简称
	SecondShortName string  `json:"shortName2"`      //第二个token简称
	FirstTokenAddr  string  `json:"tokenAddr1"`      //第一个token地址
	SecondTokenAddr string  `json:"tokenAddr2"`      //第二个token地址
	HighestPrice24h float64 `json:"highestPrice24h"` //24h最高价
	LowestPrice24h  float64 `json:"lowestPrice24h"`  //24h最低价
	Volume24h       float64 `json:"volume24h"`       //24h成交量
	Unit            string  `json:"unit"`            //单位
	PairType        int64   `json:"pairType"`        //交易对类型
	DefaultIdx      int64   `json:"defaultIdx"`      //排序索引
	LogoUrl         string  `json:"logoUrl"`         //
}

//未完成交易列表行信息
type UserOrderRowInfo struct {
	ID              int64   `json:"id"`
	FisrtShortName  string  `json:"shortName1"`  //第一个token简称
	SecondShortName string  `json:"shortName2"`  //第二个token简称
	Volume          float64 `json:"volume"`      //eg:111.11 DICE
	Price           float64 `json:"price"`       //eg:0.00245
	OrderType       int64   `json:"orderType"`   //eg:1-买入，2-卖出
	OrderTime       string  `json:"orderTime"`   //eg:2018-11-10 21:04:01
	OrderID         int64   `json:"orderID"`     //eg:TVkNuE1BYxECWq85d8UR9zsv6WppBns9iH
	Schedule        string  `json:"schedule"`    //eg:23.2
	CurTurnover     float64 `json:"curTurnover"` //成交额
	OrderStatus     int64   `json:"orderStatus"` //订单状态
	PairType        int64   `json:"pairType"`    //10还是20
}

//未完成交易列表行信息
type UserOrderRowInfoDetail struct {
	ID              int64   `json:"id"`
	UAddr           string  `json:"uAddr"`       //uAddr
	FisrtShortName  string  `json:"fShortName"`  //第一个token简称
	SecondShortName string  `json:"sShortName"`  //第二个token简称
	Volume          float64 `json:"volume"`      //eg:111.11 DICE
	Price           float64 `json:"price"`       //eg:0.00245
	OrderType       int64   `json:"orderType"`   //eg:1-买入，2-卖出
	OrderTime       string  `json:"orderTime"`   //eg:2018-11-10 21:04:01
	OrderID         string  `json:"orderID"`     //eg:TVkNuE1BYxECWq85d8UR9zsv6WppBns9iH
	Schedule        string  `json:"schedule"`    //eg:23.2
	CurTurnover     float64 `json:"curTurnover"` //成交额
	OrderStatus     string  `json:"orderStatus"` //订单状态
	FTokenAddr      string  `json:"fTokenAddr"`  //第一个token地址
	STokenAddr      string  `json:"sTokenAddr"`  //第二个token地址
	ChannelId       string  //渠道ID
}

//交易历史列表行信息
type OrderHisRowInfo struct {
	BlockID    string  `json:"blockID"`   //eg:"00000000003ea22f3ccf0ee938556dfaeaead969e27b17ddb3d7895dd07662f1"
	BuyerAddr  string  `json:"buyAddr"`   //eg:a5632020f
	SellerAddr string  `json:"sellAddr"`  //eg:56320f
	Volume     float64 `json:"volume"`    //eg:111.11 DICE
	Price      float64 `json:"price"`     //eg:0.00245
	OrderTime  string  `json:"orderTime"` //eg:2018-11-10 21:04:01
	Unit       string  `json:"unit"`      //单位
	OrderType  int     `json:"orderType"` //订单类型：0.买单 1，卖单
	PairID     int64   `json:"pairID"`    //pairID
}

//分页查询参数
type PageInfo struct {
	Start int `form:"start,default=0"`
	Limit int `form:"limit,default=20"`
}
type UserOrderListReq struct {
	Start      int    `form:"start,default=0"`
	Limit      int    `form:"limit,default=20"`
	UAddr      string `form:"userAddr,default=10" binding:"required"` //当前用户的address
	FTokenAddr string `form:"tokenAddr1"`                             //第一个token地址
	STokenAddr string `form:"tokenAddr2"`                             //第二个token地址
	Status     string `form:"status"`                                 //eg:1   1,2,3
}
type MarketHisReq struct {
	Start  int   `form:"start,default=0"`
	Limit  int   `form:"limit,default=20"`
	PairID int64 `form:"pairID" binding:"required"` //交易对ID
}
type MarketPairListInfo struct {
	Total int64                `json:"total"`
	Rows  []*MarketPairRowInfo `json:"rows"`
}

type UserOrderListInfo struct {
	Total int64               `json:"total"`
	Rows  []*UserOrderRowInfo `json:"rows"`
}

type UserOrderListInfoDetail struct {
	Total int64                     `json:"total"`
	Rows  []*UserOrderRowInfoDetail `json:"rows"`
}
type LatestOrderListInfo struct {
	Total int64              `json:"total"`
	Rows  []*OrderHisRowInfo `json:"rows"`
}

type CommonOrder struct {
	OrderID          int64   `json:"orderID"`          //订单ID，智能合约挂单后生成的唯一标示
	BlockID          int64   `json:"blockID"`          //区块ID
	ExchangeID       int64   `json:"exchangeID"`       //所属交易对ID
	OwnerAddress     string  `json:"ownerAddress"`     //挂单账户地址
	BsFlag           int32   `json:"bsFlag"`           //买卖标记  1 买入  2 卖出
	Amount           float64 `json:"amount"`           //挂单总量
	BuyTokenAddress  string  `json:"buyTokenAddress"`  //买token地址
	BuyTokenAmount   float64 `json:"buyTokenAmount"`   //买token数量
	SellTokenAddress string  `json:"sellTokenAddress"` //卖token地址
	SellTokenAmount  float64 `json:"sellTokenAmount"`  //卖token数量
	Price            float64 `json:"price"`            //订单价格
	TrxHash          string  `json:"trxHash"`          //订单hash
	OrderTime        string  `json:"orderTime"`        //订单时间，所属区块时间
	IsCancel         bool    `json:"isCancel"`         //是否撤单 初始为false，撤单为true
	Status           int64   `json:"status"`           // 订单状态
	CurTurnover      float64 `json:"curTurnover"`      //成交额
}

type OrderListResp struct {
	Buy   []*CommonOrder `json:"buy"`   // 买订单，按价格从高到低排列
	Sell  []*CommonOrder `json:"sell"`  // 卖订单 按价格从低到高排列
	Price string         `json:"price"` //当前成交价格
}

//TopPrice 返回盘口买单最高价和卖单最低价
type TopPrice struct {
	ExchangeID   int64   `json:"exchangeID"`   // 交易对ID
	BuyHighPrice float64 `json:"buyHighPrice"` // 买单最高价
	SellLowPrice float64 `json:"sellLowPrice"` // 卖单最低价
	Price        float64 `json:"price"`        //当前成交价格
}

// 添加交易对请求参数entity
type AddChannelIdReq struct {
	Hash      string `json:"hash" form:"hash"`           //订单hash
	OrderId   string `json:"orderId"  form:"orderId"`    //订单id
	ChannelId string `json:"channelId" form:"channelId"` //订单渠道ID
}
