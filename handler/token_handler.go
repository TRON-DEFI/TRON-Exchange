package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wlcy/tradehome-service/common/errno"
	"github.com/wlcy/tradehome-service/service"
)

func QueryTokenInfo(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	data, err := service.QueryTokenInfo(address)
	if err != nil {
		SendResponse(c, errno.InternalServerError, nil)
		return
	}
	SendResponse(c, nil, data)
}

func GetTronPrice(c *gin.Context) {
	logrus.Debugf("Hello /market/tron/price")
	//resp, err := service.QueryMarkets()
	resp, err := service.QueryMarketsBuffer()
	if err != nil {
		SendResponse(c, errno.InternalServerError, err)
	}
	SendResponse(c, nil, resp)
}
