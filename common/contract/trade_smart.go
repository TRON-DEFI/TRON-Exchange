package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tradehome-engin/core/tools"
	"github.com/wlcy/tradehome-engin/core/utils"
	"github.com/wlcy/tradehome-engin/web/common"
	"github.com/wlcy/tradehome-service/util"
	"math/big"
)

var tradeContract *tools.Method
var cancelContract *tools.Method
var orderIDContract *tools.Method        //合约查询最大订单id
var getOrderInfoContract *tools.Method   //合约查询订单信息
var getOrderInfoContract10 *tools.Method //合约查询trc10订单信息

var getOrderInfo *tools.Method //TODO
//var getDecimal *tools.Method
var setLNFeeLimit int64 = 10000000

var EventLogURL, SmartNode, SmartContractAddr, SmartOwnerAddr, SmartPrivateKey string
var SmartContractAddr10, SmartOwnerAddr10, SmartPrivateKey10 string

var TestNet int

func convertTrcToken(abi *core.SmartContract_ABI) {
	for _, val := range abi.Entrys {
		for _, types := range val.Inputs {
			if types.Type == "trcToken" {
				types.Type = "uint256"
			}
		}
		for _, types := range val.Outputs {
			if types.Type == "trcToken" {
				types.Type = "uint256"
			}
		}
	}
}

func InitTradeSmart() error {
	if TestNet == 1 {
		utils.TestNet = true
		utils.NetName = utils.NetShasta
	}
	defer util.CatchError()
	smartContract, err := common.GetWalletClient().GetSmartContract(SmartContractAddr)
	//fmt.Printf("ss:[%v]", smartContract.Abi.String())
	abi, err := json.Marshal(smartContract.Abi)
	if nil != err {
		fmt.Errorf("json marshal error for contract abi\n")
		return err
	}
	//fmt.Printf("abi:[%v]", string(abi))
	diceActive, err := tools.GetABI(string(abi))
	//diceActive, err := tools.GetABI(tronTradeAbiJson)
	if nil != err {
		return err
	}
	tradeContract = diceActive.GetMethod("trade")
	if nil == tradeContract {
		return errors.New("tradeContract nil")
	}
	cancelContract = diceActive.GetMethod("cancel")
	if nil == cancelContract {
		return errors.New("cancelContract nil")
	}

	orderIDContract = diceActive.GetMethod("getOrderID")
	if nil == orderIDContract {
		return errors.New("orderIDContract nil")
	}
	getOrderInfoContract = diceActive.GetMethod("getOrderInfo")
	if nil == getOrderInfoContract {
		return errors.New("getOrderInfoContract nil")
	}
	// getDecimal = diceActive.GetMethod("decimals")
	// if nil == getDecimal {
	// 	return errors.New("decimals nil")
	// }

	return nil
}

func InitTradeSmart2(order int64) error {
	var smartContract *core.SmartContract
	defer util.CatchError()
	if order >= 1000000000000 {
		smartContract, _ = common.GetWalletClient().GetSmartContract(SmartContractAddr10)
		convertTrcToken(smartContract.Abi)
	} else {
		smartContract, _ = common.GetWalletClient().GetSmartContract(SmartContractAddr)
	}
	//fmt.Printf("ss:[%v]", smartContract.Abi.String())
	abi, err := json.Marshal(smartContract.Abi)
	if nil != err {
		fmt.Errorf("json marshal error for contract abi\n")
		return err
	}
	//fmt.Printf("abi:[%v]", string(abi))
	diceActive, err := tools.GetABI(string(abi))
	//diceActive, err := tools.GetABI(tronTradeAbiJson)
	if nil != err {
		return err
	}
	if order >= 1000000000000 {
		getOrderInfoContract10 = diceActive.GetMethod("getOrderInfo")
		if nil == getOrderInfoContract10 {
			return errors.New("getOrderInfoContract10 nil")
		}
	} else {
		tradeContract = diceActive.GetMethod("tradingOrder")
		if nil == tradeContract {
			return errors.New("tradeContract nil")
		}
		cancelContract = diceActive.GetMethod("cancelOrder")
		if nil == cancelContract {
			return errors.New("cancelContract nil")
		}

		orderIDContract = diceActive.GetMethod("orderID")
		if nil == orderIDContract {
			return errors.New("orderIDContract nil")
		}

		getOrderInfoContract = diceActive.GetMethod("getOrderInfo")
		if nil == getOrderInfoContract {
			return errors.New("getOrderInfoContract nil")
		}
	}
	// getDecimal = diceActive.GetMethod("decimals")
	// if nil == getDecimal {
	// 	return errors.New("decimals nil")
	// }

	return nil
}

