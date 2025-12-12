package jsonx

import (
	"time"
)

type JType int

//goland:noinspection SpellCheckingInspection
const (
	JNULL JType = 0
	JBOOL JType = 1
	JINT  JType = 2
	JNUM  JType = 3
	JSTR  JType = 4
	JOBJ  JType = 5
	JARR  JType = 6
)

type JValue interface {
	Type() JType       // 数据类型
	IsNull() bool      // 是否空
	ToInt() int        // 强转为int
	ToLong() int64     // 强转为int64
	ToTime() time.Time // 强转为时间 time.Time
	ToDouble() float64 // 强转为float64
	ToFloat() float32  // 强转为float32
	ToBool() bool      // 强转为bool
	String() string    // 强转为string, 等效fmt.Sprint(v)
	Pretty() string    // 强转为string, obj,arr 会优化展示，其他等同String()
	ToObj() JObj       // 强转为对象，JStr 会当做json字符串解析
	ToObjPtr() *JObj   // 强转为对象指针
	ToArr() JArr       // 强转为数组，JStr 会当做json字符串解析
	ToJDoc() JDoc      // 安装json字符串格式输出，JStr 会转义加""
	ToJVal() JValue    // 返回接口操作对象
	ToGVal() any       // 返回底层对象
}

func JValEqJVal(a, b JValue) bool {
	return a.ToGVal() == b.ToGVal()
}
func JValEqGVal(a JValue, b any) bool {
	return a.ToGVal() == b
}
