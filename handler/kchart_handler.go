package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tradehome-service/common/errno"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/service"
)

// QueryKChartData ...
// @Summary Query KChart Data
// @Tags Exchange
// @Accept  json
// @Produce  json
// @Param exchange_id query string true "Exchange pair id"
// @Param time_start query string true "The start time of data"
// @Param time_end query string true "The end time of data"
// @Param granuarity query string true "The granuarity of k chart, eg: 1min/5min/1day..."
// @Success 200 {object} transcanapientity.KChartResp "KChart data"
// @Router /api/exchange/kchart [get]
func QueryKChartData(c *gin.Context) {
	req := &entity.KChartQueryParam{}
	req.ExchangeID = c.Query("exchangeId")
	req.TimeStart = c.Query("startTime")
	req.TimeEnd = c.Query("endTime")
	req.Granu = c.Query("granu")
	// key := GetQueryKey("QueryKChartData", req)
	// data := LoadBuffer(key, 5)
	// if nil != data {
	// 	c.JSON(http.StatusOK, data)
	// 	return
	// }
	if "" == req.ExchangeID || "" == req.TimeStart || "" == req.TimeEnd || "" == req.Granu {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	resp, err := service.GetKChartData(req)
	if err != nil {
		SendResponse(c, errno.InternalServerError, nil)
		return
	}
	SendResponse(c, nil, resp)
}
