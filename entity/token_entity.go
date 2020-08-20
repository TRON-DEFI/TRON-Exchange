package entity

type TokenInfo struct {
	Id          int64  `json:"id"`
	Address     string `json:"address"`
	FullName    string `json:"fullName"`
	ShortName   string `json:"shortName"`
	Circulation int64  `json:"circulation"`
	Precision   int    `json:"precision"`
	Description string `json:"description"`
	WebsiteUrl  string `json:"websiteUrl"`
	LogoUrl     string `json:"logoUrl"`
}

type TokenInfoReq struct {
	Start   int64 `json:"start"`
	Limit   int64 `json:"limit"`
	Address int64 `json:"Address"`
}

//MarketInfo 交易所信息
type MarketInfo struct {
	Rank             int64   `json:"rank"`             //:1,
	Name             string  `json:"name"`             //:"Rfinex",
	Pair             string  `json:"pair"`             //:"TRX/ETH",
	Link             string  `json:"link"`             //:"https://rfinex.com/",
	Volume           float64 `json:"volume"`           //:22144662.8099,
	VolumePercentage float64 `json:"volumePercentage"` //:19.6793615403,
	VolumeNative     float64 `json:"volumeNative"`     //:1194868733.76,
	Price            float64 `json:"price"`            //:0.0185331343806
}

type MarketCapResp struct {
	Status *CapStatus `json:"status"` //:1,
	Data   []*CapData `json:"data"`   //:1,
}
type CapStatus struct {
	Timestamp    string `json:"timestamp"`     //: "2020-08-15T15:57:52.477Z",
	ErrorCode    int64  `json:"error_code"`    //: 0,
	ErrorMessage string `json:"error_message"` //: null,
	Elapsed      int64  `json:"elapsed"`       //: 5,
	CreditCount  int64  `json:"credit_count"`  //: 0,
	Notice       string `json:"notice"`        //: null
}
type CapData struct {
	ID                int64                  `json:"id"`                 //: 3635,
	Name              string                 `json:"name"`               //: "Crypto.com Coin",
	Symbol            string                 `json:"symbol"`             //: "CRO",
	Slug              string                 `json:"slug"`               //: "crypto-com-coin",
	NumMarketPairs    int64                  `json:"num_market_pairs"`   //: 55,
	DateAdded         string                 `json:"date_added"`         //: "2018-12-14T00:00:00.000Z",
	Tags              []string               `json:"tags"`               //: [],
	MaxSupply         string                 `json:"max_supply"`         //: null,
	CirculatingSupply float64                `json:"circulating_supply"` //: 18860730593.6073,
	TotalSupply       float64                `json:"total_supply"`       //: 100000000000,
	Platform          *PlatForm              `json:"platform"`           //: {},
	CmcRank           int64                  `json:"cmc_rank"`           //: 12,
	LastUpdated       string                 `json:"last_updated"`       //: "2020-08-15T15:56:12.000Z",
	Quote             map[string]*QuoteValue `json:"quote"`              //: map{
}

type PlatForm struct {
	ID           int64  `json:"id"`            //: 1027,
	Name         string `json:"name"`          //: "Ethereum",
	Symbol       string `json:"symbol"`        //: "ETH",
	Slug         string `json:"slug"`          //: "ethereum",
	TokenAddress string `json:"token_address"` //: "0xa0b73e1ff0b80914ab6fe0444e65848c4c34450b"
}
type QuoteValue struct {
	Price            float64 `json:"price"`              //: 0.166932501693,
	Volume24h        float64 `json:"volume_24h"`         //: 75381057.3595886,
	PercentChange1h  float64 `json:"percent_change_1h"`  //: 0.665583,
	PercentChange24h float64 `json:"percent_change_24h"` //: 1.53448,
	PercentChange7d  float64 `json:"percent_change_7d"`  //: 0.187292,
	MarketCap        float64 `json:"market_cap"`         //: 3148468941.7485676,
	LastUpdated      string  `json:"last_updated"`       //: "2020-08-15T15:56:12.000Z"
}

type MarketPriceResp struct {
	Price   float64 `json:"price"`
	Updated string  `json:"updated"`
}
