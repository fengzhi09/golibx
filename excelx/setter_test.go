package excelx

import (
	"github.com/dlclark/regexp2"
	"testing"
)

func TestRegexp(t *testing.T) {
	conf := FieldConf{
		StdName:      "测试",
		Alias:        nil,
		InputType:    "",
		InputFormat:  "",
		OutputType:   "",
		OutputFormat: "",
		EnumMapper:   nil,
		EnumStrict:   false,
		ZipFields:    nil,
		SkipWords:    nil,
		Prefix:       "",
		Suffix:       "",
		Regexp:       "",
		CondEnum:     map[string]string{"yes": "顾客留言", "no": "^((?!顾客留言).)*$"},
	}
	cells := []string{"客服反馈过滤", "顾客留言过滤"}
	for _, cell := range cells {
		hits := false
		var err error
		for enum, reg := range conf.CondEnum {
			regMachine, _ := regexp2.Compile(reg, 0)
			hit, _ := regMachine.MatchString(cell)
			if hit {
				hits = true
				t.Logf("%v成功命中:%v", cell, enum)
			}
		}
		if !hits {
			t.Fatalf("正则命中失败：%v,%v", cell, err)
		}
	}
}
