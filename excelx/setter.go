package excelx

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"

	"github.com/dlclark/regexp2"
)

type _getter func(cell string, conf FieldConf) (any, error)

var (
	FT_IntDate       = "int_date"
	FT_Date          = "date"
	FT_Time          = "time"
	FT_DateTime      = "date_time"
	FT_String        = "string"
	FT_Float         = "float"
	FT_Money         = "money"
	FT_Int           = "int"
	FT_Bool          = "bool"
	FT_DateTimeArray = "[]date_time"
	FT_StringArray   = "[]string"
	FT_FloatArray    = "[]float"
	FT_IntArray      = "[]int"
	FT_BoolArray     = "[]bool"
	FT_ZipFields     = "zip_fields" // 按json键值对进行折叠 {k1:[v11,v12],k2:v2}
	FT_Enum          = "enum"
	FT_EnumArray     = "[]enum"
	FT_Regexp        = "regexp"
	FT_CondEnum      = "cond_enum"
)

var getters = map[FieldType]_getter{
	FT_DateTime: func(cell string, conf FieldConf) (any, error) {
		dt := ParseDtCell(cell)
		if dt.Unix() == 0 {
			return FormatDt(dt, conf.OutputFormat), fmt.Errorf("非法日期:%v", cell)
		}
		return FormatDt(dt, conf.OutputFormat), nil
	},
	FT_String: func(cell string, conf FieldConf) (any, error) {
		return cell, nil
	},
	FT_Float: func(cell string, conf FieldConf) (any, error) {
		return gox.AsDouble(cell), nil
	},
	FT_Money: func(cell string, conf FieldConf) (any, error) {
		v := gox.AsMoney(cell, 2)
		if v < 0 {
			return -0, fmt.Errorf("非法值：%v", cell)
		}
		return v, nil
	},
	FT_Int: func(cell string, conf FieldConf) (any, error) {
		return int64(gox.AsDouble(cell)), nil
	},
	FT_Bool: func(cell string, conf FieldConf) (any, error) {
		return ParseBool(cell), nil
	},
	FT_DateTimeArray: func(cell string, conf FieldConf) (any, error) {
		return ParseJArr(cell, func(elem string) any {
			return FormatDt(gox.ParseDt(elem), conf.OutputFormat)
		}), nil
	},
	FT_StringArray: func(cell string, conf FieldConf) (any, error) {
		return ParseJArr(cell, asStr), nil
	},
	FT_FloatArray: func(cell string, conf FieldConf) (any, error) {
		return ParseJArr(cell, asDouble), nil
	},
	FT_IntArray: func(cell string, conf FieldConf) (any, error) {
		return ParseJArr(cell, asInt), nil
	},
	FT_BoolArray: func(cell string, conf FieldConf) (any, error) {
		return ParseJArr(cell, asBool), nil
	},
	FT_IntDate: func(cell string, conf FieldConf) (any, error) {
		dt := ParseDtCell(cell)
		if dt.Unix() == 0 {
			return 0, fmt.Errorf("非法日期:%v", cell)
		}
		date, _ := strconv.Atoi(gox.FormatCnZone(gox.AsTime(dt), gox.YYYYMMDD))
		return date, nil
	},
	FT_Date: func(cell string, conf FieldConf) (any, error) {
		dt := ParseDtCell(cell)
		if dt.Unix() == 0 {
			return 0, fmt.Errorf("非法日期:%v", cell)
		}
		return gox.AsTime(dt).Format(gox.ISODate), nil
	},
	FT_Time: func(cell string, conf FieldConf) (any, error) {
		dt := ParseDtCell(cell)
		if dt.Unix() == 0 {
			return 0, fmt.Errorf("非法日期:%v", cell)
		}
		return gox.AsTime(dt).Format(gox.ISOTime), nil
	},
	FT_Enum: func(cell string, conf FieldConf) (any, error) {
		enum, err := ParseEnum(cell, conf.EnumMapper, conf.EnumStrict)
		return enum, err
	},
	FT_EnumArray: func(cell string, conf FieldConf) (any, error) {
		enums, err := ParseEnums(ParseJArr(cell, asStr).ToStrArr(), conf.EnumMapper, conf.EnumStrict)
		if err != nil {
			return nil, err
		}
		arr := jsonx.NewJArr(enums...)
		return arr, nil
	},
	FT_Regexp: func(cell string, conf FieldConf) (any, error) {
		str := cell
		isMatch := func() bool {
			pattern := regexp.MustCompile(conf.Regexp)
			matches := pattern.FindAllString(str, -1)
			for _, match := range matches {
				if match == str {
					return true
				}
			}
			return false
		}
		if conf.Regexp != "" && isMatch() {
			return str, nil
		} else {
			return nil, fmt.Errorf("格式不匹配")
		}
	},
	FT_CondEnum: func(cell string, conf FieldConf) (any, error) {
		for enum, reg := range conf.CondEnum {
			regMachine, err := regexp2.Compile(reg, 0)
			hit, err := regMachine.MatchString(cell)
			if hit {
				return enum, err
			}
		}
		return cell, nil
	},
}

func ParseEnum(raw string, mappers map[string]any, strict bool) (any, error) {
	if value, hit := mappers[raw]; hit {
		return value, nil
	} else if !strict {
		return raw, nil
	} else {
		return nil, fmt.Errorf("非法值:'%v'", raw)
	}
}

func ParseEnums(raws []string, mappers map[string]any, strict bool) ([]any, error) {
	enums := make([]any, 0)
	for _, raw := range raws {
		if value, hit := mappers[raw]; hit {
			enums = append(enums, value)
		} else if !strict {
			enums = append(enums, raw)
		} else {
			return nil, fmt.Errorf("非法值:'%v'", raw)
		}
	}
	return enums, nil
}

type converter func(elem string) any

var (
	asDouble converter = func(elem string) any { return gox.AsDouble(elem) }
	asBool   converter = func(elem string) any { return ParseBool(elem) }
	asInt    converter = func(elem string) any { return gox.AsLong(elem) }
	asStr    converter = nil
)

func ParseDtCell(cell string) time.Time {
	return gox.ParseDt(cell)
}

var formatsMapper = map[string]string{
	"":                        gox.ISODateTimeMs,
	"yyyymmdd":                gox.YYYYMMDD,
	"yyyymmddhhmmss":          gox.YYYYMMDDHHMMSS,
	"yyyy-mm-dd":              gox.ISODate,
	"yyyy-mm-dd hh:mm:ss":     gox.ISODateTime,
	"hh:mm:ss":                gox.ISOTime,
	"yyyy-mm-dd hh:mm:ss.fff": gox.ISODateTimeMs,
}

func FormatDt(dt time.Time, format string) string {
	format = strings.TrimSpace(strings.ToLower(format))
	if format == "timestamp" {
		return gox.AsStr(dt.Unix())
	}
	if goFormat, hit := formatsMapper[format]; hit {
		format = goFormat
	}
	return gox.FormatCnZone(dt, format)
}
