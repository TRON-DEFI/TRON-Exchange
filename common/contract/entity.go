package contract

import (
	"container/list"
	"github.com/tronprotocol/grpc-gateway/api"
	"sync"
)

const (
	BUY    = 0 //买
	SELL   = 1 //卖
	CANCEL = 2 //撤销
	DONE   = 3 //成交
)

//CallRecord 智能合约调用通用结构
type CallRecord struct {
	Owner     string // base58
	Contract  string // base58
	Method    string
	CallValue int64
	Data      []byte
	TrxHash   string // hex
	Err       error
	Return    *api.Return
}

//CommonOrder 订单通用结构
type CommonOrder struct {
	OrderID       int64  //订单ID，智能合约挂单后生成的唯一标示
	BlockID       int64  //区块ID
	ExchangeID    int64  //所属交易对ID
	OwnerAddress  string //挂单账户地址
	BsFlag        int32  //买卖标记  1 买入  2 卖出
	TokenAddressA string //第一个token地址
	TokenAmountA  int64  //第一个token数量
	TokenAddressB string //第二个token地址
	TokenAmountB  int64  //第二个token数量
	OrderFill     int64  //orderfill
	TurnOver      int64  //turnover
	AmountA       int64  //剩余第一个token数量 = TokenAmountA - OrderFill
	Price         int64  //订单价格
	TrxHash       string //订单hash
	OrderTime     string //订单时间，所属区块时间
	IsCancle      bool   //是否撤单  初始为false，撤单为true 暂时不用
	Status        int64  // 订单状态
}

//MatchList 买入卖出队列
type MatchList struct {
	BuyList  list.List //sorted list asc
	SellList list.List //sorted list desc
}

//TrxMatchList 盘口订单结构
type TrxMatchList struct {
	List         MatchList
	exchangeID   int64                   //队列所属交易对，用数据库主键ID标示
	nodePosition map[int64]*list.Element //记录每个节点地址，用于删除操作
	l            sync.RWMutex
}

//===================================持久层结构=================================================

//OrderBase 订单基础信息 not use
type OrderBase struct {
	OrderID      int64  //订单ID，智能合约挂单后生成的唯一标示
	OwnerAddress string //挂单账户地址
	Price        string //订单价格
	Amount       int64  //挂单数量
}

//OrderListRedis 存入redis中的盘口信息
type OrderListRedis struct {
	Buy   []*CommonOrder `json:"buy"`   // 买订单，按价格从高到低排列
	Sell  []*CommonOrder `json:"sell"`  // 卖订单 按价格从低到高排列
	Price string         `json:"price"` //当前成交价格
}

//ExchangeIDRealPrice 交易对ID实时信息 实时价格，交易总量，0点价格
type ExchangeIDRealPrice struct {
	ExchangeID int64  `json:"exchangeID"` //交易对ID
	RealPrice  string `json:"realPrice"`  // 实时价格-更新
	ZeroPrice  string `json:"zeroPrice"`  // 零点价格-更新
	Amount     int64  `json:"amount"`     // 交易总量-累加更新
}

//OrderCancelEventLog eventLog中的cancel事件
type OrderCancelEventLog struct {
	BlockNumber     int64            `json:"block_number"`     //: 12811,
	BlockTimestamp  int64            `json:"block_timestamp"`  //: 1542694011000,
	ContractAddress string           `json:"contract_address"` //: "TMREyhsHfvMZp1xqgZjZM6Xsq4BW6hrjDf",
	EventIndex      int64            `json:"event_index"`      //: 0,
	EventName       string           `json:"event_name"`       //: "Order",
	Result          *CancelEvent     `json:"result"`           //:
	ResultType      *CancelEventType `json:"result_type"`      //:
	TransactionID   string           `json:"transaction_id"`   //: "8bacb4b3a6d7c2d034d5833f650ed45e1e54ef09f2c07c0896effeb7ceef2d3d",
	ResourceNode    string           `json:"resource_Node"`    //: "FullNode"
}
type CancelEvent struct {
	OrderID     string `json:"orderID"`     //: "1",
	CreateTime  string `json:"createTime"`  //: "1542694011",
	BlockNumber string `json:"blockNumber"` //: "12811",
	Status      string `json:"status"`      //: "0"
}

type CancelEventType struct {
	OrderID     string `json:"orderID"`     //: "uint256",
	CreateTime  string `json:"createTime"`  //: "uint256",
	BlockNumber string `json:"blockNumber"` //: "uint256",
	Status      string `json:"status"`      //: "uint256"
}

//OrderNewEventLog  自己event服务接口
type OrderNewEventLog struct {
	Code    int64            `json:"code"`    //
	Message string           `json:"message"` //
	Data    []*OrderEventLog `json:"data"`    //
}

