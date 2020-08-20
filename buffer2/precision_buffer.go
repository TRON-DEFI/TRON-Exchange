package buffer2

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/mysql"
	"strconv"
	"sync"
	"time"
)

var _precisionBuffer *precisionBuffer
var _precisionOnce sync.Once

type precisionBuffer struct {
	sync.RWMutex
	pairIDFirstPrecisionMap  map[int64]int64
	pairIDSecondPrecisionMap map[int64]int64
	tokenAddrPrecisionMap    map[string]int64
}

func GetPrecisionBuffer() *precisionBuffer {
	_precisionOnce.Do(func() {
		_precisionBuffer = &precisionBuffer{}
		_precisionBuffer.loadPrecision()
		go precisionLoader()
	})
	return _precisionBuffer
}

func precisionLoader() {
	for {
		_precisionBuffer.loadPrecision()
		time.Sleep(5 * time.Second)
	}
}

func (b *precisionBuffer) GetPairIDFirstPrecisionMap() map[int64]int64 {
	return b.pairIDFirstPrecisionMap
}

func (b *precisionBuffer) GetPairIDSecondPrecisionMap() map[int64]int64 {
	return b.pairIDSecondPrecisionMap
}

func (b *precisionBuffer) GetTokenAddrPrecisionMap() map[string]int64 {
	return b.tokenAddrPrecisionMap
}

func (b *precisionBuffer) loadPrecision() {
	b.RWMutex.Lock()
	//从数据库中查询pairID、second_token_addr、对应的精度
	strSQL := fmt.Sprintf(`select id,first_token_addr,second_token_addr,first_token_precision,second_token_precision from market_pair where is_valid = 1`)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)

	if err != nil {
		log.Errorf(err, "queryPrecision error")
		return
	}
	if dataPtr == nil {
		log.Errorf(nil, "queryPrecision dataPtr is nil")
		return
	}
	//将结果赋给全局变量
	pairIDFirstPrecisionMap := make(map[int64]int64)
	pairIDSecondPrecisionMap := make(map[int64]int64)
	tokenAddrPrecisionMap := make(map[string]int64)

	for dataPtr.NextT() {
		pairID, _ := strconv.ParseInt(dataPtr.GetField("id"), 10, 64)
		fistAddr := dataPtr.GetField("first_token_addr")
		secondAddr := dataPtr.GetField("second_token_addr")
		firstPrecision, _ := strconv.ParseInt(dataPtr.GetField("first_token_precision"), 10, 64)
		secondPrecision, _ := strconv.ParseInt(dataPtr.GetField("second_token_precision"), 10, 64)
		log.Infof("pairID:%v, fistAddr:%v, secondAddr:%v, firstPrecision:%v, secondPrecision:%v", pairID, fistAddr, secondAddr, firstPrecision, secondPrecision)

		pairIDFirstPrecisionMap[pairID] = firstPrecision
		pairIDSecondPrecisionMap[pairID] = secondPrecision
		tokenAddrPrecisionMap[fistAddr] = firstPrecision
		tokenAddrPrecisionMap[secondAddr] = secondPrecision
	}
	b.pairIDFirstPrecisionMap = pairIDFirstPrecisionMap
	b.pairIDSecondPrecisionMap = pairIDSecondPrecisionMap
	b.tokenAddrPrecisionMap = tokenAddrPrecisionMap

	b.RWMutex.Unlock()
}
