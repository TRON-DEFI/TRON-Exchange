package handler

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-engin/core/utils"
	"sync"
	"time"
)

var queryBuf sync.Map

type ifData struct {
	ts   time.Time
	Data interface{}
}

// GetQueryKey ....
func GetQueryKey(ifName string, req interface{}) string {
	val, ok := req.(int64)
	if ok {
		return fmt.Sprintf("%v:%v", ifName, val)
	}
	val2, ok := req.(uint64)
	if ok {
		return fmt.Sprintf("%v:%v", ifName, val2)
	}
	val1, ok := req.(string)
	if ok {
		return fmt.Sprintf("%v:%v", ifName, val1)
	}
	return fmt.Sprintf("%v:%v", ifName, utils.ToJSONStr(req))
}

// LoadBuffer ...
func LoadBuffer(key string, ts int) interface{} {
	ret, ok := queryBuf.Load(key)
	log.Debugf("LoadBuffer key=[%v],ret=[%v], ok = [%v]\n", key, ret, ok)
	if ok {
		data, ok := ret.(*ifData)
		if ok && nil != data && time.Since(data.ts) <= time.Duration(ts)*time.Second {
			return data.Data
		}
	}
	return nil
}

// StoreBuffer ...
func StoreBuffer(key string, data interface{}) {
	queryBuf.Store(key, &ifData{time.Now(), data})
}
