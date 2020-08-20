package module

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/contract"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/util"
)

func QueryTransRecord(dateTime int64) ([]*entity.MineStatistic, error) {
	strSQL := fmt.Sprintf(`
		select id, pair_id, order_id, owner_address, deal_order_id, deal_owner_address, order_type,
		first_token_address, second_token_address, amountB, create_time
		from market_transaction
		where status=0 and mine_status=0 and create_time<%v`, dateTime)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf(err, "QueryTransRecord error sql:%v", strSQL)
		return nil, err
	}
	list := make([]*entity.MineStatistic, 0)
	for dataPtr.NextT() {
		item := &entity.MineStatistic{}
		item.Id = util.ConvertDBValueToInt64(dataPtr.GetField("id"))
		item.PairId = util.ConvertDBValueToInt64(dataPtr.GetField("pair_id"))
		item.OrderType = util.ConvertDBValueToInt64(dataPtr.GetField("order_type"))
		if item.OrderType == 0 {
			item.OrderId = util.ConvertDBValueToInt64(dataPtr.GetField("order_id"))
			item.UserAddress = dataPtr.GetField("owner_address")
		} else if item.OrderType == 1 {
			item.OrderId = util.ConvertDBValueToInt64(dataPtr.GetField("deal_order_id"))
			item.UserAddress = dataPtr.GetField("deal_owner_address")
		}
		item.FirstTokenAddress = dataPtr.GetField("first_token_address")
		item.SecondTokenAddress = dataPtr.GetField("second_token_address")
		item.TradeVolume = util.ConvertDBValueToInt64(dataPtr.GetField("amountB"))
		if contract.MineRate == 0 {
			contract.MineRate = 1000
		}
		item.RewardVolume = int64(item.TradeVolume / contract.MineRate)
		item.OrderCreateTime = util.ConvertDBValueToInt64(dataPtr.GetField("create_time"))
		list = append(list, item)
	}
	return list, nil
}

func UpdateTransRecordMineStatusById(id int64) error {
	strSQL := fmt.Sprintf(`
		update market_transaction set mine_status=%v where id=%v`,
		1, id)
	log.Info(strSQL)
	_, _, err := mysql.ExecuteSQLCommand(strSQL)
	if err != nil {
		log.Errorf(err, "UpdateTransRecordMineStatusById error:%s", strSQL)

		return err
	}
	log.Info("UpdateTransRecordMineStatusById success")

	return nil
}

func AddMineStatistic(item *entity.MineStatistic) error {
	strSQL := fmt.Sprintf(`
		insert into market_mine_statistic 
		(pair_id, order_id, order_type, user_address, first_token_address, second_token_address, trade_volume, reward_volume, status, order_create_time)
		values(%v, %v, %v, '%v', '%v', '%v', %v, %v, %v, %v)`,
		item.PairId, item.OrderId, item.OrderType, item.UserAddress, item.FirstTokenAddress, item.SecondTokenAddress,
		item.TradeVolume, item.RewardVolume, 0, item.OrderCreateTime)
	log.Info(strSQL)
	_, _, err := mysql.ExecuteSQLCommand(strSQL)
	if err != nil {
		log.Errorf(err, "AddMineStatistic error:%s", strSQL)
		return err
	}
	log.Info("AddMineStatistic success")
	return nil
}

func UpdateMineStatisticStatusByUserAddr(userAddress string) error {
	strSQL := fmt.Sprintf(`
		update market_mine_statistic set status=%v where user_address='%v'`,
		1, userAddress)
	log.Info(strSQL)
	_, _, err := mysql.ExecuteSQLCommand(strSQL)
	if err != nil {
		log.Errorf(err, "UpdateMineStatisticStatus error:%s", strSQL)

		return err
	}
	log.Info("UpdateMineStatisticStatus success")

	return nil
}

func QueryMineStatistic(orderId, orderCreateTime int64) (*entity.MineStatistic, error) {
	strSQL := fmt.Sprintf(`select id, user_address from market_mine_statistic where order_id=%v and order_create_time=%v`, orderId, orderCreateTime)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)

	if err != nil {
		log.Errorf(err, "QueryMineStatistic error")
		return nil, err
	}
	item := &entity.MineStatistic{}
	for dataPtr.NextT() {
		item.Id = util.ConvertDBValueToInt64(dataPtr.GetField("id"))
		item.UserAddress = dataPtr.GetField("user_address")
	}
	return item, nil
}

func QueryMineReward4Send() ([]*entity.MineReward, error) {
	strSQL := fmt.Sprintf(`select user_address, sum(reward_volume) as totalReward from market_mine_statistic where status=0 group by user_address`)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)

	if err != nil {
		log.Errorf(err, "QueryMineReward4Send error")
		return nil, err
	}
	list := make([]*entity.MineReward, 0)
	for dataPtr.NextT() {
		item := &entity.MineReward{}
		item.UserAddress = dataPtr.GetField("user_address")
		item.TotalReward = util.ConvertDBValueToInt64(dataPtr.GetField("totalReward"))
		list = append(list, item)
	}
	return list, nil
}