//CallOrderInfo 调用获取智能合约订单信息接口
////查询订单信息
// function getOrderInfo(uint256 _orderID)
// return (orderInfo.orderID, orderInfo.status, orderInfo.orderType, orderInfo.user,
//         orderInfo.tokenA, orderInfo.amountA  挂单量, orderInfo.tokenB, orderInfo.amountB 理论交易额,
// orderInfo.orderFills 已经成交量, orderInfo.turnover  tokenB成交额);
func CallOrderInfo(orderID int64) (*SmartOrderInfo, error) {
	smartOrderInfo := &SmartOrderInfo{}
	defer util.CatchError()
	record := &CallRecord{}

	record.Data, record.Err = getOrderInfoContract.Pack(tools.GenAbiInt(orderID))
	if nil != record.Err {
		fmt.Errorf("pack orderID err:[%v]\n", record.Err)
		return smartOrderInfo, record.Err
	}

	result, err := tools.TriggerContract(SmartOwnerAddr, SmartContractAddr, 0, record.Data, SmartNode)

	if err != nil {
		return smartOrderInfo, err
	}
	// resultJson, err := json.Marshal(result)
	// if err != nil {
	// 	return smartOrderInfo, err
	// }
	//fmt.Printf("CallOrderInfo orderID[%v], result:%s\n", orderID, result.Result.Message)

	rawData := result.ConstantResult[0]
	paramsOut, _ := getOrderInfoContract.Outputs.UnpackValues(rawData)
	//	fmt.Printf("CallOrderInfo len(paramsOut):%v\n", len(paramsOut))
	if len(paramsOut) != 11 { //返回参数必须等于10，并且顺序不变
		fmt.Errorf("call smart contract function getOrderInfo() error:[return length ! =10]\n")
		return smartOrderInfo, errors.New("return length ! =10")
	}
	smartOrderInfo.OrderID = paramsOut[0].(*big.Int).Int64()
	smartOrderInfo.Status = paramsOut[1].(*big.Int).Int64()    //订单状态，0 进行中 1 完成全部 2 取消 4 结款
	smartOrderInfo.OrderType = paramsOut[2].(*big.Int).Int64() //order type 0 买  1 卖
	smartOrderInfo.OwnerAddress = tools.GetTronAddrFromAbiAddress(paramsOut[3])
	smartOrderInfo.TokenA = tools.GetTronAddrFromAbiAddress(paramsOut[4])
	smartOrderInfo.AmountA = paramsOut[5].(*big.Int).Int64()
	smartOrderInfo.TokenB = tools.GetTronAddrFromAbiAddress(paramsOut[6])
	smartOrderInfo.AmountB = paramsOut[7].(*big.Int).Int64()
	smartOrderInfo.OrderFills = paramsOut[8].(*big.Int).Int64()
	smartOrderInfo.Turnover = paramsOut[9].(*big.Int).Int64()
	smartOrderInfo.ChannelId = paramsOut[10].(*big.Int).Int64()

	fmt.Printf("CallOrderInfo orderID[%v], smartOrderInfo:[%#v]\n", orderID, smartOrderInfo)
	return smartOrderInfo, nil

}

//CallOrderInfo 调用获取智能合约订单信息接口
////查询订单信息
// function getOrderInfo(uint256 _orderID)
// return (orderInfo.orderID, orderInfo.status, orderInfo.orderType, orderInfo.user,
//         orderInfo.tokenA, orderInfo.amountA  挂单量, orderInfo.tokenB, orderInfo.amountB 理论交易额,
// orderInfo.orderFills 已经成交量, orderInfo.turnover  tokenB成交额);
func CallOrderInfo10(orderID int64) (*SmartOrderInfo10, error) {
	var err error
	var result *api.TransactionExtention
	smartOrderInfo10 := &SmartOrderInfo10{}
	defer util.CatchError()
	record := &CallRecord{}

	record.Data, record.Err = getOrderInfoContract10.Pack(tools.GenAbiInt(orderID))
	if nil != record.Err {
		fmt.Errorf("pack orderID err:[%v]\n", record.Err)
		return smartOrderInfo10, record.Err
	}
	result, err = tools.TriggerContract(SmartOwnerAddr10, SmartContractAddr10, 0, record.Data, SmartNode)

	if err != nil {
		return smartOrderInfo10, err
	}
	// resultJson, err := json.Marshal(result)
	// if err != nil {
	// 	return smartOrderInfo, err
	// }
	//fmt.Printf("CallOrderInfo orderID[%v], result:%s\n", orderID, result.Result.Message)

	rawData := result.ConstantResult[0]
	paramsOut, _ := getOrderInfoContract10.Outputs.UnpackValues(rawData)
	//	fmt.Printf("CallOrderInfo len(paramsOut):%v\n", len(paramsOut))
	if len(paramsOut) != 10 { //返回参数必须等于10，并且顺序不变
		fmt.Errorf("call smart contract function getOrderInfo() error:[return length ! =10]\n")
		return smartOrderInfo10, errors.New("return length ! =10")
	}
	smartOrderInfo10.OrderID = paramsOut[0].(*big.Int).Int64()
	smartOrderInfo10.Status = paramsOut[1].(*big.Int).Int64()    //订单状态，0 进行中 1 完成全部 2 取消 4 结款
	smartOrderInfo10.OrderType = paramsOut[2].(*big.Int).Int64() //order type 0 买  1 卖
	smartOrderInfo10.OwnerAddress = tools.GetTronAddrFromAbiAddress(paramsOut[3])
	//smartOrderInfo.TokenA = tools.GetTronAddrFromAbiAddress(paramsOut[4])
	smartOrderInfo10.TokenA = fmt.Sprintf("%v", paramsOut[4].(*big.Int).Int64())
	smartOrderInfo10.AmountA = paramsOut[5].(*big.Int).Int64()
	// smartOrderInfo.TokenB = tools.GetTronAddrFromAbiAddress(paramsOut[6])
	smartOrderInfo10.Price = paramsOut[6].(*big.Int).Int64()
	smartOrderInfo10.AmountB = paramsOut[7].(*big.Int).Int64()
	smartOrderInfo10.OrderFills = paramsOut[8].(*big.Int).Int64()
	smartOrderInfo10.Turnover = paramsOut[9].(*big.Int).Int64()

	fmt.Printf("CallOrderInfo10 orderID[%v], smartOrderInfo10:[%#v]\n", orderID, smartOrderInfo10)
	return smartOrderInfo10, nil

}

