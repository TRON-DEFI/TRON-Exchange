package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-engin/core/tools"
	"github.com/wlcy/tradehome-engin/core/utils"
	"github.com/wlcy/tradehome-engin/web/common"
	"github.com/wlcy/tradehome-service/util"
	"math/big"
	"strings"
	"sync"
	"time"
)

var _abiEntry *tools.AbiEntry
var _userCount *tools.Method   // 返回用户数量, owner用户ID == 666
var _getUserInfo *tools.Method //根据id查询用户信息, id == 0 自动为自己，普通用户只能查自己
var _divide *tools.Method      //  将分红金额打入接受地址

var ContractAddress, OwnerAddr, PrivateKey string
var PoolAddress, PoolPrivateKey string //奖池地址
var SmartNode string
var Message = "STM Divide"
var feeLimit = int64(10000000)

func InitSmart() error {
	defer util.CatchError()
	addr := utils.GetRandFullNodeAddr()
	if strings.Contains(addr, ":") {
		SmartNode = addr
	} else {
		SmartNode = fmt.Sprintf("%v:%v", addr, utils.DefaultGrpPort)
	}
	log.Infof("init smart node:%v", SmartNode)
	smartContract, err := common.GetWalletClient().GetSmartContract(ContractAddress)
	abi, err := json.Marshal(smartContract.Abi)
	if nil != err {
		log.Infof("json marshal error for contract abi\n")
		return err
	}
	_abiEntry, err := tools.GetABI(string(abi))
	if nil != err {
		return err
	}
	_userCount = _abiEntry.GetMethod("UserCount")
	if nil == _userCount {
		return errors.New("UserCount nil")
	}
	_getUserInfo = _abiEntry.GetMethod("getUserInfo")
	if nil == _getUserInfo {
		return errors.New("getUserInfo nil")
	}
	_divide = _abiEntry.GetMethod("divide")
	if nil == _divide {
		return errors.New("divide nil")
	}
	return nil
}

//CallUserCount 获取用户数量
//UserCount() returns (uint)
//+ 返回用户数量, owner用户ID == 666
func CallUserCount() (int64, error) {
	defer util.CatchError()
	var userCount int64
	record := &CallRecord{
		Owner:    OwnerAddr,
		Contract: ContractAddress,
		Method:   "userCount",
	}

	record.Data, record.Err = _userCount.Pack()
	if nil != record.Err {
		log.Infof("pack userCount err:[%v]", record.Err)
		return 0, record.Err
	}

	result, err := tools.TriggerContract(OwnerAddr, ContractAddress, 0, record.Data, SmartNode)
	if err != nil {
		return 0, err
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		return 0, err
	}
	log.Infof("result:%v\n%s\n", string(resultJson), result.Result.Message)

	rawData := result.ConstantResult[0]
	paramsOut, _ := _userCount.Outputs.UnpackValues(rawData)
	fmt.Printf("%v--%T\n", paramsOut[0], paramsOut[0])
	if err != nil {
		log.Infof("GenTriggerSmartContract orderID err:[%v]", err)
		return 0, err
	}
	userCount = paramsOut[0].(*big.Int).Int64()
	log.Infof("call smart UserCount:%v\n", userCount)
	return userCount, nil
}

//CallUserInfo 根据id查询用户信息, id == 0 自动为自己，普通用户只能查自己
//+ returns
//+ address userAddr: 用户地址
//+ uint balance: 用户代币余额
//+ uint pledgeAmount: 用户代币质押量
//+ uint lockAmount: 用户代币锁定量
//+ uint stakingAmount: 用户TRX质押量
//+ uint unstakingAmount: 用户TRX锁定量
//+ uint rate: 用户staking代币收益率，100 == 100%
//+ uint invitee_amount: 用户所邀请的人的staking总额
//+ uint invitee_count: 用户所邀请的人数
//+ uint staking_otken_withdraw_time: 用户上次提取代币收益的时间 UTC
//+ uint parentID: 用户的邀请人ID
func CallUserInfo(userID int64) (*SmartUserInfo, error) {
	smartUserInfo := &SmartUserInfo{}
	defer util.CatchError()
	record := &CallRecord{
		Owner:    OwnerAddr,
		Contract: ContractAddress,
		Method:   "getUserInfo",
	}

	record.Data, record.Err = _getUserInfo.Pack(tools.GenAbiInt(userID))
	if nil != record.Err {
		log.Infof("pack userID err:[%v]", record.Err)
		return smartUserInfo, record.Err
	}

	result, err := tools.TriggerContract(OwnerAddr, ContractAddress, 0, record.Data, SmartNode)
	if err != nil {
		return smartUserInfo, err
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		return smartUserInfo, err
	}
	fmt.Printf("CallUserInfo result:%v\n%s\n", string(resultJson), result.Result.Message)

	rawData := result.ConstantResult[0]
	paramsOut, _ := _getUserInfo.Outputs.UnpackValues(rawData)
	if len(paramsOut) != 11 { //返回参数必须等于11，并且顺序不变
		log.Infof("call smart contract function getUserInfo() error:[return length ! =11]")
		return smartUserInfo, errors.New("return length ! =11")
	}
	smartUserInfo.UserAddr = tools.GetTronAddrFromAbiAddress(paramsOut[0])
	smartUserInfo.Balance = paramsOut[1].(*big.Int).Int64()
	smartUserInfo.PledgeAmount = paramsOut[2].(*big.Int).Int64()
	smartUserInfo.LockAmount = paramsOut[3].(*big.Int).Int64()
	smartUserInfo.StakingAmount = paramsOut[4].(*big.Int).Int64()
	smartUserInfo.UnstakingAmount = paramsOut[5].(*big.Int).Int64()
	smartUserInfo.Rate = paramsOut[6].(*big.Int).Int64()
	smartUserInfo.InviteeAmount = paramsOut[7].(*big.Int).Int64()
	smartUserInfo.InviteeCount = paramsOut[8].(*big.Int).Int64()
	smartUserInfo.StakingOtkenWithdrawTime = paramsOut[9].(*big.Int).Int64()
	smartUserInfo.ParentID = paramsOut[9].(*big.Int).Int64()

	return smartUserInfo, nil

}

