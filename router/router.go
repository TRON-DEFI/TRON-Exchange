package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tradehome-service/handler"
	"github.com/wlcy/tradehome-service/router/middleware"
	"net/http"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The Incorrect API Route")
	})

	// the health check handler
	check := g.Group("/api/check")
	{
		check.GET("/health", handler.HealthCheck)
		check.GET("/disk", handler.DiskCheck)
		check.GET("/cpu", handler.CPUCheck)
		check.GET("/ram", handler.RAMCheck)
	}

	v1 := g.Group("/api/v1")
	{
		//K线图查询
		v1.GET("/market/kline/query", handler.QueryKChartData)

		//查询交易对列表
		v1.GET("/market/pair/query", handler.MarketPairList)

		//用户交易
		v1.GET("/market/user/order", handler.UserOrderList)
		//添加订单对应渠道ID 不对外暴露
		//v1.POST("/market/channel/id/add", handler.AddChannelId)
		//最新交易记录展示
		v1.GET("/market/common/order/latest", handler.LatestOrderList)
		//根据交易对ID查询盘口信息
		//v1.GET("/api/exchange/common/orderList/:pairID", handler.GetOrderList)
		//根据交易对ID查询盘口信息2
		v1.GET("/market/common/order/list/:pairID", handler.GetOrderList2)
		//查询撮合池所有订单
		v1.GET("/market/common/order/all", handler.GetAllOrderList)
		//查询撮合池所有订单2
		//v1.GET("/market/common/order/all2", handler.GetAllOrderList2)
		//调用智能合约orderInfo
		v1.GET("/smart/order/:orderID", handler.GetSmartOrderInfo)
		//v1.GET("/smart/order2/:orderID", handler.GetSmartOrderInfo2)
		v1.GET("/market/price/:exchangeID", handler.GetTopPriceByExchangeID)
		//v1.GET("/market/price/:exchangeID", handler.GetTopPriceByExchangeID2)

		v1.GET("/market/common/tokenInfo/query", handler.QueryTokenInfo)

		v1.GET("/market/tron/price", handler.GetTronPrice)

	}

	return g
}
