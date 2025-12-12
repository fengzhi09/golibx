package gox

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

/*
excel中1970-01-01 00:00:00.000 对应1
excel中1970-03-01 00:00:00.000 对应61
多了实际不存在的1970-02-29
*/
const excel19000301 = float64(61.0000000)

// excel中1970-01-01 00:00:00.000 对应25569
const (
	excel19700101    = float64(25569.0000000)
	ExcelDFFormatter = "yyyy-mm-dd hh:mm:ss.000"
)

func InitDtFormats() {
	excelFormatters := []Formatter{
		CNFullFormatter,
		CNDateFormatter,
		CNDateShortFormatter,
		CNTimeFormatter,
		ExcelShort0Formatter,
		ExcelShortDateFRFormatter,
		ExcelShortDate0Formatter,
		ExcelShortestDateFormatter,
	}
	ExFormatters = append(excelFormatters, ExFormatters...)
}

func localZoneOffsetMs() int64 {
	return zoneOffsetMs(ZoneLocal(), ZoneUTC())
}

func cnZoneOffsetMs() int64 {
	return zoneOffsetMs(ZoneCn(), ZoneUTC())
}

func zoneOffsetMs(loc1, loc2 *time.Location) int64 {
	time1 := time.Date(1970, 1, 1, 0, 0, 0, 0, loc1).Unix()
	time2 := time.Date(1970, 1, 1, 0, 0, 0, 0, loc2).Unix()
	offset := (time2 - time1) * 1000
	return offset
}

func ParseDt(raw string) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}
	}
	dt, err := UniformDt(raw)
	if err == nil {
		return dt
	}
	dt, err = tryFloat2DT(raw)
	if err == nil {
		return dt
	}
	raw = tryFixCN(raw)
	raw = tryFill0IfShort(raw)
	dt, err = UniformDt(raw)
	if err == nil {
		return dt
	}
	return time.Time{}
}

var (
	CNFullFormat    = "2006年01月02日15时04分05秒"
	CNFullFormatter = Formatter{Name: "cn-full", Format: CNFullFormat, Can: func(dts string) bool {
		return LenUnt8(dts) >= LenUnt8(CNFullFormat) && strings.Contains(dts, "年") && strings.Contains(dts, "分")
	}, HasZone: false}

	CNDateFormat    = "2006年01月02日"
	CNDateFormatter = Formatter{Name: "cn-date", Format: CNDateFormat, Can: func(dts string) bool {
		return LenUnt8(dts) >= LenUnt8(CNDateFormat) && strings.Contains(dts, "年") && !strings.Contains(dts, "分")
	}, HasZone: false}

	CNDateShortFormat    = "01月02日"
	CNDateShortFormatter = Formatter{Name: "cn-date", Format: CNDateFormat, Can: func(dts string) bool {
		return LenUnt8(dts) >= LenUnt8(CNDateShortFormat) && strings.Contains(dts, "月") && !strings.Contains(dts, "日")
	}, HasZone: false}

	CNTimeFormat    = "15时04分05秒"
	CNTimeFormatter = Formatter{Name: "cn-time", Format: CNTimeFormat, Can: func(dts string) bool {
		return LenUnt8(dts) >= LenUnt8(CNTimeFormat) && !strings.Contains(dts, "年") && strings.Contains(dts, "分")
	}, HasZone: false}

	ExcelShort0Format    = "06-01-02 15:04:05"
	ExcelShort0Formatter = Formatter{Name: "excel-short", Format: ExcelShort0Format, Can: func(dts string) bool {
		return len(dts) >= len(ExcelShort0Format)-4 && strings.Count(dts, "/") == 2 && strings.Count(dts, ":") == 2
	}, HasZone: false}

	ExcelShortDate0Format    = "06-01-02"
	ExcelShortDate0Formatter = Formatter{Name: "excel-date-short-0", Format: ExcelShortDate0Format, Can: func(dts string) bool {
		return len(dts) == len(ExcelShortDate0Format)-2 && strings.Count(dts, "-") == 2
	}, HasZone: false}

	ExcelShortDateFRFormat    = "01-02-06"
	ExcelShortDateFRFormatter = Formatter{Name: "excel-date-short-r", Format: ExcelShortDateFRFormat, Can: func(dts string) bool {
		return len(dts) == len(ExcelShortDateFRFormat)-2 && strings.Count(dts, "-") == 2
	}, HasZone: false}

	ExcelShortestDateFormat    = "01-02"
	ExcelShortestDateFormatter = Formatter{Name: "excel-date-shortest", Format: ExcelShortestDateFormat, Can: func(dts string) bool {
		return len(dts) == len(ExcelShortDateFRFormat)-2 && strings.Count(dts, "-") == 1
	}, HasZone: false}
)

func tryFloat2DT(raw string) (time.Time, error) {
	if len(strings.Split(raw, ".")) == 2 {
		value := AsDouble(raw)
		if value >= 0 {
			value = IfElse(value < excel19000301, value+1, value).(float64)
			msTotal := (value-excel19700101)*86400000 - float64(cnZoneOffsetMs())
			return FromUnixMs(int64(msTotal)), nil
		}
	}
	return time.Time{}, fmt.Errorf("not float dt of excel")
}

func tryFixCN(raw string) string {
	if strings.Contains(raw, "年") || strings.Contains(raw, "月") || strings.Contains(raw, "日") {
		getVStr := func(curStr string, stopWords ...string) string {
			for _, stopWord := range stopWords {
				if strings.Contains(raw, stopWord) {
					return strings.TrimSpace(strings.ReplaceAll(regexp.MustCompile("\\d+"+stopWord).FindString(raw), stopWord, ""))
				}
			}
			return curStr
		}
		cY, cMon, cD := time.Now().Date()
		yStr, monStr, dStr := getVStr(AsStr(cY), "年"), getVStr(AsStr(cMon), "月"), getVStr(AsStr(cD), "日")
		hStr, minStr, secStr := getVStr("00", "点", "时"), getVStr("00", "分"), getVStr("00", "秒")
		return fmt.Sprintf("%v-%v-%v %v:%v:%v", yStr, monStr, dStr, hStr, minStr, secStr)
	}
	return raw
}

func tryFill0IfShort(raw string) string {
	raw = strings.ReplaceAll(raw, "/", "-")
	seps := []string{"-", ":"}
	dates, times := parseDatesTimes(strings.Split(raw, " "), seps)
	if len(dates) == 3 && len(dates[0]) == 2 && len(dates[1]) == 2 && len(dates[2]) == 2 {
		dates[0] = "20" + dates[0]
	}
	return strings.Join(dates, "-") + " " + strings.Join(times, ":")
}

func parseDatesTimes(parts []string, seps []string) ([]string, []string) {
	dates, times := make([]string, 0), make([]string, 0)
	for _, part := range parts {
		for i, sep := range seps {
			isDateSep := i == 0
			if strings.Contains(part, sep) {
				elems := strings.Split(part, sep)
				for _, elem := range elems {
					if len(elem) < 2 {
						elem = "0" + elem
					}
					if isDateSep {
						dates = append(dates, elem)
					} else {
						times = append(times, elem)
					}
				}
			}
		}
	}
	return dates, times
}