//OrderEventLog eventLog中的order事件
type OrderEventLog struct {
	BlockNumber     int64            `json:"block_number"`     //: 12811,
	BlockTimestamp  int64            `json:"block_timestamp"`  //: 1542694011000,
	ContractAddress string           `json:"contract_address"` //: "TMREyhsHfvMZp1xqgZjZM6Xsq4BW6hrjDf",
	EventIndex      int64            `json:"event_index"`      //: 0,
	EventName       string           `json:"event_name"`       //: "Order",
	Result          *EventResult     `json:"result"`           //:
	ResultType      *EventResultType `json:"result_type"`      //:
	TransactionID   string           `json:"transaction_id"`   //: "8bacb4b3a6d7c2d034d5833f650ed45e1e54ef09f2c07c0896effeb7ceef2d3d",
	ResourceNode    string           `json:"resource_Node"`    //: "FullNode"
}

//EventResult 事件结果
type EventResult struct {
	OrderType string `json:"orderType"`    //: "0",
	OrderID   string `json:"orderID"`      //: "1",
	AmountB   string `json:"secondAmount"` //: amountB "111000000",
	//CreateTime  string `json:"createTime"`  //: "1542694011",
	AmountA string `json:"firstAmount"` //amountA: "111",
	Price   string `json:"price"`       //: "1000000",
	//BlockNumber string `json:"blockNumber"` //: "12811",
	TokenA string `json:"firstToken"`  //tokenA: "0x3866431166cd11e6baf59d15aac2acb8515c63a8",
	TokenB string `json:"secondToken"` //:tokenB "0x0000000000000000000000000000000000000000",
	User   string `json:"user"`        //: "0xe6e81042606060d9a58ce19fcde06e3fc86a1d00",
	Status string `json:"orderStatus"` //: "0"status
}

//EventResultType 事件结果类型
type EventResultType struct {
	OrderType string `json:"orderType"`    //: "0",
	OrderID   string `json:"orderID"`      //: "1",
	AmountB   string `json:"secondAmount"` //: amountB "111000000",
	//CreateTime  string `json:"createTime"`  //: "1542694011",
	AmountA string `json:"firstAmount"` //amountA: "111",
	Price   string `json:"price"`       //: "1000000",
	//BlockNumber string `json:"blockNumber"` //: "12811",
	TokenA string `json:"firstToken"`  //tokenA: "0x3866431166cd11e6baf59d15aac2acb8515c63a8",
	TokenB string `json:"secondToken"` //:tokenB "0x0000000000000000000000000000000000000000",
	User   string `json:"user"`        //: "0xe6e81042606060d9a58ce19fcde06e3fc86a1d00",
	Status string `json:"orderStatus"` //: "0"status
}

//OrderNewTransactionLog  自己event服务接口
type OrderNewTransactionLog struct {
	Code    int64                  `json:"code"`    //
	Message string                 `json:"message"` //
	Data    []*OrderTransactionLog `json:"data"`    //
}

//==========================================
//OrderTransactionLog transaction中的返回结果，包括 Completed， Trade事件
type OrderTransactionLog struct {
	BlockNumber     int64                  `json:"block_number"`     //: 12811,
	BlockTimestamp  int64                  `json:"block_timestamp"`  //: 1542694011000,
	ContractAddress string                 `json:"contract_address"` //: "TMREyhsHfvMZp1xqgZjZM6Xsq4BW6hrjDf",
	EventIndex      int64                  `json:"event_index"`      //: 0,
	EventName       string                 `json:"event_name"`       //: "Order",
	Result          map[string]interface{} `json:"result"`           //:
	ResultType      map[string]interface{} `json:"result_type"`      //:
	TransactionID   string                 `json:"transaction_id"`   //: "8bacb4b3a6d7c2d034d5833f650ed45e1e54ef09f2c07c0896effeb7ceef2d3d",
	ResourceNode    string                 `json:"resource_Node"`    //: "FullNode"
}

//TransactionResult 事件结果
type TransactionResult struct {
	FeeAmountA  string `json:"feeAmountA"`  //: "99500000000",
	OrderFillsB string `json:"orderFillsB"` //: "0",
	BuyOrderID  string `json:"buyOrderId"`  //: "1",
	FeeAmountB  string `json:"feeAmountB"`  //: "99500000",
	AmountA     string `json:"amountA"`     //: "100000000000",
	Price       string `json:"price"`       //: "100000000",
	TradeStatus string `json:"tradeStatus"` //: "1",
	OrderFillsA string `json:"orderFillsA"` // "0",
	SellOrderID string `json:"sellOrderId"` //: "2"
}

type OrderNewCompleteEventLog struct {
	Code    int64                    `json:"code"`    //
	Message string                   `json:"message"` //
	Data    []*OrderCompleteEventLog `json:"data"`    //
}

