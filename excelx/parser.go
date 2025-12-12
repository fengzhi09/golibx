package excelx

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"
	"github.com/fengzhi09/golibx/logx"
)

func ForeachRow(ctx context.Context, xlsxFile *XlsxFile, sheetIdx, colLine int, recordCB func(record *Record)) error {
	defer xlsxFile.Release(ctx)
	sheetNames := xlsxFile.GetSheetList()
	if len(sheetNames) <= sheetIdx {
		return fmt.Errorf("idx too large")
	}

	rows, err := xlsxFile.Rows(sheetNames[sheetIdx])
	if err != nil {
		return err
	}
	defer gox.CloseRes(ctx, "rows", rows)
	rowIdx := 0
	columns := make([]string, 0)
	logx.Infof(ctx, "ForeachRow start")
	for rows.Next() {
		rowData, e := rows.Columns()
		if e != nil {
			return e
		}
		if len(rowData) <= 0 {
			continue
		}
		if rowIdx == colLine {
			for _, cell := range rowData {
				columns = append(columns, strings.TrimSpace(cell))
			}
		} else if rowIdx > colLine {
			recordCB(&Record{Idx: rowIdx, Data: rowData, Columns: columns})
		}
		rowIdx++
		if rowIdx%500 == 0 {
			logx.Infof(ctx, "ForeachRow now:%v", rowIdx)
		}
	}
	runtime.GC()
	logx.Infof(ctx, "ForeachRow done")
	return err
}

// 返回所有行的数据和列数组
func GetRecords(dataFile *XlsxFile, sheetIdx int) []*Record {
	records := make([]*Record, 0)
	_ = ForeachRow(context.Background(), dataFile, sheetIdx, 0, func(record *Record) {
		records = append(records, record)
	})
	return records
}

func ParseSheet(ctx context.Context, records []*Record, mapper *MapperConf, observer Observer) {
	fieldOpts := mapper.FieldOpts()
	for _, record := range records { // 把每一行每一列数据匹配上对应的key
		data, raw, err := ParseRecord(ctx, record, fieldOpts)
		if observer != nil {
			observer(ctx, record, mapper, data, raw, err)
		}
	}
}

func ParseSheetStream(ctx context.Context, dataFile *XlsxFile, sheetIdx int, mapper *MapperConf, observer Observer) {
	fieldOpts := mapper.FieldOpts()
	_ = ForeachRow(ctx, dataFile, sheetIdx, 0, func(record *Record) {
		data, raw, err := ParseRecord(ctx, record, fieldOpts)
		if observer != nil {
			observer(ctx, record, mapper, data, raw, err)
		}
	})
}

func IsDataEmpty(required []string, data *jsonx.JObj) (bool, string) {
	if len(required) == 0 {
		keys := data.Keys()
		for _, column := range keys {
			val := data.GetVal(column)
			isNull := val.IsNull()
			isEmptyObj := val.Type() == jsonx.JOBJ && val.ToObj().IsEmpty()
			isEmptyArr := val.Type() == jsonx.JARR && val.ToArr().IsEmpty()
			isEmptyStr := val.Type() == jsonx.JSTR && val.String() == ""
			if !(isNull || isEmptyObj || isEmptyArr || isEmptyStr) {
				return true, ""
			}
		}
		return false, "所有字段"
	} else {
		valid, emptyField, keys := true, "", required
		for _, key := range keys {
			if valid {
				val := data.GetVal(key)
				isNull := val.IsNull()
				isEmptyObj := val.Type() == jsonx.JOBJ && val.ToObj().IsEmpty()
				isEmptyArr := val.Type() == jsonx.JARR && val.ToArr().IsEmpty()
				isEmptyStr := val.Type() == jsonx.JSTR && val.String() == ""
				if isNull || isEmptyObj || isEmptyArr || isEmptyStr {
					valid, emptyField = false, key
				}
			}
		}
		return valid, emptyField
	}
}

