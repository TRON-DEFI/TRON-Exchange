package main

import (
	"github.com/adschain/log"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tradehome-engin/core/grpcclient"
	"github.com/wlcy/tradehome-engin/core/utils"
	"strings"
	"time"
)

func sendPrepareReward(users []*SmartUserInfo, totalPledge, totalBonus int64) ([]*UserRewardInfo, error) {
	results := make([]*UserRewardInfo, 0)
	for _, user := range users {
		if user == nil {
			continue
		}
		result := &UserRewardInfo{}
		result.UserAddr = user.UserAddr
		result.TotalBonus = totalBonus
		result.TotalPledged = totalPledge
		result.Pledged = user.PledgeAmount
		result.TrxHash = ""
		result.TrxResult = ""
		if user.UserAddr == "" || user.PledgeAmount == 0 {
			result.Bonus = 0

			results = append(results, result)
			continue
		}
		reward := user.PledgeAmount * totalBonus / totalPledge
		result.Bonus = reward
		results = append(results, result)
	}
	return results, nil
}

func sendReward(users []*SmartUserInfo, totalPledge, totalBonus int64) ([]*UserRewardInfo, error) {
	results := make([]*UserRewardInfo, 0)
	for _, user := range users {
		if user == nil {
			continue
		}
		result := &UserRewardInfo{}
		result.UserAddr = user.UserAddr
		result.TotalBonus = totalBonus
		result.TotalPledged = totalPledge
		result.Pledged = user.PledgeAmount
		if user.UserAddr == "" || user.PledgeAmount == 0 {
			result.Bonus = 0
			result.TrxHash = ""
			result.TrxResult = ""
			results = append(results, result)
			continue
		}
		reward := user.PledgeAmount * totalBonus / totalPledge
		trxHash := Transfer(PoolAddress, user.UserAddr, reward, Message, PoolPrivateKey)
		trxResult, _ := checkTrx(trxHash)
		result.Bonus = reward
		result.TrxHash = trxHash
		result.TrxResult = trxResult
		results = append(results, result)
	}
	return results, nil
}

func sendRealReward(users []*UserRewardInfo) []*UserRewardInfo {
	results := make([]*UserRewardInfo, 0)
	for _, user := range users {

		if user == nil {
			continue
		}
		result := &UserRewardInfo{}
		result.UserAddr = user.UserAddr
		result.TotalBonus = user.TotalBonus
		result.TotalPledged = user.TotalPledged
		result.Bonus = user.Bonus
		result.Pledged = user.Pledged

		if user.UserAddr == "" || user.Pledged == 0 || user.TotalBonus == 0 || user.TotalPledged == 0 {
			result.TrxHash = ""
			result.TrxResult = ""
			results = append(results, result)
			continue
		}
		trxHash := Transfer(PoolAddress, user.UserAddr, result.Bonus, Message, PoolPrivateKey)
		trxResult, _ := checkTrx(trxHash)
		result.TrxHash = trxHash
		result.TrxResult = trxResult
		results = append(results, result)
	}
	return results
}
func GetAccountBalance(addr string) (int64, error) {
	var ret int64
	var err error
	tryCnt := 10000
	for tryCnt > 0 {
		tryCnt--
		ret, err = getAccountBalance(addr)
		//fmt.Printf("CallOrderInfo result:smartOrderInfo:[%#v] err:[%v]", ret, err)
		if err != nil {
			// log.Errorf("call smart transfer err:[%v]", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return ret, err
	}
	return ret, err
}

//getAccountBalance 查询用户余额
func getAccountBalance(addr string) (int64, error) {
	account, err := grpcclient.GetRandomWallet().GetAccount(addr)
	if err != nil {
		log.Infof("get account:[%v] info err:%v", addr, err)
		return 0, err
	}
	return account.Balance, nil
}

//Transfer TRX转账操作
func Transfer(from, to string, amount int64, message string, pk string) string {
	var hash string
	tran := new(core.TransferContract)
	tran.OwnerAddress = utils.Base58DecodeAddr(from) //
	tran.ToAddress = utils.Base58DecodeAddr(to)      //
	tran.Amount = amount                             // sun
	trx, err := utils.BuildTransaction(core.Transaction_Contract_TransferContract, tran, []byte(message))

	if nil != err {
		log.Infof("build trx failed:%v\n", err)
		return hash
	}
	client := grpcclient.GetRandomWallet()
	block, err := client.GetNowBlock()
	if nil != err {
		log.Infof("getBlock failed:%v\n", err)
		return hash
	}

	blockHash := utils.CalcBlockHash(block)
	trx.RawData.RefBlockHash = blockHash[8:16]
	trx.RawData.RefBlockBytes = utils.BinaryBigEndianEncodeInt64(block.BlockHeader.RawData.Number)[6:8]
	trx.RawData.Timestamp = time.Now().UTC().UnixNano() / 1000000
	trx.RawData.Expiration = time.Now().UTC().Add(5*time.Minute).UnixNano() / 1000000

	//fmt.Println("3xxxxxxxxxxxxxxx")

	sign, err := utils.SignTransaction(trx, pk)
	if nil != err {
		log.Infof("sign failed:%v\n", err)
		return hash
	}

	trx.Signature = append(trx.Signature, sign)

	ret, err := client.BroadcastTransaction(trx)

	if nil != ret {
		log.Infof("%v\n%s\n%v\n", err, ret.Message, ret.Code)
		return hash
	}
	hash = utils.HexEncode(utils.CalcTransactionHash(trx))

	log.Infof("transaction hash:%v", hash)
	return hash
}

func checkTrxResult(trxHash string) bool {
	result, _ := checkTrx(trxHash)
	if strings.ToUpper(result) == "DEFAULT" || strings.ToUpper(result) == "SUCCESS" {
		return true
	}
	return false
}

func checkTrx(trxHash string) (string, string) {
	client := grpcclient.GetRandomWallet()
	if nil != client {
		defer client.Close()
	}
	var trxInfo *core.TransactionInfo
	for {
		trxInfo, _ = client.GetTransactionInfoByID(trxHash)
		if nil != trxInfo && nil != trxInfo.Receipt {
			break
		}
		time.Sleep(3 * time.Second)
	}
	log.Infof("FUCK--%v\n-->%v-->%s\n", utils.ToJSONStr(trxInfo), trxInfo.Receipt.Result.String(), string(trxInfo.ResMessage))
	return trxInfo.Receipt.Result.String(), string(trxInfo.ResMessage)

}
