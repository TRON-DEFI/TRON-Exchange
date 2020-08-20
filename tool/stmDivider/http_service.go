package main

import (
	"github.com/adschain/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(NoCache)
	g.Use(Options)
	g.Use(Secure)
	g.Use(mw...)

	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The Incorrect API Route")
	})

	v1 := g.Group("/api/v0")
	{
		v1.GET("/reward/pool", getPoolReward)    //获取矿池余额
		v1.GET("/reward/query", getRewardList)   // 获取待发放奖励记录
		v1.GET("/reward/send", sendManualReward) //发放奖励
	}
	return g
}

type commonResp struct {
	Data  interface{} `json:"data"`
	Error error       `json:"err"`
}

func getRewardList(c *gin.Context) {
	results, err := queryPreparedAndSendReward()
	resp := &commonResp{
		Data:  results,
		Error: err,
	}
	c.JSON(http.StatusOK, resp)
}
func getPoolReward(c *gin.Context) {
	balance, err := getPoolTotalReward()
	resp := &commonResp{
		Data:  balance,
		Error: err,
	}
	c.JSON(http.StatusOK, resp)
}
func sendManualReward(c *gin.Context) {
	preparedList, err := queryPreparedAndSendReward()
	if err != nil {
		log.Infof("queryPreparedReward err:%v", err)
		c.JSON(http.StatusOK, err)
		return
	}
	results := sendRealReward(preparedList)
	err = UpdateRewardResult(results)
	resp := &commonResp{
		Data:  results,
		Error: err,
	}
	c.JSON(http.StatusOK, resp)
}
