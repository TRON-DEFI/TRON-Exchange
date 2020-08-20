package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tradehome-service/common/errno"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, msg := errno.DecodeErr(err)
	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
