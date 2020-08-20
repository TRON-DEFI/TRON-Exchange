package entity

import (
	"fmt"
	"github.com/wlcy/tradehome-engin/lib/log"
	"github.com/wlcy/tradehome-engin/lib/mysql"
	"github.com/wlcy/tradehome-engin/lib/util"
	"github.com/wlcy/tradehome-engin/web/common"
	"strconv"
	"time"
)

//ExchangeTransactionContract Exchange Transaction Contract
type ExchangeTransactionContract struct {
	OwnerAddress string `json:"owner_address,omitempty"`
	ExchangeID   int64  `json:"exchange_id,omitempty"`
	TokenID      string `json:"token_id,omitempty"`
	Quant        int64  `json:"quant,omitempty"`
	Expected     int64  `json:"expected,omitempty"`
}

//ExchangeCreateContract ...
type ExchangeCreateContract struct {
	OwnerAddress       string `json:"owner_address"`
	FirstTokenID       string `json:"first_token_id"`
	FirstTokenBalance  int64  `json:"first_token_balance"`
	SecondTokenID      string `json:"second_token_id"`
	SecondTokenBalance int64  `json:"second_token_balance"`
	Transaction        string `json:"transaction"`
}

//ExchangePairsDetail ...
type ExchangePairsDetail struct {
	TotalPairs    int32                 `json:"total,omitempty"`
	ExchangePairs []*ExchangePairDetail `json:"data,omitempty"`
}

// ExchangePairDetail ...
type ExchangePairDetail struct {
	ExchangeID         int64   `json:"exchange_id,omitempty"`
	CreatorAddress     string  `json:"creator_address,omitempty"`
	CreateTime         int64   `json:"create_time,omitempty"`
	FirstTokenID       string  `json:"first_token_id,omitempty"`
	FirstTokenBalance  int64   `json:"first_token_balance,omitempty"`
	SecondTokenID      string  `json:"second_token_id,omitempty"`
	SecondTokenBalance int64   `json:"second_token_balance,omitempty"`
	ExchangeName       string  `json:"exchange_name,omitempty"`
	Price              float64 `json:"price,omitempty"`
	Volume             int64   `json:"volume,omitempty"`
	UpDownPercent      string  `json:"up_down_percent,omitempty"`
	High               float64 `json:"high,omitempty"`
	Low                float64 `json:"low,omitempty"`
}

// KChartQueryParam ...
type KChartQueryParam struct {
	ExchangeID string `json:"exchange_id"`
	TimeStart  string `json:"time_start"`
	TimeEnd    string `json:"time_end"`
	Granu      string `json:"granularity"`
}

//KChartResp ...
type KChartResp struct {
	ExchangeID string        `json:"exchangeId"`
	TimeStart  string        `json:"startTime"`
	TimeEnd    string        `json:"endTime"`
	Granu      string        `json:"gran"`
	Data       []*KChartData `json:"data"`
}

