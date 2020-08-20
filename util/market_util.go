package util

import (
	"fmt"
	"github.com/adschain/log"
	"math"
	"time"
)

//ExchangeFactor exchange pair price stored to DB and got from DB's factor 10E8
const ExchangeFactor int64 = 100000000

//CalcPrice Calculate the price of first token according to the balance of the two tokens,
//Note that the price is mutiplied by 100000000 (10E8)for store into DB.
func CalcPrice(fristTokenBalance int64, secondTokenBalance int64) int64 {
	if 0 == fristTokenBalance || 0 == secondTokenBalance {
		return 0
	}
	price := int64((float64(secondTokenBalance) / float64(fristTokenBalance)) * float64(ExchangeFactor))
	return price
}

//CalcQuant Calculate the quantity of the other token according to bancor protocol
func CalcQuant(sellTokenQuant int64, sellTokenBalance int64, theOtherTokenBalance int64) int64 {
	//supply is 1E18 at present
	supply := int64(1000000000000000000)
	supplyQuant := int64(float64(-supply) * float64(1.0-math.Pow(1.0+float64(sellTokenQuant)/float64(sellTokenBalance+sellTokenQuant), 0.0005)))
	buyTokenQuant := int64(float64(theOtherTokenBalance) * float64(math.Pow(1.0+float64(supplyQuant)/float64(supply), 2000.0)-1.0))
	return buyTokenQuant
}

//Abs 返回int64的绝对值
func Abs(num int64) int64 {
	if num < 0 {
		return -num
	}
	return num
}

//GetBenchmarkTimeAndSpan find a proper time point for querying data from kgraph database as the start time
//@param isForward forward means start time later than tempTime, backward means start time before than tempTime
//@param tempTime is the time given by front end
//@param granu is the granularity of the kgraph, that is 1 minute/5 minutes/15 minutes...
func GetBenchmarkTimeAndSpan(granu string, tempTime time.Time, isForward bool) (int, time.Time, error) {
	log.Debugf("Begin getBenchmarkTimeAndSpan Granu:[%v], tempTime:[%v]", granu, tempTime.UTC().Format(DATETIMEFORMAT))
	var span int
	if "1min" == granu {
		span = 1
		tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute(), 0, 0, time.UTC)
	} else if "5min" == granu {
		span = 5
		r := tempTime.Minute() % span
		if isForward {
			if 0 == r {
				tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute(), 0, 0, time.UTC)
			} else {
				tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute()-r+span, 0, 0, time.UTC)
			}
		} else {
			tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute()-r, 0, 0, time.UTC)
		}
	} else if "15min" == granu {
		span = 15
		r := tempTime.Minute() % span
		if isForward {
			if 0 == r {
				tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute(), 0, 0, time.UTC)
			} else {
				tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute()-r+span, 0, 0, time.UTC)
			}
		} else {
			tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute()-r, 0, 0, time.UTC)
		}
	} else if "30min" == granu {
		span = 30
		r := tempTime.Minute() % span
		if isForward {
			if 0 == r {
				tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute(), 0, 0, time.UTC)
			} else {
				tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute()-r+span, 0, 0, time.UTC)
			}
		} else {
			tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), tempTime.Minute()-r, 0, 0, time.UTC)
		}
	} else if "1h" == granu {
		span = 60
		aTime := time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), 0, 0, 0, time.UTC)
		if aTime.Equal(tempTime) {
			tempTime = aTime
		} else {
			if isForward {
				tempTime = aTime.Add(time.Hour)
			} else {
				tempTime = aTime
			}
		}
	} else if "4h" == granu {
		span = 240
		r := tempTime.Hour() % 4
		if isForward { // get integral hour forward
			if 0 == r {
				aTime := time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour(), 0, 0, 0, time.UTC)
				if aTime.Equal(tempTime) { // right integral hour and can be devided by 4 (eg.16:00)
					tempTime = aTime
				} else { // not integral hour and hours can be devided by 4(eg. 16:20)
					tempTime = aTime.Add(4 * time.Hour)
				}
			} else { //hours cannot devided by 4, then subtract the remainder and add 4 hours
				tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour()-r+4, 0, 0, 0, time.UTC)
			}
		} else { // get inegral hour backward
			tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), tempTime.Hour()-r, 0, 0, 0, time.UTC)
		}
	} else if "1d" == granu {
		span = 1
		tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), 0, 0, 0, 0, time.UTC)
	} else if "5d" == granu {
		span = 5
		tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), 0, 0, 0, 0, time.UTC)
	} else if "1w" == granu {
		span = 7
		weekDayNO := int(tempTime.Weekday())
		if 0 == weekDayNO {
			tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), 0, 0, 0, 0, time.UTC)
		} else {
			tempTime = time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day()-int(weekDayNO)+span, 0, 0, 0, 0, time.UTC)
		}
	} else if "1m" == granu {
		span = 31
		log.Debugf("tempTime is: %v\n", tempTime.Format(DATEFORMAT))
		aTime := time.Date(tempTime.Year(), tempTime.Month(), 1, 0, 0, 0, 0, time.UTC)
		log.Debugf("aTime is: %v\n", aTime.Format(DATEFORMAT))
		if aTime.Equal(tempTime) {
			tempTime = aTime
		} else {
			tempTime = time.Date(tempTime.Year(), tempTime.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		}
		log.Debugf("tempTime is: %v\n", tempTime.Format(DATEFORMAT))
	} else {
		log.Errorf(fmt.Errorf("granularity is not right"),"Granularity is not right:[%v]", granu,)
		return 0, time.Now().UTC(), fmt.Errorf("granularity is not right")
	}
	log.Debugf("End findBenchmarkMin benchmarkTime:[%v], span: [%v]\n", tempTime.UTC().Format(DATETIMEFORMAT), span)
	return span, tempTime.UTC(), nil
}
