package gox

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	IntDateMin  int64 = 20000101       // 2000-01-01
	IntDateMax  int64 = 22010101       // 2201-01-01
	UtsSecMin   int64 = 946656000      // time.Date(2000,01,01,0,0,0,0).Unix()
	UtsSecMax   int64 = 7289625600     // time.Date(2201,01,01,0,0,0,0).Unix()
	UtsMsMin    int64 = 946656000000   // time.Date(2000,01,01,0,0,0,0).UnixNano()/1e6
	UtsMsMax    int64 = 7258089600000  // time.Date(2201,01,01,0,0,0,0).UnixNano()/1e6
	LongDateMin int64 = 20000101000000 // 2000-01-01 00:00:00
	LongDateMax int64 = 22010101000000 // 2201-01-01 00:00:00
)

const (
	ISODateTimeNsTZF  = "2006-01-02T15:04:05.999999999Z07:00"
	ISODateTimeNsTZF8 = "2006-01-02T15:04:05.99999999Z07:00"
	ISODateTimeNsTZF7 = "2006-01-02T15:04:05.9999999Z07:00"
	ISODateTimeNsTZF6 = "2006-01-02T15:04:05.999999Z07:00"
	ISODateTimeNsTZF5 = "2006-01-02T15:04:05.99999Z07:00"
	RFC1123Z          = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC850            = "Monday, 02-Jan-06 15:04:05 MST"  // 变长最后考虑
	RFC1123           = "Mon, 02 Jan 2006 15:04:05 MST"
	RubyDateZ         = "Mon Jan 02 15:04:05 -0700 2006"
	ISODateTimeNsTZ   = "2006-01-02T15:04:05.999999999Z"
	ISODateTimeNsZ    = "2006-01-02 15:04:05.999999999Z"
	ISODateTimeNsT    = "2006-01-02T15:04:05.999999999"
	ISODateTimeNs     = "2006-01-02 15:04:05.999999999"
	ISODateTimeMsTZF  = "2006-01-02T15:04:05.999Z07:00"
	UnixDateZ         = "Mon Jan _2 15:04:05 MST 2006"
	ISODateTimeMs2TZF = "2006-01-02T15:04:05.99Z07:00"
	ISODateTimeMs1TZF = "2006-01-02T15:04:05.9Z07:00"
	RFC3339Z          = "2006-01-02T15:04:05Z07:00"
	ANSIC             = "Mon Jan _2 15:04:05 2006"
	ISODateTimeMsTZ   = "2006-01-02T15:04:05.999Z"
	ISODateTimeMsZ    = "2006-01-02 15:04:05.999Z"
	ISODateTimeMsT    = "2006-01-02T15:04:05.999"
	ISODateTimeMs     = "2006-01-02 15:04:05.999"
	ISODateTimeMs2TZ  = "2006-01-02T15:04:05.99Z"
	ISODateTimeMs2Z   = "2006-01-02 15:04:05.99Z"
	ISODateTimeMs2T   = "2006-01-02T15:04:05.99"
	ISODateTimeMs2    = "2006-01-02 15:04:05.99"
	ISODateTimeMs1TZ  = "2006-01-02T15:04:05.9Z"
	ISODateTimeMs1Z   = "2006-01-02 15:04:05.9Z"
	RFC822Z           = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	ISODateTimeMs1T   = "2006-01-02T15:04:05.9"
	ISODateTimeMs1    = "2006-01-02 15:04:05.9"
	RFC822            = "02 Jan 06 15:04 MST"
	ISODateTime       = "2006-01-02 15:04:05"
	YYYYMMDDHHMMSS    = "20060102150405"
	ISOTimeMs         = "15:04:05.999"
	ISODate           = "2006-01-02"
	ISOTime           = "15:04:05"
	YYYYMMDD          = "20060102"
	HHMMSS            = "150405"
)

type Formatter struct {
	Name    string
	Format  string
	Can     func(s string) bool
	HasZone bool
}

var (
	numbers     = []string{"0", "1", "3", "4", "5", "6", "7", "8", "9"}
	zoneSymbols = []string{"+", "-", "Z"}
)

// var zoneWords = []string{"+", "-", "Z"}

