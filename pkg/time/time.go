package time

import (
	"strconv"
	"strings"
	"time"
)

func MilliSecond(t time.Time) int64 {
	return t.UnixNano() / time.Millisecond.Nanoseconds()
}

func TimeFromMilliSecond(ms int64) time.Time {
	return time.Unix(0, ms*time.Millisecond.Nanoseconds())
}

func NowMS() int64 {
	return MilliSecond(time.Now())
}

// Getyyyy get the format of time yyyy for node uses
func Getyyyy(t time.Time) string {
	y, _, _ := t.Date()
	return formatWithZeroPadding(4, y)
}

// Getyyyymm get the format of time yyyymm for node uses
func Getyyyymm(t time.Time) string {
	y, m, _ := t.Date()
	return formatWithZeroPadding(4, y) + formatWithZeroPadding(2, int(m))
}

// Getyyyymmdd get the format of time yyyymmdd for node uses
func Getyyyymmdd(t time.Time) string {
	y, m, d := t.Date()
	return formatWithZeroPadding(4, y) + formatWithZeroPadding(2, int(m)) + formatWithZeroPadding(2, d)
}

func formatWithZeroPadding(l int, d int) string {
	s := strconv.Itoa(d)
	if len(s) < l {
		return strings.Repeat("0", l-len(s)) + s
	}
	return s
}