//OrderCompleteEventLog eventLog中的order事件
type OrderCompleteEventLog struct {
	BlockNumber     int64                    `json:"block_number"`     //: 12811,
	BlockTimestamp  int64                    `json:"block_timestamp"`  //: 1542694011000,
	ContractAddress string                   `json:"contract_address"` //: "TMREyhsHfvMZp1xqgZjZM6Xsq4BW6hrjDf",
	EventIndex      int64                    `json:"event_index"`      //: 0,
	EventName       string                   `json:"event_name"`       //: "Order",
	Result          *TransactionComplete     `json:"result"`           //:
	ResultType      *TransactionCompleteType `json:"result_type"`      //:
	TransactionID   string                   `json:"transaction_id"`   //: "8bacb4b3a6d7c2d034d5833f650ed45e1e54ef09f2c07c0896effeb7ceef2d3d",
	ResourceNode    string                   `json:"resource_Node"`    //: "FullNode"
}

//TransactionComplete ...
type TransactionComplete struct {
	OrderID string `json:"orderID"`     //: "1",  Completed 事件
	Status  string `json:"orderStatus"` //: "0"   Completed 事件
}

//TransactionCompleteType ...
type TransactionCompleteType struct {
	OrderID string `json:"orderID"`     //: "1",  Completed 事件
	Status  string `json:"orderStatus"` //: "0"   Completed 事件
}

//TransactionResultType 事件结果类型
type TransactionResultType struct {
	FeeAmountA  string `json:"feeAmountA"`  //: "99500000000",
	OrderFillsB string `json:"orderFillsB"` //: "0",
	BuyOrderID  string `json:"buyOrderId"`  //: "1",
	FeeAmountB  string `json:"feeAmountB"`  //: "99500000",
	AmountA     string `json:"amountA"`     //: "100000000000",
	Price       string `json:"price"`       //: "100000000",
	TradeStatus string `json:"tradeStatus"` //: "1",
	OrderFillsA string `json:"orderFillsA"` // "0",
	SellOrderID string `json:"sellOrderId"` //: "2"
}

//SmartOrderInfo 合约中查询到的订单信息
type SmartOrderInfo struct {
	OrderID      int64  `json:"orderID"`    //: "99500000000",
	Status       int64  `json:"status"`     //: "0",
	OrderType    int64  `json:"orderType"`  //: "1",
	OwnerAddress string `json:"user"`       //: "99500000",
	TokenA       string `json:"token1"`     //: "100000000000",
	TokenB       string `json:"token2"`     //: "100000000",
	AmountA      int64  `json:"amount1"`    //: "1",
	AmountB      int64  `json:"amount2"`    // "0",
	OrderFills   int64  `json:"orderFills"` //: "2"
	Turnover     int64  `json:"turnover"`   //: "2"
	ChannelId    int64  `json:"channelId"`
}

//SmartOrderInfo 合约中查询到的订单信息
type SmartOrderInfo10 struct {
	OrderID      int64  `json:"orderID"`    //: "99500000000",
	Status       int64  `json:"status"`     //: "0",
	OrderType    int64  `json:"orderType"`  //: "1",
	OwnerAddress string `json:"user"`       //: "99500000",
	TokenA       string `json:"tokenA"`     //: "100000000000",
	TokenB       string `json:"tokenB"`     //: "100000000",  20token交易所有该字段，没有Price
	Price        int64  `json:"price"`      //: "100000000",	10token交易所有该字段，没有tokenB
	AmountA      int64  `json:"amountA"`    //: "1",
	AmountB      int64  `json:"amountB"`    // "0",
	OrderFills   int64  `json:"orderFills"` //: "2"
	Turnover     int64  `json:"turnover"`   //: "2"
}

//ExchangeInfo 交易对精度信息
type ExchangeInfo struct {
	ExchangeID      int64 `json:"exchangeID"`       //
	FirstPrecision  int64 `json:"first_precision"`  //
	SecondPrecision int64 `json:"second_precision"` //
}

type SmartOrderInfoTwo struct {
	OrderID      int64   `json:"orderID"`    //: "99500000000",
	Status       int64   `json:"status"`     //: "0",
	OrderType    int64   `json:"orderType"`  //: "1",
	OwnerAddress string  `json:"user"`       //: "99500000",
	TokenA       string  `json:"token1"`     //: "100000000000",
	TokenB       string  `json:"token2"`     //: "100000000",
	AmountA      float64 `json:"amount1"`    //: "1",
	AmountB      float64 `json:"amount2"`    // "0",
	OrderFills   float64 `json:"orderFills"` //: "2"
	Turnover     float64 `json:"turnover"`   //: "2"
	ChannelId    int64   `json:"channelId"`
}
