package main

import (
	"encoding/base64"
	"fmt"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tradehome-engin/core/grpcclient"
	"github.com/wlcy/tradehome-engin/core/utils"
	"github.com/wlcy/tradehome-service/util"
	"testing"
	"time"
)

func init() {
	if err := Init("config.yaml"); err != nil {
		panic(err)
	}
	InitGo(1000)
}

//合约divide方法测试
func TestCallDivide(t *testing.T) {
	_, err := CallDivide(PoolAddress)
	fmt.Printf("CallDivide for address: %v err:%v", PoolAddress, err)
}

//合约userCount测试
func TestCallUserCount(t *testing.T) {
	userCount, err := CallUserCount()
	fmt.Printf("CallUserCount, userCount: %v err:%v", userCount, err)
}

//合约userInfo测试
func TestCallAllUserInfo(t *testing.T) {
	userInfos, pledgeTotal, err := getUserInfoSimple()
	userStr, _ := util.JSONObjectToString(userInfos)
	fmt.Printf("CallUserInfo, userInfos: %v, totalPledge:%v, err:%v", userStr, pledgeTotal, err)
}

//查询userInfo计算结果并写入数据
func TestQueryAndSavePreparedReward(t *testing.T) {
	totalBonus := int64(200000000)
	err := queryAndSavePreparedReward(totalBonus)
	fmt.Printf("queryAndSavePreparedReward err:%v", err)
}

//发放奖励，并更新数据库结果
func TestSendManualReward(t *testing.T) {
	preparedList, err := queryPreparedAndSendReward()
	preparedListStr, _ := util.JSONObjectToString(preparedList)
	fmt.Printf("queryPreparedReward preparedList:%v, err:%v", preparedListStr, err)

	results := sendRealReward(preparedList)
	err = UpdateRewardResult(results)
	fmt.Printf("UpdateRewardResult preparedList err:%v", err)

}
func TestCheckTrx(t *testing.T) {
	msg, msg2 := checkTrx("bb7688dbd0c4f1fcf93e960ef8ffd9b2cd11a91388faaabd8687df36da103b9d")
	fmt.Printf("msg1:%v,msg2:%v", msg, msg2)
}

//用户余额查询
func TestGetAccountBalance(t *testing.T) {
	balance, err := GetAccountBalance("TKaTjCYPyU3NZWgGckQEVJyb4N6c82gfcU")
	fmt.Printf("GetAccountBalance result:%v-%v\n", balance, err)
}

//trx转账测试
func TestTransfer(t *testing.T) {

	tran := new(core.TransferContract)
	tran.OwnerAddress = utils.Base58DecodeAddr("TDXkW6HBqJHBChAvL9fvWQuXDdqrHCEEw9") // me
	tran.ToAddress = utils.Base58DecodeAddr("TKaTjCYPyU3NZWgGckQEVJyb4N6c82gfcU")    // 闫诤
	tran.Amount = 1000000                                                            // sun
	_ = tran
	trx, err := utils.BuildTransaction(core.Transaction_Contract_TransferContract, tran, []byte("test"))
	//转账到新地址，余额需要是转账金额+ 0.1 trx，否则报错
	// 	contract validate error : Validate TransferContract error, balance is not sufficient.
	// CONTRACT_VALIDATE_ERROR

	if nil != err {
		fmt.Printf("build trx failed:%v\n", err)
		return
	}
	client := grpcclient.GetRandomWallet()
	block, err := client.GetNowBlock()
	if nil != err {
		fmt.Printf("getBlock failed:%v\n", err)
		return
	}

	blockHash := utils.CalcBlockHash(block)
	trx.RawData.RefBlockHash = blockHash[8:16]
	trx.RawData.RefBlockBytes = utils.BinaryBigEndianEncodeInt64(block.BlockHeader.RawData.Number)[6:8]
	trx.RawData.Timestamp = time.Now().UTC().UnixNano() / 1000000
	trx.RawData.Expiration = time.Now().UTC().Add(5*time.Minute).UnixNano() / 1000000

	fmt.Println("3xxxxxxxxxxxxxxx")

	sign, err := utils.SignTransaction(trx, "")
	if nil != err {
		fmt.Printf("sign failed:%v\n", err)
		return
	}

	trx.Signature = append(trx.Signature, sign)

	ret, err := client.BroadcastTransaction(trx)
	if nil != ret {
		fmt.Printf("%v\n%s\n%v\n", err, ret.Message, ret.Code)
	}
	hash := utils.HexEncode(utils.CalcTransactionHash(trx))

	fmt.Printf("transaction hash:%v", hash)

}

// 私钥加密
func TestAAA(t *testing.T) {
	var aeskey = []byte("dblyztradestmvid")
	pass := []byte("F1BEFC60BA1216E20F89998CDD18CFC42316B6A928456109559A64CA6DCBECCF")
	xpass, err := util.AesEncrypt(pass, aeskey)
	if err != nil {
		fmt.Println(err)
		return
	}

	pass64 := base64.StdEncoding.EncodeToString(xpass)
	fmt.Printf("加密后:%v\n", pass64)

	bytesPass, err := base64.StdEncoding.DecodeString(pass64)
	if err != nil {
		fmt.Println(err)
		return
	}

	tpass, err := util.AesDecrypt(bytesPass, aeskey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("解密后:%s\n", tpass)
}
