package util

import (
	"time"
)

/**字符串->时间对象*/
func Str2Time(formatTimeStr string) (time.Time, error) {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型
	return theTime, err

	//ts := theTime.Unix()
	//dt := time.Unix(ts, 0)
	//fmt.Printf("过期时间计算:%s\n", dt.Format("2006-01-02 15:04:05"))
}