//CallOrderID 获取最新订单号方法
func CallOrderID() (int64, error) {
	defer util.CatchError()
	var nowOrderID int64
	record := &CallRecord{}
	record.Data, record.Err = orderIDContract.Pack()
	if nil != record.Err {
		fmt.Errorf("pack orderID err:[%v]\n", record.Err)
		return 0, record.Err
	}

	result, err := tools.TriggerContract(SmartOwnerAddr, SmartContractAddr, 0, record.Data, SmartNode)
	if err != nil {
		return 0, err
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		return 0, err
	}
	fmt.Printf("result:%v\n%s\n", string(resultJson), result.Result.Message)

	rawData := result.ConstantResult[0]
	paramsOut, _ := orderIDContract.Outputs.UnpackValues(rawData)
	fmt.Printf("%v--%T\n", paramsOut[0], paramsOut[0])
	if err != nil {
		fmt.Errorf("GenTriggerSmartContract orderID err:[%v]\n", err)
		return 0, err
	}
	nowOrderID = paramsOut[0].(*big.Int).Int64()
	return nowOrderID, nil
}

//SmartCancelOrder 调用合约-撤销订单  not use 前端直接调用
func SmartCancelOrder(orderID int64) (*CallRecord, error) {
	if nil == cancelContract {
		return nil, errors.New("withDrawContract nil")
	}

	record := &CallRecord{}
	record.Data, record.Err = cancelContract.Pack(tools.GenAbiInt(orderID))
	if nil != record.Err {
		return nil, record.Err
	}
	result, err := tools.TriggerContract(SmartOwnerAddr, SmartContractAddr, 0, record.Data, SmartNode)
	if err != nil {
		return nil, err
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	fmt.Printf("result:%v\n%s\n", string(resultJson), result.Result.Message)

	rawData := result.ConstantResult[0]
	paramsOut, _ := tradeContract.Outputs.UnpackValues(rawData)
	fmt.Printf("%v--%T\n", paramsOut[0], paramsOut[0])
	//number := paramsOut[0].(*big.Int)
	return nil, nil
}

//CallGetDecimal 获取合约账户精度
func CallGetDecimal(address string) (int64, error) {
	smartContract, err := common.GetWalletClient().GetSmartContract(address)
	//fmt.Printf("ss:[%v]", smartContract.Abi)
	//_ = smartContract
	abi, err := json.Marshal(smartContract.Abi)
	if nil != err {
		fmt.Errorf("json marshal error for contract abi\n")
		return 0, err
	}
	//fmt.Printf("abi:[%v]", string(abi))
	addressActive, err := tools.GetABI(string(abi))
	if nil != err {
		fmt.Errorf("get abi from string [%v] error:[%v]\n", string(abi), err)
		return 0, err
	}
	getDecimal := addressActive.GetMethod("decimals")
	if nil == getDecimal {
		fmt.Errorf("get dicimals() is nil\n")
		return 0, errors.New("decimals nil")
	}
	data, err := getDecimal.Pack()
	if nil != err {
		fmt.Errorf("get Pack() err:[%v]\n", err)
		return 0, err
	}
	result, err := tools.TriggerContract(address, address, 0, data, SmartNode)
	if err != nil {
		fmt.Errorf("TriggerContract err:[%v]\n", err)
		return 0, err
	}
	// resultJSON, _ := json.Marshal(result)
	// log.Debugf("result:%v\n%s\n", string(resultJSON), result.Result.Message)

	rawData := result.ConstantResult[0]
	paramsOut, _ := getDecimal.Outputs.UnpackValues(rawData)
	decimals := paramsOut[0].(uint8)
	return int64(decimals), nil
}
