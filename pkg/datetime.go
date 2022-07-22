package pkg

import (
	"time"
)

var (
	TimeLayout  = "2006-01-02 15:04:05"
	SingaporeTz = time.FixedZone("Singapore", 3600*8)
)

// 获取系统当前时间(秒)
func GetSecUnixTime() int64 {
	return time.Now().In(SingaporeTz).Unix()
}

// 获取系统当前时间(毫秒)
func GetMsecUnixTime() int64 {
	return time.Now().In(SingaporeTz).UnixNano() / 1e6
}

// 获取系统当前天0点时间
func GetDayZeroTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).In(SingaporeTz)
}

// 获取系统当前月第一天0点时间
func GetMonthFirstDayZeroTime() time.Time {
	now := time.Now()
	firstDateTime := now.AddDate(0, 0, -now.Day()+1)
	return time.Date(firstDateTime.Year(), firstDateTime.Month(), firstDateTime.Day(), 0, 0, 0, 0, firstDateTime.Location()).In(SingaporeTz)
}

// 获取系统当前年第一天0点时间
func GetYearFirstDayZeroTime() time.Time {
	now := time.Now()
	currentYear, _, _ := now.Date()
	currentLocation := now.Location()
	return time.Date(currentYear, time.January, 1, 0, 0, 0, 0, currentLocation).In(SingaporeTz)
}

// 获取系统当前时间(纳秒)
func GetNanoUnixTime() int64 {
	return time.Now().In(SingaporeTz).UnixNano()
}

// 年月日
func GetUnixDay() string {
	return time.Now().In(SingaporeTz).Format("2006-01-02")
}
func GetUnixDayTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, SingaporeTz)
}

// 时间戳转日期
func TimestampToDate(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}

	return time.Unix(timestamp, 0).In(SingaporeTz).Format(TimeLayout)
}

// 带时区的当前时间戳(秒)
func Now() time.Time {
	return time.Now().In(SingaporeTz)
}

// 带时区的当前时间戳(毫秒)
func UnixMilli(milliUnix int64) time.Time {
	return time.Unix(0, milliUnix*1e6).In(SingaporeTz)
}

// 带时区的当前时间戳(纳秒)
func UnixNanoUnix(nanoUnix int64) time.Time {
	return time.Unix(0, nanoUnix).In(SingaporeTz)
}
