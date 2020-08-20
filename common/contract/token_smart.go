package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wlcy/tradehome-engin/core/tools"
	"github.com/wlcy/tradehome-engin/web/common"
	"github.com/wlcy/tradehome-service/util"
)

var _abiEntry *tools.AbiEntry
var _transfer *tools.Method

var MineRate int64
var MineContractAddress, MineOwnerAddr, MinePrivateKey string

func InitMineSmart() error {
	defer util.CatchError()
	smartContract, err := common.GetWalletClient().GetSmartContract(MineContractAddress)
	abi, err := json.Marshal(smartContract.Abi)
	if nil != err {
		fmt.Errorf("json marshal error for contract abi\n")
		return err
	}
	_abiEntry, err := tools.GetABI(string(abi))
	if nil != err {
		return err
	}
	_transfer = _abiEntry.GetMethod("transfer")
	if nil == _transfer {
		return errors.New("transfer nil")
	}

	return nil
}

func TransferAsset2(toAddr string, amount int64) (*CallRecord, error) {
	privateKey := MinePrivateKey
	feeLimit := int64(10000000)
	if nil == _transfer {
		return nil, errors.New("_transfer nil")
	}

	record := &CallRecord{}
	record.Data, record.Err = _transfer.Pack(tools.GenAbiAddress(toAddr), tools.GenAbiInt(amount))
	if nil != record.Err {
		return nil, record.Err
	}
	ctxType, ctx, err := tools.GenTriggerSmartContract(MineOwnerAddr, MineContractAddress, record.CallValue, record.Data)
	if nil != err {
		return nil, record.Err
	}
	record.TrxHash, record.Return, record.Err = tools.BroadcastCtxWithFeeLimit(ctxType, ctx, nil, privateKey, 0, feeLimit)
	if nil != record.Err {
		return nil, record.Err
	}

	return record, nil
}
