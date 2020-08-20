package buffer

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/util"
	"sync"
	"time"
)

/*
store all marketBuffer data in memory
load from db every 30 seconds
*/

var _marketBuffer *marketBuffer
var onceMarketOnce sync.Once

//GetMarketBuffer ...
func GetMarketBuffer() *marketBuffer {
	return getMarketBuffer()
}

// getMarketBuffer
func getMarketBuffer() *marketBuffer {
	onceMarketOnce.Do(func() {
		_marketBuffer = &marketBuffer{}
		_marketBuffer.load()

		go marketBufferLoader()
	})
	return _marketBuffer
}

func marketBufferLoader() {
	for {
		_marketBuffer.load()
		time.Sleep(5 * time.Second)
	}
}

type marketBuffer struct {
	sync.RWMutex

	marketInfoList []*entity.MarketInfo
	price          float64

	updateTime string
}

func (w *marketBuffer) GetMarket() *entity.MarketPriceResp {
	resp := &entity.MarketPriceResp{}
	w.RLock()
	resp.Price = w.price
	resp.Updated = w.updateTime
	w.RUnlock()
	return resp
}

func (w *marketBuffer) load() {
	marketURL := "https://web-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=6&start=12&convert=USD#"
	resp, err := util.SendRequest(marketURL, "GET", "", nil)
	if err != nil {
		logrus.Error(err)
		return
	}
	var market entity.MarketCapResp
	err = json.Unmarshal(resp.Bytes(), &market)
	if err != nil {
		logrus.Error(err)
		return
	}
	if &market == nil || market.Status.ErrorCode != 0 {
		logrus.Error("request coinmarket api err!")
		return
	}
	var price float64
	for _, data := range market.Data {
		if data.Symbol == "TRX" {
			if _, ok := data.Quote["USD"]; ok {
				price = data.Quote["USD"].Price
			}
			break
		}
	}

	logrus.Infof("market price in buffer done.")
	w.Lock()
	w.price = price
	w.updateTime = time.Now().Local().Format(util.DATETIMEFORMAT)
	w.Unlock()
}

//
//func (w *marketBuffer) load() {
//	marketInfos := make([]*entity.MarketInfo, 0)
//	marketURL := "https://coinmarketcap.com/currencies/tron/"
//	_, body, errs := gorequest.New().Get(marketURL).End()
//	if errs != nil && len(errs) > 0 {
//		logrus.Error(errs)
//		return
//	}
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
//	if err != nil {
//		logrus.Error(err)
//		return
//	}
//	doc.Find("#__next > tbody > tr").Each(func(i int, s *goquery.Selection) {
//		marketInfo := &entity.MarketInfo{}
//		node := strconv.Itoa(i + 1)
//		rank, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(1)").Html()
//		name, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(2)").Attr("data-sort")
//		pair, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(3)").Attr("data-sort")
//		link, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(3) > a").Attr("href")
//		volume, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4) > span[class=volume]").Attr("data-usd")
//		volumeNative, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4) > span[class=volume]").Attr("data-native")
//		price, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(5) > span[class=price]").Attr("data-usd")
//		volumePercentage, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(6)").Attr("data-sort")
//		marketInfo.Rank = util.ConvertStringToInt64(rank, 0)
//		marketInfo.Name = name
//		marketInfo.Pair = pair
//		marketInfo.Link = link
//		marketInfo.Volume = util.ConvertStringToFloat(volume, 0)
//		marketInfo.VolumeNative = util.ConvertStringToFloat(volumeNative, 0)
//		marketInfo.VolumePercentage = util.ConvertStringToFloat(volumePercentage, 0)
//		marketInfo.Price = util.ConvertStringToFloat(price, 0)
//		marketInfos = append(marketInfos, marketInfo)
//	})
//
//	logrus.Infof("market in buffer : parse page data done.")
//	w.Lock()
//	w.marketInfoList = marketInfos
//	w.updateTime = time.Now().Local().Format(util.DATETIMEFORMAT)
//	w.Unlock()
//}
