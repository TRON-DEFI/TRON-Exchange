package util

import (
	"fmt"
	"testing"
)

func TestGetRandomString(t *testing.T) {
	fmt.Printf("GetRandomString result:%s\n", GetRandomString(16))
}

func TestMD5(t *testing.T) {
	fmt.Printf("MD5 result:%s\n", MD5("1g1hlg42"+"#123@HGD"))
}

func TestIsMailFormat(t *testing.T) {
	fmt.Printf("isMailFormat result:%v", IsMailFormat("6767654@qq.com"))
}
