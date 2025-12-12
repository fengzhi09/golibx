package excelx

import (
	"context"
	"github.com/fengzhi09/golibx/jsonx"
	"strings"
)

type FieldType = string

type MapperConf struct {
	SheetIdx int         `json:"sheet_idx" toml:"sheet_idx" mapstructure:"sheet_idx"`
	ColLine  int         `json:"col_line" toml:"col_line" mapstructure:"col_line"`
	Mappers  []FieldConf `json:"mappers" toml:"mappers" mapstructure:"mappers"`
	Required []string    `json:"required" toml:"required" mapstructure:"required"`
}

type FieldOpts map[string][]FieldConf

func (mc *MapperConf) FieldOpts() FieldOpts {
	mappers := map[string][]FieldConf{}
	if mc != nil && mc.Mappers != nil {
		for _, mapper := range mc.Mappers {
			stdName := strings.ToLower(mapper.StdName)
			if _, hit := mappers[stdName]; !hit {
				mappers[stdName] = []FieldConf{}
			}
			mappers[stdName] = append(mappers[stdName], mapper)
			for _, alias := range mapper.Alias {
				aliasLow := strings.ToLower(alias)
				if _, hit := mappers[aliasLow]; !hit {
					mappers[alias] = []FieldConf{}
				}
				mappers[aliasLow] = append(mappers[aliasLow], mapper)
			}
		}
	}
	return mappers
}

func (mc *MapperConf) FindLikelyColumns(field string, columns []string) []string {
	columnsMaybe := make([]string, 0)
	if mappers, hit := mc.FieldOpts()[field]; hit {
		for _, mapper := range mappers {
			for _, column := range columns {
				if NameMatch(column, field, mapper.Alias) {
					columnsMaybe = append(columnsMaybe, column)
				}
			}
			if len(columnsMaybe) == 0 {
				columnsMaybe = append(columnsMaybe, mapper.Alias...)
			}
		}
	}
	if len(columnsMaybe) == 0 {
		columnsMaybe = append(columnsMaybe, field)
	}
	return columnsMaybe

}

type FieldConf struct {
	StdName      string            `json:"std_name" toml:"std_name" mapstructure:"std_name"` //标准名称
	Alias        []string          `json:"alias" toml:"alias" mapstructure:"alias"`          //别名
	InputType    FieldType         `json:"input_type" toml:"input_type" mapstructure:"input_type"`
	InputFormat  string            `json:"input_format" toml:"input_format" mapstructure:"input_format"`
	OutputType   FieldType         `json:"output_type" toml:"output_type" mapstructure:"output_type"`
	OutputFormat string            `json:"output_format" toml:"output_format" mapstructure:"output_format"`
	EnumMapper   map[string]any    `json:"enum_mapper" toml:"enum_mapper" mapstructure:"enum_mapper"` //枚举值映射
	EnumStrict   bool              `json:"enum_strict" toml:"enum_strict" mapstructure:"enum_strict"` //枚举值映射
	ZipFields    map[string]string `json:"zip_fields" toml:"zip_fields" mapstructure:"zip_fields"`    //字段压缩
	SkipWords    []string          `json:"skip_words" toml:"skip_words" mapstructure:"skip_words"`    //需要移除的字符串
	Prefix       string            `json:"prefix" toml:"prefix" mapstructure:"prefix"`                //需要增加的前缀
	Suffix       string            `json:"suffix" toml:"suffix" mapstructure:"suffix"`                //需要增加的后缀
	Regexp       string            `json:"regexp" toml:"regexp" mapstructure:"regexp"`                //正则
	CondEnum     map[string]string `json:"cond_enum" toml:"cond_enum" mapstructure:"cond_enum"`       //条件枚举
}

type Record struct {
	Idx     int
	Data    []string
	Columns []string
}

type Observer func(ctx context.Context, record *Record, conf *MapperConf, data *jsonx.JObj, raw *jsonx.JObj, err error)
