package buffer

import (
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/buffer2"
	"github.com/wlcy/tradehome-service/module"
	"github.com/wlcy/tradehome-service/util"
	"time"
)

//启动定时任务，获取交易对0点价格，并写入exchange_zero_price
// 由于两个服务，所以数据库中两条同样的记录
func StartTimer() {
	go func() {
		for {
			log.Debugf("start fetch zero click price for exchangeID")
			now := time.Now().UTC()
			// 计算下一个零点
			next := now.Add(time.Hour * 24).UTC()
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			log.Infof("get zero price for exchange pair started:[%v]", time.Now().UTC().Format(util.DATETIMEFORMAT))
			exchanges := buffer2.GetMarketPairBuffer().GetMarketPairs()
			log.Infof("get exchange len:[%v]", len(exchanges))
			for _, exchange := range exchanges {
				log.Infof("start update zero price for exchangeID:[%v]", exchange.ID)
				zeroPrice, err := module.GetExchangeZeroPrice(exchange.ID)
				if nil == err && zeroPrice > 0 {
					err = module.InsertZeroPrice(exchange.ID, zeroPrice)
					if nil != err {
						log.Errorf(err, "insert zero price in db error")
					} else {
						log.Info("fetch zero price for exchangeID success")
					}
				} else {
					log.Error("get exchange zeroPrice error", nil)
				}
			}

		}
	}()
}
