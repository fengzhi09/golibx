package gox

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func AsTime(val any) time.Time {
	if val == nil {
		return FromUnixMs(0)
	}
	str := ""
	switch val.(type) {
	case int8, int16, int32, int, int64, uint8, uint16, uint32, uint, uint64:
		long := val.(int64)
		if IntDateMin <= long && long < IntDateMax {
			year, month, date := int(long/10000), int(long/100%100), int(long%100)
			return time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.Local)
		}
		if UtsSecMin <= long && long < UtsSecMax {
			return time.Unix(long, 0)
		}
		if UtsMsMin <= long && long < UtsMsMax {
			return FromUnixMs(long)
		}
		if LongDateMin <= long && long < LongDateMax {
			long1, long2 := long/1e6, long%1e6
			year, month, date := int(long1/10000), int(long1/100%100), int(long1%100)
			hour, minute, second := int(long2/10000), int(long2/100%100), int(long2%100)
			return time.Date(year, time.Month(month), date, hour, minute, second, 0, time.Local)
		}
		str = fmt.Sprintf("%d", long)
		break
	case string:
		str = val.(string)
		break
	default:
		str = fmt.Sprintf("%v", val)
		break
	}
	if len(str) > 0 {
		ts, err := UniformDt(str)
		if err == nil {
			return ts
		}
	}
	return time.Unix(0, 0)
}

func AsStr(val any) string {
	switch val.(type) {
	case string:
		return val.(string)
	case int8, int16, int32, int, int64, uint8, uint16, uint32, uint, uint64:
		return fmt.Sprintf("%d", val)
	case float32:
		return Ftoa(val.(float32))
	case float64:
		return Dtoa(val.(float64))
	case bool:
		return IfElse(val.(bool), "true", "false").(string)
	default:
		vt := reflect.TypeOf(val)
		switch vt.Kind() {
		case reflect.Bool:
			return IfElse(val.(bool), "true", "false").(string)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%d", val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fmt.Sprintf("%d", val)
		case reflect.Float32, reflect.Float64:
			return Ftoa(val.(float32))
		case reflect.String:
			return Dtoa(val.(float64))
		}
		return fmt.Sprintf("%v", val)
	}
}

func AsBool(val any) bool {
	switch val.(type) {
	case int8, int16, int32, int, int64, uint8, uint16, uint32, uint, uint64:
		return AsLong(val) > 0
	case bool:
		return val.(bool)
	case string:
		return strings.ToLower(AsStr(val)) == "true"
	default:
		return val != nil
	}
}

func AsFloat(val any) float32 {
	return float32(AsDouble(val))
}

func AsDouble(val any) float64 {
	switch val.(type) {
	case int8:
		return float64(val.(int8))
	case int16:
		return float64(val.(int8))
	case int32:
		return float64(val.(int8))
	case int:
		return float64(val.(int))
	case int64:
		return float64(val.(int64))
	case uint8:
		return float64(val.(uint8))
	case uint16:
		return float64(val.(uint16))
	case uint32:
		return float64(val.(uint32))
	case uint:
		return float64(val.(uint))
	case uint64:
		return float64(val.(uint64))
	case string:
		v, _ := strconv.ParseFloat(val.(string), 64)
		return v
	default:
		return AsDouble(AsStr(val))
	}
}

func AsInt(val any) int {
	return int(AsLong(val))
}

func AsUInt(val any) uint {
	return uint(AsULong(val))
}

func AsLong(val any) int64 {
	switch val.(type) {
	case int8:
		return int64(val.(int8))
	case int16:
		return int64(val.(int16))
	case int32:
		return int64(val.(int32))
	case int:
		return int64(val.(int))
	case int64:
		return val.(int64)
	case uint8:
		return int64(val.(uint8))
	case uint16:
		return int64(val.(uint16))
	case uint32:
		return int64(val.(uint32))
	case uint:
		return int64(val.(uint))
	case uint64:
		return int64(val.(uint64))
	case float32:
		return int64(val.(float32))
	case float64:
		return int64(val.(float64))
	case string:
		v, _ := strconv.Atoi(val.(string))
		return int64(v)
	default:
		return AsLong(AsStr(val))
	}
}

func AsULong(val any) uint64 {
	switch val.(type) {
	case int8:
		return uint64(val.(int8))
	case int16:
		return uint64(val.(int8))
	case int32:
		return uint64(val.(int8))
	case int:
		return uint64(val.(int))
	case int64:
		return uint64(val.(int64))
	case uint8:
		return uint64(val.(uint8))
	case uint16:
		return uint64(val.(uint16))
	case uint32:
		return uint64(val.(uint32))
	case uint:
		return uint64(val.(uint))
	case uint64:
		return val.(uint64)
	case float32:
		return uint64(val.(float32))
	case float64:
		return uint64(val.(float64))
	case string:
		v, _ := strconv.Atoi(val.(string))
		return uint64(v)
	default:
		return AsULong(AsStr(val))
	}
}