var Formatters = []Formatter{
	{"ISODateTimeNsTZF", ISODateTimeNsTZF, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNsTZF) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 29, zoneSymbols...)
	}, true},
	{"ISODateTimeNsTZF8", ISODateTimeNsTZF8, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNsTZF8) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 28, zoneSymbols...)
	}, true},
	{"ISODateTimeNsTZF7", ISODateTimeNsTZF7, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNsTZF7) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 27, zoneSymbols...)
	}, true},
	{"ISODateTimeNsTZF6", ISODateTimeNsTZF6, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNsTZF6) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 26, zoneSymbols...)
	}, true},
	{"ISODateTimeNsTZF5", ISODateTimeNsTZF5, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNsTZF5) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 25, zoneSymbols...)
	}, true},
	{"RFC1123Z", RFC1123Z, func(dts string) bool {
		return len(dts) >= len(RFC1123Z)
	}, true},
	{"RFC1123", RFC1123, func(dts string) bool {
		return len(dts) >= len(RFC1123)
	}, true},
	{"RubyDateZ", RubyDateZ, func(dts string) bool {
		return len(dts) >= len(RubyDateZ) && StrAtMatch(dts, 3, " ")
	}, true},
	{"ISODateTimeNsTZ", ISODateTimeNsTZ, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNsTZ) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 29, zoneSymbols...)
	}, true},
	{"ISODateTimeNsZ", ISODateTimeNsZ, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNsZ) && StrAtMatch(dts, 29, zoneSymbols...)
	}, true},
	{"ISODateTimeNsT", ISODateTimeNsT, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNsT) && StrAtMatch(dts, 10, "T")
	}, false},
	{"ISODateTimeNs", ISODateTimeNs, func(dts string) bool {
		return len(dts) >= len(ISODateTimeNs) && StrAtMatch(dts, 10, " ")
	}, false},
	{"ISODateTimeMsTZF", ISODateTimeMsTZF, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMsTZF) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 23, zoneSymbols...)
	}, true},
	{"UnixDateZ", UnixDateZ, func(dts string) bool {
		return len(dts) >= len(UnixDateZ) && StrAtMatch(dts, 3, " ")
	}, true},
	{"ISODateTimeMs2TZF", ISODateTimeMs2TZF, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs2TZF) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 22, zoneSymbols...)
	}, true},
	{"ISODateTimeMs1TZF", ISODateTimeMs1TZF, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs1TZF) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 21, zoneSymbols...)
	}, true},
	{"RFC3339Z", RFC3339Z, func(dts string) bool {
		return len(dts) >= len(RFC3339Z)
	}, true},
	{"ANSIC", ANSIC, func(dts string) bool {
		return len(dts) >= len(ANSIC) && strings.Count(dts, " ") >= 3
	}, false},
	{"ISODateTimeMsTZ", ISODateTimeMsTZ, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMsTZ) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 23, zoneSymbols...)
	}, true},
	{"ISODateTimeMsZ", ISODateTimeMsZ, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMsZ) && StrAtMatch(dts, 23, zoneSymbols...)
	}, true},
	{"ISODateTimeMsT", ISODateTimeMsT, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMsT) && StrAtMatch(dts, 10, "T")
	}, true},
	{"ISODateTimeMs", ISODateTimeMs, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs) && StrAtMatch(dts, 10, " ")
	}, true},
	{"ISODateTimeMs2TZ", ISODateTimeMs2TZ, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs2TZ) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 22, zoneSymbols...)
	}, true},
	{"ISODateTimeMs2Z", ISODateTimeMs2Z, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs2Z) && StrAtMatch(dts, 22, zoneSymbols...)
	}, true},
	{"ISODateTimeMs2T", ISODateTimeMs2T, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs2T) && StrAtMatch(dts, 10, "T")
	}, false},
	{"ISODateTimeMs2", ISODateTimeMs2, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs2) && StrAtMatch(dts, 10, " ") && StrAtMatch(dts, 21, numbers...)
	}, false},
	{"ISODateTimeMs1TZ", ISODateTimeMs1TZ, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs1Z) && StrAtMatch(dts, 10, "T") && StrAtMatch(dts, 21, zoneSymbols...)
	}, true},
	{"ISODateTimeMs1Z", ISODateTimeMs1Z, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs1Z) && StrAtMatch(dts, 10, " ") && StrAtMatch(dts, 21, zoneSymbols...)
	}, true},
	{"RFC822Z", RFC822Z, func(dts string) bool {
		return len(dts) >= len(RFC822Z) && strings.Count(dts, " ") >= 3 && strings.Contains(dts, "Z")
	}, true},
	{"ISODateTimeMs1T", ISODateTimeMs1T, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs1T) && StrAtMatch(dts, 10, "T")
	}, true},
	{"ISODateTimeMs1", ISODateTimeMs1, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs1) && StrAtMatch(dts, 10, " ")
	}, true},
	{"ISODateTimeMs1", ISODateTimeMs1, func(dts string) bool {
		return len(dts) >= len(ISODateTimeMs1) && strings.Count(dts, " ") == 1
	}, true},
	{"RFC822", RFC822, func(dts string) bool {
		return len(dts) >= len(RFC822) && strings.Count(dts, " ") == 3
	}, true},
	{"ISODateTime", ISODateTime, func(dts string) bool {
		return len(dts) >= len(ISODateTime) && strings.Count(dts, " ") == 1
	}, true},
	{"YYYYMMDDHHMMSS", YYYYMMDDHHMMSS, func(dts string) bool {
		return len(dts) >= len(YYYYMMDDHHMMSS)
	}, true},
	{"ISODate", ISODate, func(dts string) bool {
		return len(dts) >= len(ISODate)
	}, false},
	{"ISOTimeMs", ISOTimeMs, func(dts string) bool {
		return len(dts) >= len(ISOTimeMs)
	}, false},
	{"YYYYMMDD", YYYYMMDD, func(dts string) bool {
		return len(dts) >= len(YYYYMMDD) && !strings.Contains(dts, ":")
	}, false},
	{"ISOTime", ISOTime, func(dts string) bool {
		return len(dts) >= len(ISOTime) && strings.Contains(dts, ":")
	}, false},
	{"HHMMSS", HHMMSS, func(dts string) bool {
		return len(dts) >= len(HHMMSS)
	}, false},
}

