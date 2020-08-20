package main

import (
	"github.com/adschain/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/wlcy/tradehome-service/util"
	"net/http"
	"time"
)

var (
	cfg = pflag.StringP("config", "c", "config.yaml", "server config file path")
)

/*
 1. 每天零点，调用合约接口 divide 将合约分红池金额转账到奖金池账户，用于发放奖励
 2. 遍历合约用户，合约接口 userCount， getUserInfo, 获取用户质押量，计算质押总量及每个用户应分的reward
 3. 将结果写入数据库表，并不发放
 4. 调用发放接口，完成发放
*/
func main() {
	pflag.Parse()
	if err := Init(*cfg); err != nil {
		panic(err)
	}
	InitGo(1000)

	gin.SetMode(viper.GetString("run_mode"))

	g := gin.New()

	Load(
		g,
	)
	go func() {
		for {
			log.Infof("start divide stm bonus...")
			now := time.Now().UTC()
			// 计算下一个
			//零点
			next := now.Add(time.Hour * 24).UTC()
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			log.Infof("start divide stm bonus started:[%v]", time.Now().UTC().Format(util.DATETIMEFORMAT))
			//奖池金额打入PoolAddress
			trxHash, err := CallDivide(PoolAddress)

			//查询奖池金额并写入奖池金额提取记录
			poolReward := &DividePoolInfo{}
			poolReward.CallerAddress = OwnerAddr
			poolReward.PoolAddress = PoolAddress
			poolReward.TrxHash = trxHash

			if err != nil {
				poolReward.TrxResult = "1"
				InsertPoolReward(poolReward)
				log.Infof("CallDivide for address: %v err:%v", PoolAddress, err)
				return
			}
			poolReward.TrxResult = "0"
			totalBonus, err := getPoolTotalReward()
			poolReward.TotalAmount = totalBonus
			InsertPoolReward(poolReward)
			if err != nil || totalBonus == 0 {
				log.Infof("getPoolTotalReward ==0 or err:%v", err)
				return
			}

			queryAndSavePreparedReward(totalBonus)
			return
		}
	}()
	log.Infof("start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Infof(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

func getPoolTotalReward() (int64, error) {
	//校验账户余额
	totalBonus, err := GetAccountBalance(PoolAddress)
	if err != nil {
		log.Infof("GetAccountBalance for address: %v err:%v", PoolAddress, err)
		return 0, err
	}
	return totalBonus, nil
}

// 查询合约，保存中间结果进入数据
func queryAndSavePreparedReward(totalBonus int64) error {
	//查询合约用户数量
	users, totalPledge, err := getUserInfoSimple()
	if err != nil {
		log.Infof("getUserInfoSimple err:%v", err)
		return err
	}
	//每个用户扣除0.1trx作为转账手续费
	totalBonus = totalBonus - int64(len(users)*100000)
	//计算并直接发放奖励
	//rewardResults, _ := sendReward(users, totalPledge, totalBonus)
	//计算奖励并写入数据库
	rewardResults, _ := sendPrepareReward(users, totalPledge, totalBonus)
	err = InsertRewardStandby(rewardResults)
	if err != nil {
		log.Infof("InsertRewardStandby err:%v", err)
		return err
	}
	return nil
}

func getUserInfoSimple() ([]*SmartUserInfo, int64, error) {
	users := make([]*SmartUserInfo, 0)
	var totalPledge int64
	//查询合约用户数量
	userTotal, err := CallUserCount()
	if err != nil || userTotal == 0 {
		log.Infof("CallUserCount nil or err:%v ", err)
		return nil, totalPledge, nil
	}
	//遍历合约用户信息
	//users := getSmartUserInfos(0, userTotal)

	for i := int64(0); i <= userTotal; i++ {
		uid := 666 + i
		user, err := getSmartUserInfo(uid)
		if err != nil || user == nil {
			log.Infof("getSmartUserInfo for id:%v nil or err:%v", i, err)
			continue
		}
		totalPledge += user.PledgeAmount
		users = append(users, user)
	}
	return users, totalPledge, nil
}
