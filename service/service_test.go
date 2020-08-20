package service

import (
	"fmt"
	"github.com/wlcy/tradehome-service/buffer"
	"github.com/wlcy/tradehome-service/util"
	"testing"
)

func TestGetPrice(t *testing.T) {
	tt := buffer.GetMarketBuffer().GetMarket()
	str, _ := util.JSONObjectToString(tt)
	fmt.Println(str)
}