var ExFormatters = []Formatter{
	{"RFC850", RFC850, func(dts string) bool {
		return len(dts) >= len(RFC850)
	}, true},
}

type Zone = time.Location

func ZoneCn() *Zone {
	return time.FixedZone("cn", 8*3600)
}

func ZoneLocal() *Zone {
	return time.Local
}

func ZoneUTC() *Zone {
	return time.UTC
}

func Now() time.Time {
	return time.Now()
}

func NowStr() string {
	return Now().Format(ISODateTime)
}

func NowMsStr() string {
	return Now().Format(ISODateTimeMs)
}

func FormatUTCZone(dt time.Time, format string) string {
	return TimeFormatLoc(dt, format, ZoneUTC())
}

func FormatLocalZone(dt time.Time, format string) string {
	return TimeFormatLoc(dt, format, ZoneLocal())
}

func FormatCnZone(dt time.Time, format string) string {
	return TimeFormatLoc(dt, format, ZoneCn())
}

func TimeFormatLoc(dt time.Time, format string, loc *Zone) string {
	if loc != nil {
		dt = dt.In(loc)
	}
	return dt.Format(format)
}

func UniformDt(dts string) (time.Time, error) {
	loc := time.Local
	dts = strings.ReplaceAll(dts, "/", "-")
	//if StrAtMatch(dts, 10, "T") && strings.Contains(dts, "Z") {
	//	/*
	//		"2006-01-02T15:04:05.999Z"
	//		"2006-01-02T15:04:05.99Z"
	//		"2006-01-02T15:04:05.9Z"
	//		"2006-01-02T15:04:05Z"
	//	*/
	//	loc = time.UTC
	//	dts = strings.ReplaceAll(dts, "T", " ")
	//	dts = strings.ReplaceAll(dts, "Z", "")
	//}

	long := AsLong(dts)
	if IntDateMin <= long && long < IntDateMax {
		year, month, date := int(long/10000), int(long/100%100), int(long%100)
		return time.Date(year, time.Month(month), date, 0, 0, 0, 0, loc), nil
	}
	if UtsSecMin <= long && long < UtsSecMax {
		return FromUnixMs(long * 1000), nil
	}
	if UtsMsMin <= long && long < UtsMsMax {
		return FromUnixMs(long), nil
	}
	if LongDateMin <= long && long < LongDateMax {
		long1, long2 := long/1e6, long%1e6
		year, month, date := int(long1/10000), int(long1/100%100), int(long1%100)
		hour, minute, second := int(long2/10000), int(long2/100%100), int(long2%100)
		return time.Date(year, time.Month(month), date, hour, minute, second, 0, loc), nil
	}
	firstName, firstErr := "", ""
	for _, _formatter := range Formatters {
		tmp := SubStr(dts, 0, len(_formatter.Format))
		var (
			ti  time.Time
			err error
		)
		if _formatter.HasZone {
			ti, err = time.Parse(_formatter.Format, tmp)
		} else {
			ti, err = time.ParseInLocation(_formatter.Format, tmp, loc)
		}
		if err != nil {
			if _formatter.Can(dts) && firstName == "" {
				firstName, firstErr = _formatter.Name, err.Error()
			}
			continue
		}
		return ti, nil
	}
	for _, _formatter := range ExFormatters {
		tmp := dts
		var (
			ti  time.Time
			err error
		)
		if _formatter.HasZone {
			ti, err = time.Parse(_formatter.Format, tmp)
		} else {
			ti, err = time.ParseInLocation(_formatter.Format, tmp, loc)
		}
		if err != nil {
			if _formatter.Can(dts) {
				firstName, firstErr = _formatter.Name, err.Error()
			}
			continue
		}
		return ti, nil
	}

	_debugf("last use %v err %v\n", firstName, firstErr)
	jsonS := "{\"time\":\"" + dts + "\"}"
	type Tmp struct {
		Time time.Time `json:"time"`
	}
	ptr := Tmp{Time: time.Unix(long, 0)}
	err := Unmarshal([]byte(jsonS), ptr)
	return ptr.Time, err
}

