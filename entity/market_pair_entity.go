package entity

type MarketPair struct {
	ID              int64   `json:"id"`
	FirstTokenName  string  `json:"firstTokenName"`  //eg:dice's token name
	FirstTokenAddr  string  `json:"firstTokenAddr"`  //eg:dice's token ContractAddress
	SecondTokenName string  `json:"secondTokenName"` //eg:the token name of trx
	SecondTokenAddr string  `json:"secondTokenAddr"` //eg:the token ContractAddress of trx
	Price           float64 `json:"price"`           //eg: 0.000245
	Unit            string  `json:"unit"`            //eg:DICE
	FirstPrecision  int64   //第一个token精度
	SecondPrecision int64   //第二个token精度
	FisrtShortName  string  //第一个token简称
	SecondShortName string  //第二个token简称
	CreatedAt       string
	UpdateAt        string
	PairType        int64
	DefaultIdx      int64 // 可设置的索引
}

//交易对24h交易信息
type PriceInf struct {
	HighestPrice24h float64 `json:"highestPrice24h"` //24h最高价
	LowestPrice24h  float64 `json:"lowestPrice24h"`  //24h最低价
	Volume24h       float64 `json:"volume24h"`       //24h成交量
	Amount24h       float64 `json:"amount24h"`       //24h成交量
}

type OrderRecord struct {
	Schedule    string `json:"schedule"`    //eg:23.2
	CurTurnover string `json:"curTurnover"` //成交额
}

type MarketPairInfo struct {
	FirstTokenName  string `json:"firstTokenName"`  //eg:dice's token name
	SecondTokenName string `json:"secondTokenName"` //eg:the token name of trx
	FirstPrecision  int64  //第一个token精度
	SecondPrecision int64  //第二个token精度
}
