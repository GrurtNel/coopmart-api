package utils

import (
	"strconv"
	"time"
)

func Float32ToPercentString(num float32) string {
	num = num * 100
	return strconv.FormatFloat(float64(num), 'f', 2, 32) + "%"
}
func UnixToDate(unix int64) string {
	tm := time.Unix(unix/1000, 0)
	return tm.Format("02-01-2006 15:04:05")
}
func FormatFloatPoint(num float32) float32 {
	var temp = strconv.FormatFloat(float64(num), 'f', 2, 32)
	var res, _ = strconv.ParseFloat(temp, 32)
	return float32(res)
}

func FormatDecimalPoint(num float32) float32 {
	var temp = strconv.FormatFloat(float64(num*10), 'f', 2, 32)
	var res, _ = strconv.ParseFloat(temp, 32)
	return float32(res)
}
