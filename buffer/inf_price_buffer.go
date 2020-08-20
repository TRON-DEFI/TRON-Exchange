package buffer

import (
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/module"
	"github.com/wlcy/tradehome-service/util"
	"sync"
	"time"
)

var _priceInfBuffer *priceInfBuffer
var oncePirceInfBufferOnce sync.Once

type priceInfBuffer struct {
	sync.RWMutex
	priceInfs map[int64]*entity.PriceInf
}

func GetPirceInfBuffer() *priceInfBuffer {
	oncePirceInfBufferOnce.Do(func() {
		_priceInfBuffer = &priceInfBuffer{}
		_priceInfBuffer.loadPriceInfsBuffer()
		go priceInfsBufferLoader()

	})
	return _priceInfBuffer
}

//对外获取数据函数
func (b *priceInfBuffer) GetPriceInfs() map[int64]*entity.PriceInf {
	return b.priceInfs
}

// 定时重载数据
func priceInfsBufferLoader() {
	for {
		_priceInfBuffer.loadPriceInfsBuffer()
		time.Sleep(30 * time.Second)
	}
}

//加载数据
func (b *priceInfBuffer) loadPriceInfsBuffer() {
	b.RLock()
	//缓存前100条
	b.priceInfs = module.PriceInfs(util.TimeStampBefore24h())
	b.RUnlock()
}
