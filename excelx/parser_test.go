// nolint
package excelx

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"
)

type (
	emapper = map[string]any
	mapper  = map[string]string
	kv      = [2]string
	row     = []kv
	_case   struct {
		row    row
		conf   FieldConf
		isWant func(val jsonx.JValue) string
		reason string
	}
)

const (
	patternCn = "(年度)|(月度)|(季度)"
	patternEn = "(\\d{6})|(\\d{4}Q[1234])|(\\d{4})"
)

//nolint:funlen
func TestParseField(t *testing.T) {
	gox.InitDtFormats()
	_cases := map[string]*_case{
		"one-int-date":          mockCase(row{kv{"日期", "2022-05-12"}}, "日期", FT_IntDate, "20220512"),
		"one-int-date-long":     mockCase(row{kv{"日期", "2022/05/12"}}, "日期", FT_IntDate, "20220512"),
		"one-int-date-cn":       mockCase(row{kv{"日期", "2022年05月12日"}}, "日期", FT_IntDate, "20220512"),
		"one-int-date-short":    mockCase(row{kv{"日期", "22/05/12"}}, "日期", FT_IntDate, "20220512"),
		"one-int-date-float-2":  mockCase(row{kv{"日期", "44693.00"}}, "日期", FT_IntDate, "20220512"),
		"one-int-date-float-10": mockCase(row{kv{"日期", "44693.4579861111"}}, "日期", FT_IntDate, "20220512"),
		"one-int-date-YYYYMMDD": mockCase(row{kv{"日期", "20220512"}}, "日期", FT_IntDate, "20220512"),

		"one-str":                      mockCase(row{kv{"a", "a"}}, "a", FT_String, "a"),
		"one-str-cn":                   mockCase(row{kv{"列名", "单元格"}}, "列名", FT_String, "单元格"),
		"one-date-time-float":          mockCase(row{kv{"日期", "44693.4579861111"}}, "日期", FT_DateTime, "2022-05-12 10:59:29.999"),
		"one-date-time-str":            mockCase(row{kv{"日期", "2022-05-12 10:59:30"}}, "日期", FT_DateTime, "2022-05-12 10:59:30"),
		"one-date-time-str-short":      mockCase(row{kv{"日期", "22/5/12 10:59:30"}}, "日期", FT_DateTime, "2022-05-12 10:59:30"),
		"one-date-time-str-cn":         mockCase(row{kv{"日期", "2022年05月12日 10时59分30秒"}}, "日期", FT_DateTime, "2022-05-12 10:59:30"),
		"one-date-time-str-long":       mockCase(row{kv{"日期", "2022/05/12 10:59:30"}}, "日期", FT_DateTime, "2022-05-12 10:59:30"),
		"one-date-time-ts":             mockCase(row{kv{"日期", "1652324370"}}, "日期", FT_DateTime, "2022-05-12 10:59:30"),
		"one-date-time-YYYYMMDDhhmmss": mockCase(row{kv{"日期", "20220512105930"}}, "日期", FT_DateTime, "2022-05-12 10:59:30"),

		"one-float":   mockCase(row{kv{"小数", "1.2"}}, "小数", FT_Float, "1.2"),
		"one-float-1": mockCase(row{kv{"小数", "1.2"}}, "小数", FT_Float, "1.2"),
		"one-float-2": mockCase(row{kv{"小数", "0.1234567890123456"}}, "小数", FT_Float, "0.1234567890123456"),
		"one-float-3": mockCase(row{kv{"小数", "1.23e01"}}, "小数", FT_Float, "12.3"),
		"one-float-4": mockCase(row{kv{"小数", "-1"}}, "小数", FT_Float, "-1"),

		"one-int-1":        mockCase(row{kv{"整数", "1"}}, "整数", FT_Int, "1"),
		"one-int-1.2":      mockCase(row{kv{"整数", "1.2"}}, "整数", FT_Int, "1"),
		"one-int-1.6":      mockCase(row{kv{"整数", "1.6"}}, "整数", FT_Int, "1"),
		"one-int-neg-1":    mockCase(row{kv{"整数", "-1"}}, "整数", FT_Int, "-1"),
		"one-int-max-long": mockCase(row{kv{"整数", "1000000000000000000"}}, "整数", FT_Int, "1000000000000000000"),

		"one-bool-yes":   mockCase(row{kv{"布尔", "yes"}}, "布尔", FT_Bool, "true"),
		"one-bool-no":    mockCase(row{kv{"布尔", "no"}}, "布尔", FT_Bool, "false"),
		"one-bool-true":  mockCase(row{kv{"布尔", "true"}}, "布尔", FT_Bool, "true"),
		"one-bool-false": mockCase(row{kv{"布尔", "false"}}, "布尔", FT_Bool, "false"),
		"one-bool-是":     mockCase(row{kv{"布尔", "是"}}, "布尔", FT_Bool, "true"),
		"one-bool-否":     mockCase(row{kv{"布尔", "否"}}, "布尔", FT_Bool, "false"),
		"one-bool-1":     mockCase(row{kv{"布尔", "1"}}, "布尔", FT_Bool, "true"),
		"one-bool-0":     mockCase(row{kv{"布尔", "0"}}, "布尔", FT_Bool, "false"),

		"arr-str-[]": mockCase(row{kv{"a", "[\"a\",\"b\"]"}}, "a", FT_StringArray, "[\"a\",\"b\"]"),
		"arr-str-;":  mockCase(row{kv{"a", "a;b"}}, "a", FT_StringArray, "[\"a\",\"b\"]"),
		"arr-str-,":  mockCase(row{kv{"a", "a,b"}}, "a", FT_StringArray, "[\"a\",\"b\"]"),

		"arr-data-time": mockCase(row{kv{"arr-data-time", "22/5/12 10:59:30;2022-05-12 10:59:30"}}, "arr-data-time", FT_DateTimeArray, "[\"2022-05-12 10:59:30\",\"2022-05-12 10:59:30\"]"),
		"arr-float":     mockCase(row{kv{"arr-float", "1;1.2;1.6;-1.2;0.1234567890123456;1.23e01"}}, "arr-float", FT_FloatArray, "[1,1.2,1.6,-1.2,0.1234567890123456,12.3]"),
		"arr-int":       mockCase(row{kv{"arr-int", "1;9223372036854775807;-1"}}, "arr-int", FT_IntArray, "[1,9223372036854775807,-1]"),
		"arr-bool":      mockCase(row{kv{"arr-bool", "0;是;yes;false"}}, "arr-bool", FT_BoolArray, "[false,true,true,false]"),

		"arr-zip_fields": mockZip(row{kv{"a", "1"}, kv{"b", "2"}}, "zip", FT_ZipFields, mapper{"a": "A", "b": "B", "c": "C"}),
		"one-enum":       mockEnum(row{kv{"状态", "开始"}}, "状态", FT_Enum, true, emapper{"开始": 1}),
		"arr-enum":       mockEnum(row{kv{"状态", "开始,等待,停止"}}, "状态", FT_EnumArray, false, emapper{"开始": 1, "停止": "2"}),

		"one-alias-1-1":       mockAlias(row{kv{"1", "1"}, kv{"2", "2"}}, "列名0", FT_String, "1", "1"),
		"one-alias-a,b-a,b-b": mockAlias(row{kv{"a", "a"}, kv{"b", "b"}}, "列名1", FT_String, "b", "a", "b"),
		"one-alias-a,b-b,a-b": mockAlias(row{kv{"a", "a"}, kv{"b", "b"}}, "列名2", FT_String, "b", "b", "a"),
		"one-alias-b,a-b,a-a": mockAlias(row{kv{"b", "b"}, kv{"a", "a"}}, "列名2", FT_String, "a", "b", "a"),

		"regexp-cn-周期1-年度":      mockRegexp(row{kv{"周期1", "年度"}}, "周期1", FT_Regexp, patternCn, ""),
		"regexp-cn-周期2-月度":      mockRegexp(row{kv{"周期2", "年度"}}, "周期2", FT_Regexp, patternCn, ""),
		"regexp-cn-周期3-年度1":     mockRegexp(row{kv{"周期3", "年度1"}}, "周期3", FT_Regexp, patternCn, "不匹配"),
		"regexp-en-周期4-2022":    mockRegexp(row{kv{"周期4", "2022"}}, "周期4", FT_Regexp, patternEn, ""),
		"regexp-en-周期5-202201":  mockRegexp(row{kv{"周期5", "202201"}}, "周期5", FT_Regexp, patternEn, ""),
		"regexp-en-周期6-2022Q1":  mockRegexp(row{kv{"周期6", "2022Q1"}}, "周期6", FT_Regexp, patternEn, ""),
		"regexp-en-周期7-2022Q1f": mockRegexp(row{kv{"周期7", "2022Q1f"}}, "周期7", FT_Regexp, patternEn, "不匹配"),
	}
	for cname, case_ := range _cases {
		t.Run(cname, func(tt *testing.T) {
			field := case_.conf.StdName
			data := &jsonx.JObj{}
			var err error
			for _, kv_ := range case_.row {
				err = ParseField(kv_[1], kv_[0], case_.conf, data)
				if err != nil {
					break
				}
			}
			if case_.reason != "" {
				if err == nil || !strings.Contains(err.Error(), case_.reason) {
					tt.Errorf("reason want=%v but got=%v", case_.reason, err)
				}
			} else {
				got := data.GetVal(field)
				if diff := case_.isWant(got); diff != "" {
					tt.Errorf("value got=%v diff=%v", got.String(), diff)
				}
			}
		})
	}
}