func ParseRecord(ctx context.Context, record *Record, opts FieldOpts) (*jsonx.JObj, *jsonx.JObj, error) {
	data, raw := &jsonx.JObj{}, &jsonx.JObj{}
	var err error = nil
	emptyCnt := 0

	if len(record.Columns) == 0 {
		return data, raw, fmt.Errorf("空行")
	}
	for col, cell := range record.Data {
		if cell == "" {
			emptyCnt++
		}
		if col >= len(record.Columns) {
			return data, raw, fmt.Errorf("太多列")
		}
		column := record.Columns[col]
		column = strings.ToLower(column)
		fieldConfs, hit := opts[column]
		if hit {
			var lastErr error = nil
			errCnt := 0
			for _, fieldConf := range fieldConfs {

				fieldErr := ParseField(cell, column, fieldConf, data)
				if fieldErr != nil {
					lastErr = fieldErr
					errCnt++
				}
			}
			if errCnt == len(fieldConfs) {
				err = lastErr
			}
		}
		raw.PutStr(column, cell)
	}
	if emptyCnt == len(record.Columns) {
		return data, raw, fmt.Errorf("空行")
	}
	return data, raw, err
}

func NameMatch(column, stdName string, alias []string) bool {
	if alias == nil {
		alias = make([]string, 0)
	}
	column = strings.ToLower(column)
	alias = append(alias, stdName)
	for _, alias_ := range alias {
		alias_ = strings.TrimSpace(strings.ToLower(alias_))
		if alias_ == column {
			return true
		}
	}
	return false
}

func ParseField(cell string, column string, conf FieldConf, data *jsonx.JObj) error {
	field, typeWant := conf.StdName, conf.OutputType
	if !NameMatch(column, field, conf.Alias) {
		return fmt.Errorf("列'%v':字段不匹配", column)
	}
	for _, word := range conf.SkipWords {
		cell = strings.ReplaceAll(cell, word, "")
	}
	if conf.Prefix != "" && !strings.HasPrefix(cell, conf.Prefix) {
		cell = conf.Prefix + cell
	}
	if conf.Suffix != "" && !strings.HasSuffix(cell, conf.Suffix) {
		cell = cell + conf.Suffix
	}
	if getter, hit := getters[typeWant]; hit {
		val, err := getter(cell, conf)
		if err != nil {
			return fmt.Errorf("列'%v':%v", column, err)
		}
		data.Put(field, val)
		return nil
	} else {
		switch typeWant {
		case FT_ZipFields:
			obj := data.GetObj(field)
			if subField, yes := conf.ZipFields[column]; yes {
				obj.PutStr(subField, cell)
			} else {
				obj.Put(column, cell)
			}
			data.Put(field, obj)
			return nil
		}
	}
	return fmt.Errorf("列'%v':类型配置不支持", column)
}

func ParseMapper(path string) *MapperConf {
	ctx := context.Background()
	mapper := &MapperConf{}

	if err := gox.ReadJson(ctx, path, mapper); err != nil {
		logx.Errorf(ctx, "配置解析失败: %v, err:%v", path, err)
		return nil
	}
	logx.Debugf(ctx, "配置为 %v", mapper)
	return mapper
}

func ParseBool(raw string) bool {
	raw = strings.ToLower(raw)
	if raw == "yes" || raw == "是" || raw == "正确" || raw == "1" {
		return true
	}
	if raw == "no" || raw == "否" || raw == "错误" || raw == "0" || raw == "" {
		return false
	}
	return gox.AsBool(raw)
}

func ParseJArr(raw string, converter converter) *jsonx.JArr {
	arr := &jsonx.JArr{}
	if strings.HasPrefix(raw, "[") {
		*arr = jsonx.ParseJArr(raw)
	} else {
		parts := make([]string, 0)
		if len(parts) == 0 && strings.Contains(raw, ";") {
			parts = strings.Split(raw, ";")
		}
		if len(parts) == 0 && strings.Contains(raw, ",") {
			parts = strings.Split(raw, ",")
		}
		if len(parts) == 0 {
			parts = append(parts, raw)
		}

		for _, v := range parts {
			if converter != nil {
				arr.Add(converter(v))
			} else {
				arr.Add(v)
			}
		}
	}
	return arr
}
