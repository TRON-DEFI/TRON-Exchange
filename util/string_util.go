package util

import (
	"crypto/md5"
	"da3/system/apiSystem"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 生成随机字符串
func GetRandomString(n int) string {
	const symbols = "0123456789abcdefghjkmnopqrstuvwxyzABCDEFGHJKMNOPQRSTUVWXYZ"
	const symbolsIdxBits = 6 // symbols共58(111010b)个, 6bits可表示所有可能的Index
	const symbolsIdxMask = 1<<symbolsIdxBits - 1
	const symbolsIdxMax = 63 / symbolsIdxBits // number of symbol indices fitting in 63 bits

	prng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i, cache, remain := n-1, prng.Int63(), symbolsIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = prng.Int63(), symbolsIdxMax
		}
		if idx := int(cache & symbolsIdxMask); idx < len(symbols) {
			b[i] = symbols[idx]
			i--
		}
		cache >>= symbolsIdxBits
		remain--
	}
	return string(b)
}

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 校验邮箱格式
func IsMailFormat(mail string) bool {
	var mailRe = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)
	return mailRe.MatchString(mail)
}

func NowStr() string {
	tm := time.Unix(time.Now().Unix(), 0)
	return tm.Format("2006-01-02 15:04:05")
}

func IsEmptyStr(needCheck string) bool {
	if needCheck != "" && len(needCheck) > 0 {
		return false
	}
	return true
}

// ParsingJSONFromString 解析json结构：eg：
// {"BitTorrent":0,"Bithumb":0,"HuobiToken":0,"IPFS":0,"James":0,"MacCoin":0,"NBACoin":0,"Skypeople":0,"TRXTestCoin":0,"binance":0,"ofoBike":0}
func ParsingJSONFromString(jstr string) map[string]int64 {
	if jstr == "" {
		return nil
	}
	jsonMap := make(map[string]int64, 0)
	jsonStr := jstr[1 : len(jstr)-1] //去除前后{}
	for _, param := range strings.Split(jsonStr, ",") {
		if param != "" {
			for key, value := range strings.Split(param, ":") {
				if key == 0 {
					value = strings.Replace(value, "\"", "", -1)
					value = strings.Replace(value, "'", "", -1)
					jsonMap[value] = 0
				} else {
					ValueInt, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						ValueInt = 0
					}
					if ValueInt > 0 {
						jsonMap[value] = ValueInt
					}
				}

			}
		}
	}
	return jsonMap
}

//Distinct 清除重复的信息
func Distinct(value []string) (retData []string, distinct map[string]int) {
	retData = make([]string, 0, len(value)) //返回的去重后的信息
	distinct = make(map[string]int, len(value))

	for _, item := range value {
		if v, ok := distinct[item]; ok {
			v++
			distinct[item] = v
		} else {
			distinct[item] = 1
		}
	}

	for key := range distinct {
		retData = append(retData, key)
	}
	return retData, distinct
}

//SetDefaultVal 如果src为空，则返回defaultVal
func SetDefaultVal(src, defaultVal string) string {
	if src == "" {
		src = defaultVal
	}
	return src
}

//JSONObjectToString 将jason对象转换为string
func JSONObjectToString(v interface{}) (string, error) {
	//检查参数是否有效
	if nil == v {
		return "", apiSystem.ErrParam
	}

	var strJSONString string
	buffer, err := json.Marshal(v)
	if err == nil {
		strJSONString = string(buffer)
	}
	return strJSONString, err
}