func AsMap(val any) map[string]any {
	obj := map[string]any{}
	switch val.(type) {
	case []byte, string:
		_ = Unmarshal(val.([]byte), &obj)
	default:
		obj = val.(map[string]any)
	}
	return obj
}

func AsStrMap(val any) map[string]string {
	obj := map[string]string{}
	switch val.(type) {
	case []byte, string:
		_ = Unmarshal(val.([]byte), &obj)
	default:
		obj = val.(map[string]string)
	}
	return obj
}

func AsArray(val any) []any {
	var arr []any
	switch val.(type) {
	case []byte, string:
		_ = Unmarshal(val.([]byte), &arr)
	default:
		arr = val.([]any)
	}
	return arr
}

func AsStrArr(val any) []string {
	var arr []string
	switch val.(type) {
	case []byte, string:
		_ = Unmarshal(val.([]byte), &arr)
	default:
		arr = val.([]string)
	}
	return arr
}

func IsLong(d float64) bool {
	offset := d - float64(int64(d))
	limit := EpsilonD
	return -limit <= offset && offset <= limit
}

func IsInt(f float32) bool {
	offset := f - float32(int(f))
	limit := EpsilonF
	return -limit <= offset && offset <= limit
}

func Dtoa(d float64) string {
	if IsLong(d) {
		return Ltoa(int64(d))
	}
	return fmt.Sprintf("%g", d)
}

func Ftoa(f float32) string {
	if IsInt(f) {
		return Itoa(int(f))
	}
	return fmt.Sprintf("%g", f)
}

func Ltoa(l int64) string {
	return fmt.Sprintf("%d", l)
}

func Itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

// 输入val(值),ratio(进位)
func RoundS(val string, precision int) string {
	part0, part1 := splitNum(val, precision)
	if precision > 0 {
		return fmt.Sprintf("%v.%v", part0, FillLeft(fmt.Sprint(part1), "0", precision))
	}
	newV := fmt.Sprintf("%v.%v", part0, FillLeft(fmt.Sprint(part1), "0", -precision))
	return AsStr(AsLong(math.Round(AsDouble(newV))))
}

func RoundV(val string, precision int) float64 {
	return AsDouble(RoundS(val, precision))
}

func RoundD(val float64, precision int) float64 {
	return RoundV(fmt.Sprint(val), precision)
}

func RoundL(val string, precision int) int64 {
	return int64(RoundV(val, precision))
}

func RoundI(val string, precision int) int {
	return int(RoundV(val, precision))
}

func splitNum(val string, precision int) (int64, int64) {
	if precision >= 0 {
		value := AsMoney(val, precision)
		ratio := int64(math.Pow10(precision))
		remain := AsMoney(val, precision+1) - value*10
		value += IfElse(remain >= 5, int64(1), int64(0)).(int64)
		return value / ratio, value % ratio
	} else {
		value := int64(AsDouble(val))
		ratio := int64(math.Pow10(-precision))
		remain := int64(AsDouble(val)*10) % 10
		value += IfElse(remain >= 5, int64(1), int64(0)).(int64)
		return value / ratio, value % ratio
	}
}

func FillLeft(raw, placeholder string, size int) string {
	v := raw
	for i := len(raw); i < size; i++ {
		v = placeholder + v
	}
	return v[:size]
}

func FillRight(raw, placeholder string, size int) string {
	v := raw
	for i := len(raw); i < size; i++ {
		v += placeholder
	}
	return v[:size]
}

// 输入val(值),ratio(进位)
func AsMoney(val string, precision int) int64 {
	val = strings.ToLower(val)
	if items := strings.Split(val, "e"); len(items) == 2 {
		vStr, pStr := items[0], items[1] // 获取值和进位的文本
		vVal := AsLong(strings.ReplaceAll(vStr, ".", ""))
		// 还需移位=科学进位+截断进位-已移位
		vEOffset := AsInt(pStr) /*科学进位*/ + precision /*截断进位*/ - (len(vStr) - 1 - strings.Index(vStr, ".")) /*已移位*/
		zoom := int64(math.Pow10(IfElse(vEOffset >= 0, vEOffset, 0-vEOffset).(int)))
		return IfElse(vEOffset >= 0, vVal*zoom, vVal/zoom).(int64)
	}
	replaceMap := map[string]string{"..": ".", ",": "", "元": ""}
	for k, v := range replaceMap {
		val = strings.ReplaceAll(val, k, v)
	}
	parts := strings.Split(val, ".")
	getPart := func(idx int) string {
		if idx >= 0 && idx < len(parts) {
			return parts[idx]
		}
		return "0"
	}
	iPart, fPart := getPart(0), FillRight(getPart(1), "0", precision)
	lVal := fmt.Sprintf("%v%v", iPart, fPart)
	return AsLong(lVal)
}

func GetStrMapKeys(src map[string]string) []string {
	ret := make([]string, len(src))
	i := 0
	for key := range src {
		ret[i] = key
		i++
	}
	return ret
}
