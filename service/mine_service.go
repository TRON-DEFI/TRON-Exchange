package service

import (
	"errors"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/contract"
	"github.com/wlcy/tradehome-service/module"
	"time"
)

func HandleMineStatistic(dateTime int64) error {
	list, err := module.QueryTransRecord(dateTime)
	if err != nil {
		return err
	}
	for _, item := range list {
		temp, _ := module.QueryMineStatistic(item.OrderId, item.OrderCreateTime)
		if temp.Id == 0 {
			if err := module.AddMineStatistic(item); err == nil {
				module.UpdateTransRecordMineStatusById(item.Id)
			}
		}
	}
	return nil
}

func HandleMineReward4Send() error {
	count := int(0)
exec:
	start := time.Now()
	list, err := module.QueryMineReward4Send()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	for _, item := range list {
		if item.TotalReward > 0 {
			if _, err := contract.TransferAsset2(item.UserAddress, item.TotalReward); err == nil {
				module.UpdateMineStatisticStatusByUserAddr(item.UserAddress)
			} else {
				log.Errorf(err, "TransferAsset2 error")
			}
		} else {
			module.UpdateMineStatisticStatusByUserAddr(item.UserAddress)
		}
	}
	cost := time.Since(start)
	log.Infof("SendBonus end, costTime=%v", cost)
	count++
	if count >= 5 {
		err := errors.New("TransferAsset2 more than 5 times")
		log.Errorf(err, "TransferAsset2 more than 5 times")
		return err
	}
	goto exec
}
