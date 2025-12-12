package gox

import (
	"testing"
	"time"
)

func TestFloat2DT(t *testing.T) {
	t.Logf("name loc=%v,cn=%v\n", ZoneLocal(), ZoneCn())
	t.Logf("offset loc=%v,cn=%v\n", localZoneOffsetMs(), cnZoneOffsetMs())
	InitDtFormats()
	tests := map[string]string{
		"44693.4579861": "2022-05-12 10:59:29.999",
		"0.00":          "1899-12-31 00:00:00.000",
		"1.00":          "1900-01-01 00:00:00.000",
		"2.00":          "1900-01-02 00:00:00.000",
		"31.00":         "1900-01-31 00:00:00.000",
		"32.00":         "1900-02-01 00:00:00.000",
		"59.00":         "1900-02-28 00:00:00.000",
		"60.00":         "1900-03-01 00:00:00.000",
		"61.00":         "1900-03-01 00:00:00.000",
		"366.00":        "1900-12-31 00:00:00.000",
		"367.00":        "1901-01-01 00:00:00.000",
		"3654.00":       "1910-01-01 00:00:00.000",
		"7306.00":       "1920-01-01 00:00:00.000",
		"10959.00":      "1930-01-01 00:00:00.000",
		"14611.00":      "1940-01-01 00:00:00.000",
		"18264.00":      "1950-01-01 00:00:00.000",
		"21916.00":      "1960-01-01 00:00:00.000",
		"25568.00":      "1969-12-31 00:00:00.000",
		"25569.00":      "1970-01-01 00:00:00.000",
		"1.0000001":     "1900-01-01 00:00:00.009",
	}
	for src, want := range tests {
		dt, err := tryFloat2DT(src)
		if err != nil {
			t.Errorf("src=%v want=%v but err=%v", src, want, err)
		} else {
			got := FormatCnZone(dt, ISODateTimeMs)
			if len(got) == 19 {
				got = got + ".000"
			}
			if got != want {
				t.Errorf("src=%v want=%v but got=%v", src, want, got)
			} else {
				t.Logf("src=%v want=%v and got=%v", src, want, got)
			}
		}
	}
}

func TestParseDt(t *testing.T) {
	curYear := AsStr(time.Now().Year())
	tests := map[string]string{
		"44693.4579861":           "2022-05-12 10:59:29.999",
		"2022-05-12 10:59:29.999": "2022-05-12 10:59:29.999",
		"2022-05-12":              "2022-05-12 00:00:00.000",
		"2022-5-12":               "2022-05-12 00:00:00.000",
		"22-5-12":                 "2022-05-12 00:00:00.000",
		//"5-12-22":   "2022-05-12 00:00:00.000", // can not tell apart form 12-5-22
		"2022年5月12日":  "2022-05-12 00:00:00.000",
		"2022年05月12日": "2022-05-12 00:00:00.000",
		"22年5月12日":    "2022-05-12 00:00:00.000",
		"22年05月12日":   "2022-05-12 00:00:00.000",
		"5月12日22年":    "2022-05-12 00:00:00.000",
		"5月12日":       curYear + "-05-12 00:00:00.000",
	}
	InitDtFormats()
	for src, want := range tests {
		t.Run(src, func(tt *testing.T) {
			got := FormatCnZone(ParseDt(src), ISODateTimeMs)
			if len(got) == 19 {
				got = got + ".000"
			}
			if got != want {
				tt.Errorf("want=%v but got=%v", want, got)
			}
		})
	}
}