//CallDivide
//+ 分红接口, 将分红额度从余额中提走, 最多每24小时调用一次, 更新 lastDivideTime 及 totalDivideAmount
//+ 分红额度 == totalStaking * 5 / 100;
//+ address addr: 分红金额接受地址，用与给用户转账的资金账户
func CallDivide(address string) (string, error) {
	var trxHash string
	defer util.CatchError()
	record := &CallRecord{
		Owner:    OwnerAddr,
		Contract: ContractAddress,
		Method:   "divide",
	}

	record.Data, record.Err = _divide.Pack(tools.GenAbiAddress(address))
	if nil != record.Err {
		log.Infof("pack address err:[%v]", record.Err)
		return trxHash, record.Err
	}
	//必须使用Broadcast方法
	//result, err := tools.TriggerContract(OwnerAddr, ContractAddress, 0, record.Data, SmartNode)
	ctxType, ctx, err := tools.GenTriggerSmartContract(OwnerAddr, ContractAddress, record.CallValue, record.Data)
	if nil != err {
		return trxHash, record.Err
	}
	record.TrxHash, record.Return, record.Err = tools.BroadcastCtxWithFeeLimit(ctxType, ctx, nil, PrivateKey, 0, feeLimit)
	if nil != record.Err {
		return trxHash, record.Err
	}
	resultJson, err := json.Marshal(record.Return)
	if err != nil {
		return trxHash, err
	}
	result := checkTrxResult(record.TrxHash)
	if result {
		trxHash = record.TrxHash
	}
	fmt.Printf("CallDivide result:%v---%v\n%v\n%v\n", record.TrxHash, result, string(resultJson), record.Err)

	return trxHash, nil

}

//多线程获取用户信息
func getSmartUserInfos(start, end int64) []*SmartUserInfo {
	var wg sync.WaitGroup
	resultInfo := make([]*SmartUserInfo, end-start+1)
	now := time.Now().UTC()
	num := int64(0)
	for i := start; i <= end; i++ {
		num = num + 1
		wg.Add(1)
		go func(index, offset int64) {
			defer wg.Done()
			getGo()
			defer releaseGo()
			smartUserInfo, err := getSmartUserInfo(index) //获取订单信息，100次机会
			// fmt.Printf("CallOrderInfo result:orderID:[%v],smartOrderInfo:[%#v] err:[%v]", smartOrderInfo, err)
			if err != nil || smartUserInfo == nil || (smartUserInfo != nil && smartUserInfo.UserAddr == "") {
				log.Infof("call smart userInfo:[%v] err:[%v]\n", index, err)
				return
			}
			fmt.Printf("offset:[%v],userAddress:[%v],done\n", offset-1, smartUserInfo.UserAddr)
			resultInfo[offset] = smartUserInfo

		}(i, num)
	}
	wg.Wait()
	log.Infof("get userInfo from smart done:[%v],timecost:[%v]\n", len(resultInfo), time.Now().UTC().Sub(now))
	return resultInfo
}

func getSmartUserInfo(userID int64) (*SmartUserInfo, error) {
	var ret *SmartUserInfo
	var err error
	tryCnt := 1000
	for tryCnt > 0 {
		tryCnt--
		ret, err = CallUserInfo(userID)
		//fmt.Printf("CallOrderInfo result:smartOrderInfo:[%#v] err:[%v]", ret, err)
		if err != nil || ret == nil || (ret != nil && ret.UserAddr == "") {
			// log.Errorf("call smart transfer err:[%v]", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return ret, err
	}
	return ret, err
}

var chanGo chan bool

func InitGo(size int) {
	chanGo = make(chan bool, size)

	for i := 0; i < size; i++ {
		chanGo <- true
	}
}

func getGo() {
	<-chanGo
}

func releaseGo() {
	chanGo <- true
}
