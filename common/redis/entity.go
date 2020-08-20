package redis

//OrderListRedisKey 盘口订单信息，存储买卖队列中的订单id，ownerAddress, 挂单价格，挂单数量 val:map[exchange]OrderListRedis，不设有效期
var OrderListRedisKey = "trxmarket:orderlist"
var OrderListRedisKey10 = "trxmarket10:orderlist"

//ExchangeIDRealPriceRedisKey 保存所有交易对信息的实时价格，交易总量，0点价格 val:OrderRealPrice 不舍有效期，撤单成功后，删除key
var ExchangeIDRealPriceRedisKey = "trxmarket:exchangeprice:exchangeID:%v"

//OrderBlockIDOffset 事件服务器中order事件中的最新区块信息
var OrderBlockIDOffset = "trxmarket:order:newblockIDoffset"

//OrderListRedis 存入redis中的盘口信息
type OrderListRedis struct {
	Buy   []*CommonOrder `json:"buy"`   // 买订单，按价格从高到低排列
	Sell  []*CommonOrder `json:"sell"`  // 卖订单 按价格从低到高排列
	Price string         `json:"price"` //当前成交价格
}

//CommonOrder 订单通用结构
type CommonOrder struct {
	OrderID       int64  `json:"orderID"`       //订单ID，智能合约挂单后生成的唯一标示
	BlockID       int64  `json:"blockID"`       //区块ID
	ExchangeID    int64  `json:"exchangeID"`    //所属交易对ID
	OwnerAddress  string `json:"ownerAddress"`  //挂单账户地址
	BsFlag        int32  `json:"bsFlag"`        //买卖标记  1 买入  2 卖出
	TokenAddressA string `json:"tokenAddress1"` //第一个token地址
	TokenAmountA  int64  `json:"tokenAmount1"`  //第一个token数量
	TokenAddressB string `json:"tokenAddress2"` //第二个token地址
	TokenAmountB  int64  `json:"tokenAmount2"`  //第二个token数量
	OrderFill     int64  `json:"orderFill"`     //orderfill
	TurnOver      int64  `json:"turnOver"`      //turnover
	AmountA       int64  `json:"amount1"`       //剩余第一个token数量 = TokenAmount1 - OrderFill
	Price         int64  `json:"price"`         //订单价格
	TrxHash       string `json:"trxHash"`       //订单hash
	OrderTime     string `json:"orderTime"`     //订单时间，所属区块时间
	IsCancel      bool   `json:"isCancel"`      //是否撤单  初始为false，撤单为true 暂时不用
	Status        int64  `json:"status"`        // 订单状态
}

//ExchangeIDRealPrice 交易对ID实时信息 实时价格，交易总量，0点价格
type ExchangeIDRealPrice struct {
	ExchangeID int64  `json:"exchangeID"` //交易对ID
	RealPrice  string `json:"realPrice"`  // 实时价格-更新
	ZeroPrice  string `json:"zeroPrice"`  // 零点价格-更新
	Amount     int64  `json:"amount"`     // 交易总量-累加更新
}