//KChartData ...
type KChartData struct {
	Time   string `json:"time"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

//ExchangeOper 交易对注资和撤资
type ExchangeOper struct {
	OwnerAddress string `json:"owner_address"`
	ExchangeID   int64  `json:"exchange_id"`
	TokenID      string `json:"token_id"`
	Quant        int64  `json:"quant"`
	Transaction  string `json:"transaction"`
}

//ExchangeAuthReq ...
type ExchangeAuthReq struct {
	Sort    string `json:"sort,omitempty"`    // 按照区块高度倒序
	Limit   int64  `json:"limit,omitempty"`   // 每页记录数
	Start   int64  `json:"start,omitempty"`   // 记录的起始序号
	Address string `json:"address,omitempty"` // 按照地址精确查询
}

//ExchangeAuthResp ...
type ExchangeAuthResp struct {
	Total int64               `json:"total"`
	Data  []*ExchangeAuthInfo `json:"data"`
}

//ExchangeAuthInfo ...
type ExchangeAuthInfo struct {
	OwnerAddress  string `json:"owner_address"`
	FirstTokenID  string `json:"first_token_id"`
	SecondTokenID string `json:"second_token_id"`
	CreateTime    string `json:"create_time"`
}

//PostExchangeAuthResp ...
type PostExchangeAuthResp struct {
	Data string `json:"data"`
}

//ExchangeListReq ...
type ExchangeListReq struct {
	Sort       string `json:"sort,omitempty"`       // 按照区块高度倒序
	Limit      int64  `json:"limit,omitempty"`      // 每页记录数
	Start      int64  `json:"start,omitempty"`      // 记录的起始序号
	Address    string `json:"address,omitempty"`    // 按照交易人精确查询
	Name       string `json:"name,omitempty"`       //交易对名称查询
	ExchangeID string `json:"exchangeID,omitempty"` //交易对ID
}

//ExchangeListResp ...
type ExchangeListResp struct {
	Total int64                      `json:"total"`
	Data  []*ExchangeTransactionInfo `json:"data"`
}

//ExchangeTransactionInfo ...
type ExchangeTransactionInfo struct {
	ExchangeID     string  `json:"exchangeID"`     // :"1",//交易对ID
	BlockID        int64   `json:"blockID"`        // :"1",//区块ID
	CreatorAddress string  `json:"creatorAddress"` // :"TFA1qpUkQ1yBDw4pgZKx25wEZAqkjGoZo1",//创建者
	CreatorName    string  `json:"creatorName"`    // :"Sesamesee",//创建者
	TrxHash        string  `json:"trx_hash"`       // :"0afa11cbfa9b4707b1308addc48ea31201157a989db92fe75750c068f0cc14e0",//交易hash
	CreateTime     int64   `json:"createTime"`     // :"1536416859000",//创建时间
	TokenID        string  `json:"tokenID"`        // :"IGG/MEETONE",//交易对名称
	Quant          float64 `json:"quant"`          // :23,//对第一个token的交易数量
	Confirmed      bool    `json:"confirmed"`      // :0,//0 未确认，1 已确认
}

//ExchangeReportResp ...
type ExchangeReportResp struct {
	Total int64             `json:"total"`
	Data  []*ExchangeReport `json:"data"`
}

//ExchangeReport 首页统计交易对及其24h交易量，涨幅
type ExchangeReport struct {
	ExchangeID         int64   `json:"exchange_id"`          // :"1",//交易对ID
	ExchangeName       string  `json:"exchange_name"`        // :"IGG/MEETONE",//交易对名称
	CreatorAddress     string  `json:"creator_address"`      //:"TFA1qpUkQ1yBDw4pgZKx25wEZAqkjGoZo1",//创建者
	FirstTokenID       string  `json:"first_token_id"`       //:"IGG",//第一个tokenID
	FirstTokenBalance  float64 `json:"first_token_balance"`  //:10000,//第一个token 余额
	SecondTokenID      string  `json:"second_token_id"`      //:"TRX",//第二个tokenID
	SecondTokenBalance float64 `json:"second_token_balance"` //:200,//第二个token 余额
	CreateTime         int64   `json:"create_time"`          //:"1536416859000",//创建时间
	Price              string  `json:"price"`                //:0.0023,//交易对价格
	Volume             string  `json:"volume"`               //:3345.342,//24H成交量
	SVolume            string  `json:"svolume"`              //:3345.342,//第二个token 24H成交量
	UpDownPercent      string  `json:"up_down_percent"`      //:"-6.64%",//涨幅
	High               string  `json:"high"`                 //:0.0025,//最高价格
	Low                string  `json:"low"`                  //:0.002,//最低价格
}

//ExchangeCalc 交易对预估
type ExchangeCalc struct {
	ExchangeID    int64 `json:"exchangeID"`
	BuyTokenQuant int64 `json:"buyTokenQuant"`
}

//ExchangeTransactionDetail ...
type ExchangeTransactionDetail struct {
	ExchangeID       int64   `json:"exchangeID"`         // :"1",//交易对ID
	BlockID          int64   `json:"blockID"`            // :"1",//区块ID
	CreatorAddress   string  `json:"creatorAddress"`     // :"TFA1qpUkQ1yBDw4pgZKx25wEZAqkjGoZo1",//创建者
	TrxHash          string  `json:"trx_hash"`           // :"0afa11cbfa9b4707b1308addc48ea31201157a989db92fe75750c068f0cc14e0",//交易hash
	CreateTime       int64   `json:"createTime"`         // :"1536416859000",//创建时间
	FirstTokenID     string  `json:"first_token_id"`     //:"IGG",//第一个tokenID
	FirstTokenQuant  float64 `json:"first_token_quant"`  //:10000,//第一个token 交易量
	SecondTokenID    string  `json:"second_token_id"`    //:"TRX",//第二个tokenID
	SecondTokenQuant float64 `json:"second_token_quant"` //:200,//第二个token 交易量
	Price            string  `json:"price"`              //:价格
	Confirmed        int64   `json:"confirmed"`          // :0,//0 未确认，1 已确认
}

//Insert ...
func (w *ExchangeTransactionDetail) Insert() error {
	if nil == w || "" == w.TrxHash {
		return util.NewErrorMsg(util.Error_common_no_data)
	}

	f64, err := strconv.ParseFloat(w.Price, 64)
	if nil != err {
		log.Errorf("Price is not legal: [%v]", w.Price)
		return err
	}
	//multi 1E8 when store into DB
	price := int64(f64 * float64(common.ExchangeFactor))
	strSQL := fmt.Sprintf(`insert into tron_ext.exchange_transaction_detail(trx_hash, create_time, confirmed, owner_address, exchange_id, first_token_id, first_quant, second_token_id, second_quant, price)
	 values('%v', '%v','%v', '%v', '%v', '%v','%v', '%v', '%v', '%v')`,
		w.TrxHash,
		time.Now().UTC().UnixNano()/1000000,
		w.Confirmed,
		w.CreatorAddress,
		w.ExchangeID,
		w.FirstTokenID,
		w.FirstTokenQuant,
		w.SecondTokenID,
		w.SecondTokenQuant,
		price)
	insertID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("Insert exchange transaction fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("Insert exchange transaction success, transaction hash: [%v]", insertID)
	return nil
}

type ChannelId struct {
	ID        int64  `json:"id"`
	Hash      string `json:"hash"`
	OrderId   string `json:"orderId"`
	ChannelId string `json:"channelId"`
	CreatedAt string
}

//ExchangeInfo 交易对精度信息
type ExchangeInfo struct {
	ExchangeID      int64 `json:"exchangeID"`       //
	FirstPrecision  int64 `json:"first_precision"`  //
	SecondPrecision int64 `json:"second_precision"` //
}
