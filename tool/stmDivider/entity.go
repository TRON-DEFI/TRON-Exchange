package main

import "github.com/tronprotocol/grpc-gateway/api"

//CallRecord 智能合约调用通用结构
type CallRecord struct {
	Owner     string // base58
	Contract  string // base58
	Method    string
	CallValue int64
	Data      []byte
	TrxHash   string // hex
	Err       error
	Return    *api.Return
}

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
//SmartUserInfo 合约中查询到的用户信息
type SmartUserInfo struct {
	UserAddr                 string `json:"userAddr"`                    //: "99500000000",
	Balance                  int64  `json:"balance"`                     //: "0",
	PledgeAmount             int64  `json:"pledgeAmount"`                //: "1",
	LockAmount               int64  `json:"lockAmount"`                  //: "99500000",
	StakingAmount            int64  `json:"stakingAmount"`               //: "100000000000",
	UnstakingAmount          int64  `json:"unstakingAmount"`             //: "100000000",
	Rate                     int64  `json:"rate"`                        //: "1",
	InviteeAmount            int64  `json:"invitee_amount"`              // "0",
	InviteeCount             int64  `json:"invitee_count"`               //: "2"
	StakingOtkenWithdrawTime int64  `json:"staking_otken_withdraw_time"` //: "2"
	ParentID                 int64  `json:"parentID"`                    //: "2"
}

type UserRewardInfo struct {
	ID           int64  `json:"id"`          //: 用户地址,
	UserAddr     string `json:"userAddr"`    //: 用户地址,
	Pledged      int64  `json:"pledge"`      //: 用户抵押量,
	TotalPledged int64  `json:"totalPledge"` //当前总抵押量
	TotalBonus   int64  `json:"totalBonus"`  //: 当前奖池金额,
	Bonus        int64  `json:"bonus"`       //: 用户可分得金额,
	TrxHash      string `json:"hash"`        //: 转账交易hash,
	TrxResult    string `json:"result"`      //: 转账交易结果,
}

// 矿池信息  stm_pool
type DividePoolInfo struct {
	ID            int64  `json:"id"`             //: 用户地址,
	CallerAddress string `json:"caller_address"` //: 用户地址,
	PoolAddress   string `json:"pool_address"`   //: 用户抵押量,
	TotalAmount   int64  `json:"total_amount"`   //当前总抵押量
	TrxHash       string `json:"trx_hash"`       //: 当前奖池金额,
	TrxResult     string `json:"trx_result"`     //: 用户可分得金额,
	CTime         string `json:"c_time"`         //: 转账交易hash,
	MTime         string `json:"m_time"`         //: 转账交易结果,
}
