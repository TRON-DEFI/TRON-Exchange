package buffer

import (
	"encoding/json"
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/convert"
	"github.com/wlcy/tradehome-service/entity"
	"strconv"
	"sync"
	"time"
)

var _orderInfoBuffer *orderInfoBuffer
var onceOrderBufferOnce sync.Once

type orderInfoBuffer struct {
	sync.RWMutex
	orderRecords map[int64]*entity.OrderRecord
}

func GetOrderInfoBuffer() *orderInfoBuffer {
	onceOrderBufferOnce.Do(func() {
		_orderInfoBuffer = &orderInfoBuffer{}
		_orderInfoBuffer.loadOrderRecordBuffer()
		go orderRecordBufferLoader()

	})
	return _orderInfoBuffer
}

//对外获取数据函数
func (b *orderInfoBuffer) GetOrderRecord() map[int64]*entity.OrderRecord {
	return b.orderRecords
}

// 定时重载数据
func orderRecordBufferLoader() {
	for {
		_priceInfBuffer.loadPriceInfsBuffer()
		time.Sleep(30 * time.Second)
	}
}

//加载数据
func (b *orderInfoBuffer) loadOrderRecordBuffer() {
	b.Lock()
	b.orderRecords = OrderRecords()
	log.Info("loadOrderRecordBuffer success")
	if b.orderRecords != nil {
		logData, _ := json.Marshal(b.orderRecords)
		log.Debugf("loadOrderRecordBuffer data :[%v]", string(logData))
	}

	b.Unlock()
}

func OrderRecords() map[int64]*entity.OrderRecord {
	strSQL := `select order_id,order_type,first_token_balance,second_token_balance,cur_turnover,first_token_address from market_order `
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)

	if err != nil {
		log.Errorf(err, "OrderRecords error")
		return nil
	}
	if dataPtr == nil {
		log.Error("OrderRecords dataPtr is nil", nil)
		return nil
	}
	return convert2OrderRecordMap(dataPtr)
}

func convert2OrderRecordMap(dataPtr *mysql.TronDBRows) map[int64]*entity.OrderRecord {
	data := make(map[int64]*entity.OrderRecord, 0)

	for dataPtr.NextT() {
		orderRecord := &entity.OrderRecord{}
		orderID, _ := strconv.ParseInt(dataPtr.GetField("order_id"), 10, 64)
		firstTokenAddr := dataPtr.GetField("first_token_address")
		//按精度处理
		tempCurTurnover, _ := strconv.ParseFloat(dataPtr.GetField("cur_turnover"), 10)
		orderRecord.CurTurnover = fmt.Sprintf("%f", convert.DealVolumeByAddr(tempCurTurnover, firstTokenAddr))
		orderRecord.Schedule = getSchedule(dataPtr.GetField("first_token_balance"), dataPtr.GetField("second_token_balance"), dataPtr.GetField("order_type"), orderRecord.CurTurnover)

		data[orderID] = orderRecord
	}

	return data
}
func getSchedule(fTokenBalanceStr string, sTokenBalanceStr string, orderTypeStr string, curTurnover string) string {
	orderType, _ := strconv.Atoi(orderTypeStr)
	ct, _ := strconv.ParseInt(curTurnover, 10, 64)
	var tokenBalance int64
	//买单-第一个的 tokenBalance，卖单-第二个的 tokenBalance
	if orderType == 0 {
		tokenBalance, _ = strconv.ParseInt(fTokenBalanceStr, 10, 64)

	} else {
		tokenBalance, _ = strconv.ParseInt(sTokenBalanceStr, 10, 64)
	}
	schedule := float64(ct / tokenBalance)
	log.Debugf("getSchedule error :[%f]\n", schedule)
	return fmt.Sprintf("%f", schedule)
}
