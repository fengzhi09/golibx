package jsonx

import (
	"github.com/fengzhi09/golibx/gox"
	"testing"
	"time"
)

func Test_ToTime(t *testing.T) {
	data := "{" +
		"\"int_date\":20220406," +
		"\"int_date_time\":20060102150405," +
		"\"ts_sec_int64\":1649233212," +
		"\"ts_ms_int64\":\"1649233212047\"," +
		"\"YYYYMMDD_str\":\"20220406\"," +
		"\"YYYYMMDDhhmmss_str\":\"20060102150405\"," +
		"\"go_time_str\":\"2022-04-06T08:20:08.682Z\"," +
		"\"go_time_str_short\":\"2022-04-06T08:20:08Z\"," +
		"\"go_time_ns_str\":\"2022-04-06T08:20:08.123456789Z\"," +
		"\"iso_time_ms_str\":\"2022-04-06 08:20:09.047\"," +
		"\"iso_time_str\":\"2022-04-06 08:20:21\"," +
		"\"ISODateTimeNsTZF\":\"2022-04-07T20:36:18.4869717+08:00\"," +
		"\"RFC1123Z\":\"Thu, 07 Apr 2022 20:39:23 +0800\"," +
		"\"RFC850\":\"Thursday, 07-Apr-22 20:39:23 CST\"," +
		"\"RFC1123\":\"Thu, 07 Apr 2022 20:39:23 CST\"," +
		"\"RubyDateZ\":\"Thu Apr 07 20:39:23 +0800 2022\"," +
		"\"UnixDateZ\":\"Thu Apr  7 20:39:23 CST 2022\"," +
		"\"RFC3339Z\":\"2022-04-07T20:39:23+08:00\"," +
		"\"ANSIC\":\"Thu Apr  7 20:39:23 2022\"," +
		"\"ISODateTimeMs\":\"2022-04-07 20:39:23.969\"," +
		"\"ISODateTimeMs2\":\"2022-04-07 20:39:23.96\"," +
		"\"RFC822Z\":\"07 Apr 22 20:39 +0800\"," +
		"\"ISODateTimeMs1\":\"2022-11-02 18:04:53.9\"," +
		"\"RFC822\":\"07 Apr 22 20:39 CST\"," +
		"\"ISODateTime\":\"2022-04-07 20:39:23\"," +
		"\"YYYYMMDDHHMMSS\":\"20220407203923\"," +
		"\"ISODate\":\"2022-04-07\"," +
		"\"YYYYMMDD\":\"20220407\"," +
		"\"empty\":\"\"" +
		"}"

	obj := ParseJObj(data)
	keys := obj.Keys()
	for _, key := range keys {
		ti := obj.GetTime(key)
		if ti.Year() < 2000 && key != "empty" {
			t.Errorf("%v got time=%v str=%v ", key, ti, obj.GetVal(key).String())
		}
	}

}

func Test_FormatTime(t *testing.T) {
	ti := time.Now()
	for _, formatter := range gox.Formatters {
		t.Logf("%v got %v", formatter.Name, ti.Format(formatter.Format))
	}
}
func TestName(t *testing.T) {
	ticK := time.Now()
	ticS := "2022-04-07T20:36:18.4869717+08:00"
	tickCopy, e := time.Parse(time.RFC3339Nano, ticS)
	t.Log(ticK)
	t.Log(ticS)
	t.Log(tickCopy)
	t.Log(e)
}

func Test_AsTime(t *testing.T) {
	data := "{" +
		"\"int_date\":20220406," +
		"\"int_date_time\":20060102150405," +
		"\"ts_sec_int64\":1649233212," +
		"\"ts_ms_int64\":\"1649233212047\"," +
		"\"YYYYMMDD_str\":\"20220406\"," +
		"\"YYYYMMDDhhmmss_str\":\"20060102150405\"," +
		"\"go_time_str\":\"2022-04-06T08:20:08.682Z\"," +
		"\"go_time_str_short\":\"2022-04-06T08:20:08Z\"," +
		"\"go_time_ns_str\":\"2022-04-06T08:20:08.123456789Z\"," +
		"\"iso_time_ms_str\":\"2022-04-06 08:20:09.047\"," +
		"\"iso_time_str\":\"2022-04-06 08:20:21\"," +
		"\"ISODateTimeNsTZF\":\"2022-04-07T20:36:18.4869717+08:00\"," +
		"\"RFC1123Z\":\"Thu, 07 Apr 2022 20:39:23 +0800\"," +
		"\"RFC850\":\"Thursday, 07-Apr-22 20:39:23 CST\"," +
		"\"RFC1123\":\"Thu, 07 Apr 2022 20:39:23 CST\"," +
		"\"RubyDateZ\":\"Thu Apr 07 20:39:23 +0800 2022\"," +
		"\"UnixDateZ\":\"Thu Apr  7 20:39:23 CST 2022\"," +
		"\"RFC3339Z\":\"2022-04-07T20:39:23+08:00\"," +
		"\"ANSIC\":\"Thu Apr  7 20:39:23 2022\"," +
		"\"ISODateTimeMs\":\"2022-04-07 20:39:23.969\"," +
		"\"ISODateTimeMs2\":\"2022-04-07 20:39:23.96\"," +
		"\"RFC822Z\":\"07 Apr 22 20:39 +0800\"," +
		"\"ISODateTimeMs1\":\"2022-11-02 18:04:53.9\"," +
		"\"RFC822\":\"07 Apr 22 20:39 CST\"," +
		"\"ISODateTime\":\"2022-04-07 20:39:23\"," +
		"\"YYYYMMDDHHMMSS\":\"20220407203923\"," +
		"\"ISODate\":\"2022-04-07\"," +
		"\"YYYYMMDD\":\"20220407\"," +
		"\"empty\":\"\"" +
		"}"

	obj := ParseJObj(data)
	keys := obj.Keys()
	for _, key := range keys {
		txt := obj.GetStr(key)
		ti := gox.AsTime(txt)
		if ti.Year() < 2000 && key != "empty" {
			t.Errorf("%v got time=%v str=%v ", key, ti.Format(gox.ISODateTimeMs), txt)
		}
	}
}

func Test_Float(t *testing.T) {
	v1, v2, v3, v4, v5, v6, v7 := 0.0, 1.01, 23.67, 1e9, 123456789, 1e12, 1234567890123
	t.Logf("v1 %v | %f | %d | %e | %g | %s", v1, v1, int(v1), v1, v1, gox.AsStr(v1))
	t.Logf("v2 %v | %f | %d | %e | %g | %s", v2, v2, int(v2), v2, v2, gox.AsStr(v2))
	t.Logf("v3 %v | %f | %d | %e | %g | %s", v3, v3, int(v3), v3, v3, gox.AsStr(v3))
	t.Logf("v4 %v | %f | %d | %e | %g | %s", v4, v4, int(v4), v4, v4, gox.AsStr(v4))
	t.Logf("v5 %v | %f | %d | %e | %g | %s", v5, float64(v5), int(v5), float64(v5), float64(v5), gox.AsStr(v5))
	t.Logf("v6 %v | %f | %d | %e | %g | %s", v6, float64(v6), int(v6), float64(v6), float64(v6), gox.AsStr(v6))
	t.Logf("v5 %v | %f | %d | %e | %g | %s", v7, float64(v7), int(v7), float64(v7), float64(v7), gox.AsStr(v7))

}
