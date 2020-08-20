package entity

type MineStatistic struct {
	Id                 int64  `json:"id"`
	PairId             int64  `json:"pairId"`
	OrderId            int64  `json:"orderId"`
	OrderType          int64  `json:"orderType"`
	UserAddress        string `json:"userAddress"`
	FirstTokenAddress  string `json:"firstTokenAddress"`
	SecondTokenAddress string `json:"secondTokenAddress"`
	TradeVolume        int64  `json:"tradeVolume"`
	RewardVolume       int64  `json:"rewardVolume"`
	status             int    `json:"status"`
	OrderCreateTime    int64  `json:"orderCreateTime"`
}

type MineReward struct {
	UserAddress string `json:"userAddress"`
	TotalReward int64  `json:"totalReward"`
}
