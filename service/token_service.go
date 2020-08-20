package service

import (
	"github.com/wlcy/tradehome-service/buffer"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/module"
)

func QueryTokenInfo(address string) (*entity.TokenInfo, error) {
	return module.QueryTokenInfo(address)
}

//QueryMarketsBuffer ... 从buffer获取市场信息
func QueryMarketsBuffer() (*entity.MarketPriceResp, error) {

	marketBuffer := buffer.GetMarketBuffer()
	markets := marketBuffer.GetMarket()

	return markets, nil
}
