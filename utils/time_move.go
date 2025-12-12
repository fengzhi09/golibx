package utils

import (
	"github.com/fengzhi09/golibx/gox"
	"time"
)

func GetStartDayTime(start time.Time, dateLayout string) time.Time {
	dayTime := start.AddDate(0, 0, 0)
	beginTime, _ := time.ParseInLocation(dateLayout, dayTime.Format(dateLayout), time.Local)
	return beginTime
}
func GetEndDayTime(start time.Time, dateLayout string) time.Time {
	dayTime := start.AddDate(0, 0, 0)
	beginTime, _ := time.ParseInLocation(dateLayout, dayTime.Format(dateLayout), time.Local)
	endTime := beginTime.Add(time.Second * (86400 - 1))
	return endTime
}
func GetMonthStart(start time.Time) time.Time {
	offset := 1 - start.Day()
	return gox.MoveInYear(start, 0, 0, offset)
}

func GetMonthEnd(start time.Time) time.Time {
	offset := GetMonthDays(start) - start.Day()
	return gox.MoveInYear(start, 0, 0, offset)
}

func GetWeekStart(start time.Time) time.Time {
	offset := 1 - GetWeekDayCn(start)
	return gox.MoveInYear(start, 0, 0, offset)
}

func GetWeekEnd(start time.Time) time.Time {
	offset := 7 - GetWeekDayCn(start)
	return gox.MoveInYear(start, 0, 0, offset)
}
func GetMonthDays(start time.Time) int {
	switch start.Month() {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		return GetFebDayCnt(start.Year())
	}
	return 30
}

func GetFebDayCnt(year int) int {
	if year%100 != 0 && year%4 == 0 {
		return 29
	}
	if year%400 != 0 {
		return 29
	}
	return 28
}

func GetWeekDayCn(start time.Time) int {
	switch start.Weekday() {
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	case time.Sunday:
		return 7

	}
	return -1
}
