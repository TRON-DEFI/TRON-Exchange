package main

import (
	"crypto/tls"
	"errors"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/config"
	"github.com/wlcy/tradehome-service/router"
	"github.com/wlcy/tradehome-service/task"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var (
	cfg = pflag.StringP("config", "c", "", "server config file path")
)

func main() {
	pflag.Parse()

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	gin.SetMode(viper.GetString("run_mode"))

	g := gin.New()

	router.Load(
		g,
		// log record
		//middleware.Logging(),
	)

	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("the router has no response, or it might took too long to start up", err)
		}
		log.Info("the router has been deployed successfully")
	}()

	go task.Async()

	log.Infof("start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())

}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(viper.GetString("url") + "/api/check/health")

		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		log.Info("Waiting for the router, retry in 1 second")
		time.Sleep(time.Second)
	}
	return errors.New("connect to the router error")
}
