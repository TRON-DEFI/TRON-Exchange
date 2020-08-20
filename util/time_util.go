package util

import "time"

func TimeStampBefore24hV1() string {
	now := time.Now().UTC()
	d, _ := time.ParseDuration("-24h")
	d24h := now.Add(d).UTC()
	return d24h.Format("2006-01-02 15:04:05")
}

func TimeStampBefore24h() int64 {
	now := time.Now().UTC()
	d, _ := time.ParseDuration("-24h")
	d24h := now.Add(d).UTC()
	return d24h.UnixNano() / 1000000
}