func mockCase(row row, field string, fType FieldType, want string) *_case {
	return &_case{
		row: row,
		conf: FieldConf{
			StdName:    field,
			Alias:      []string{},
			OutputType: fType,
		},
		isWant: func(val jsonx.JValue) string {
			got := strings.ReplaceAll(val.String(), " ", "")
			target := strings.ReplaceAll(want, " ", "")
			return gox.IfElse(got == target, "", want).(string)
		},
		reason: "",
	}
}

func mockAlias(row row, field string, fType FieldType, want string, alias ...string) *_case {
	return &_case{
		row: row,
		conf: FieldConf{
			StdName:    field,
			Alias:      alias,
			OutputType: fType,
		},
		isWant: func(val jsonx.JValue) string {
			got := strings.ReplaceAll(val.String(), " ", "")
			target := strings.ReplaceAll(want, " ", "")
			return gox.IfElse(got == target, "", want).(string)
		},
		reason: "",
	}
}

func mockZip(row row, field string, fType FieldType, zips mapper) *_case {
	alias, hits := make([]string, 0), map[string]bool{}
	for k := range zips {
		alias = append(alias, k)
	}
	for _, kv_ := range row {
		hits[kv_[0]] = true
	}
	return &_case{
		row: row,
		conf: FieldConf{
			StdName:    field,
			Alias:      alias,
			OutputType: fType,
			ZipFields:  zips,
		},
		isWant: func(val jsonx.JValue) string {
			got := strings.ReplaceAll(val.String(), " ", "")
			for k, v := range zips {
				if !strings.Contains(got, v) && hits[k] {
					return fmt.Sprintf("field=%v not zipped as %v.%v", k, field, v)
				}
			}
			return ""
		},
		reason: "",
	}
}