var _enableDebugTime = false

func _debugf(format string, a ...any) {
	if _enableDebugTime {
		fmt.Printf(format, a...)
	}
}

func EnableDebugTime() {
	_enableDebugTime = true
}

func MaxTime(times ...time.Time) time.Time {
	res := time.Unix(0, 0)
	for _, dt := range times {
		if res.Before(dt) {
			res = dt
		}
	}
	return res
}

func MinTime(times ...time.Time) time.Time {
	res := time.Unix(UtsSecMax, 0)
	for _, dt := range times {
		if res.After(dt) {
			res = dt
		}
	}
	return res
}

func IntDay(value time.Time) int {
	return AsInt(value.Format(YYYYMMDD))
}

func MaxMinTime(values ...time.Time) (time.Time, time.Time) {
	vMin := time.Unix(UtsSecMax, 0)
	vMax := time.Unix(0, 0)
	for _, dt := range values {
		if dt.Before(vMin) {
			vMin = dt
		}
		if dt.After(vMax) {
			vMax = dt
		}
	}
	return vMax, vMin
}

func MoveInDay(start time.Time, h, m, s int) time.Time {
	return start.Add(time.Second * time.Duration(h*3600+m*60+s))
}

func MoveInYear(start time.Time, years, months, days int) time.Time {
	return start.AddDate(years, months, days)
}

func UnixMilli(t time.Time) int64 {
	return t.Unix()*1e3 + int64(t.Nanosecond())/1e6
}

func FromUnixMs(msec int64) time.Time {
	return time.Unix(msec/1e3, (msec%1e3)*1e6)
}

func ParseInterval(conf string) (int, string) {
	reValue := "\\d+(.\\d+)?"
	reUnit := "[dhms]"
	pattern, _ := regexp.Compile(fmt.Sprintf("(%s)?(%s)", reValue, reUnit))
	matches := pattern.FindAllStringSubmatch(conf, -1)
	if len(matches) > 0 {
		match := matches[0]
		if len(match) > 1 {
			value := AsDouble(match[1])
			unit := match[2]
			switch unit {
			case "d":
				return int(value * 86400), "d"
			case "h":
				return int(value * 3600), "h"
			case "m":
				return int(value * 60), "m"
			case "s":
				return int(value), "s"
			}
		}
	}
	return 3600, "h" // 默认1小时
}

const (
	UnitSec     = 1
	UnitMinute  = 60
	UnitHour    = UnitMinute * 60
	UnitDay     = UnitHour * 24
	UnitMonth   = UnitDay * 31
	UnitYear    = UnitDay * 365
	MUnitMs     = 1
	MUnitSec    = 1e3
	MUnitMinute = UnitMinute * 1e3
	MUnitHour   = UnitHour * 1e3
	MUnitDay    = UnitDay * 1e3
)

func TimeUnitStr(ms int64) string {
	if ms < MUnitSec {
		// Less than 1 second, show in milliseconds
		return fmt.Sprintf("%dms", ms)
	} else if ms < MUnitMinute {
		// Less than 1 minute, show in seconds
		return fmt.Sprintf("%.1fs", float64(ms)/MUnitSec)
	} else if ms < MUnitHour {
		// Less than 1 hour, show in minutes
		return fmt.Sprintf("%.1fm", float64(ms)/MUnitMinute)
	} else if ms < MUnitDay {
		// Less than 1 day, show in hours - use 23.9 as max for hours
		if ms >= MUnitDay-(MUnitMinute/2) {
			return "23.9h"
		}
		return fmt.Sprintf("%.1fh", float64(ms)/MUnitHour)
	} else {
		// 1 day or more, show in days
		return fmt.Sprintf("%.1fd", float64(ms)/MUnitDay)
	}
}
