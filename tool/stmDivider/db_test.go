package main

import (
	"fmt"
	"github.com/wlcy/tradehome-service/util"
	"testing"
)

func init() {
	if err := Init("config.yaml"); err != nil {
		panic(err)
	}
	InitGo(1000)
}

func TestInsertPrepareReward(t *testing.T) {
	users := make([]*SmartUserInfo, 0)
	var totalPledge int64
	totalBonus := int64(200000000)

	for i := int64(0); i <= 5; i++ {
		user := &SmartUserInfo{}
		user.UserAddr = fmt.Sprintf("14321215321-%v", i)
		user.PledgeAmount = int64(3000000000) / (i + 1)
		totalPledge = totalPledge + user.PledgeAmount
		users = append(users, user)
	}
	rewardResults, _ := sendPrepareReward(users, totalPledge, totalBonus)
	err := InsertRewardStandby(rewardResults)
	fmt.Printf("queryAndSavePreparedReward err:%v", err)
}

func TestQueryPrepareReward(t *testing.T) {
	results, err := queryPreparedAndSendReward()
	userStr, _ := util.JSONObjectToString(results)
	fmt.Printf("queryPreparedAndSendReward:%v, err:%v", userStr, err)

	for idx, result := range results {
		result.TrxHash = fmt.Sprintf("rewqtewqrewqrewrew-%v", idx)
		var ss string
		if idx%2 == 0 {
			ss = "DEFULT"
		} else {
			ss = "SUCCESS"
		}
		result.TrxResult = ss
	}
	err = UpdateRewardResult(results)
}

func TestInsertPoolReward(t *testing.T) {
	dd := &DividePoolInfo{
		CallerAddress: "caller1",
		PoolAddress:   "pool1",
		TotalAmount:   1000000023212,
		TrxHash:       "215432rewqtqwrweqt-1",
		TrxResult:     "0",
	}
	InsertPoolReward(dd)
}