func mockEnum(row row, field string, fType FieldType, strict bool, mapper emapper) *_case {
	invalids := make([]string, 0)
	for k := range mapper {
		invalids = append(invalids, k)
	}
	return &_case{
		row: row,
		conf: FieldConf{
			StdName:    field,
			OutputType: fType,
			EnumStrict: strict,
			EnumMapper: mapper,
		},
		isWant: func(val jsonx.JValue) string {
			got := strings.ReplaceAll(val.String(), " ", "")
			if strict {
				for _, k := range invalids {
					if strings.Contains(got, k) {
						return fmt.Sprintf("enum=%v not mapper as %v", k, mapper[k])
					}
				}
			}
			return ""
		},
		reason: "",
	}
}

func mockRegexp(row row, field string, fType FieldType, pattern, reason string) *_case {
	return &_case{
		row: row,
		conf: FieldConf{
			StdName:    field,
			OutputType: fType,
			Regexp:     pattern,
		},
		isWant: func(val jsonx.JValue) string {
			got := strings.ReplaceAll(val.String(), " ", "")
			match, err := regexp.MatchString(pattern, got)
			if err != nil {
				return err.Error()
			}
			return gox.IfElse(match, "", fmt.Sprintf("不匹配 %v", pattern)).(string)
		},
		reason: reason,
	}
}
